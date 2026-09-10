package repos

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"example.com/german/backend/internal/models"
	"github.com/lib/pq"
)

type PublicCategoryRecord struct {
	ID   int64
	Name string
}

type PublicAnswerRecord struct {
	ID   int64
	Text string
}

type PublicQuestionRecord struct {
	ID      int64
	Text    string
	Answers []PublicAnswerRecord
}

type TicketAnswerRecord struct {
	QuestionID int64
	AnswerID   int64
}

type CreateTicketRecord struct {
	TrackID       string
	ApplicantType models.ApplicantType
	CategoryID    int64
	Priority      models.TicketPriority
	Status        models.TicketStatus
	Description   string
	Answers       []TicketAnswerRecord
	Attachments   []NewAttachmentRecord
	CrisisContact string
	Crisis        bool
	CreatedAt     time.Time
}

type CreatedTicketRecord struct {
	TrackID   string
	Status    models.TicketStatus
	CreatedAt time.Time
}

type CreatedReviewRecord struct {
	ID       int64
	TicketID int64
	Rating   int
	Text     string
}

func (repository *TicketRepository) EnsureCategory(ctx context.Context, name string) error {
	_, err := repository.db.ExecContext(ctx, `
		INSERT INTO categories (name)
		SELECT $1::text
		WHERE NOT EXISTS (
			SELECT 1 FROM categories WHERE lower(name) = lower($1::text)
		)`, name)
	if err != nil {
		return fmt.Errorf("ensure required category: %w", err)
	}
	return nil
}

func (repository *TicketRepository) ListPublicCategories(ctx context.Context) ([]PublicCategoryRecord, error) {
	rows, err := repository.db.QueryContext(ctx, `SELECT id, name FROM categories ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("list public categories: %w", err)
	}
	defer rows.Close()

	result := make([]PublicCategoryRecord, 0)
	for rows.Next() {
		var category PublicCategoryRecord
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			return nil, fmt.Errorf("scan public category: %w", err)
		}
		result = append(result, category)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate public categories: %w", err)
	}
	return result, nil
}

func (repository *TicketRepository) ListPublicQuestions(ctx context.Context, categoryID int64) ([]PublicQuestionRecord, error) {
	var exists bool
	if err := repository.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM categories WHERE id = $1)`, categoryID).Scan(&exists); err != nil {
		return nil, fmt.Errorf("check public category: %w", err)
	}
	if !exists {
		return nil, ErrCategoryNotFound
	}

	rows, err := repository.db.QueryContext(ctx, `
		SELECT q.id, q.text, a.id, a.text
		FROM questions AS q
		LEFT JOIN answers AS a ON a.question_id = q.id
		WHERE q.cat_id = $1
		ORDER BY q.id, a.id`, categoryID)
	if err != nil {
		return nil, fmt.Errorf("list public questions: %w", err)
	}
	defer rows.Close()

	result := make([]PublicQuestionRecord, 0)
	positions := make(map[int64]int)
	for rows.Next() {
		var questionID int64
		var questionText string
		var answerID sql.NullInt64
		var answerText sql.NullString
		if err := rows.Scan(&questionID, &questionText, &answerID, &answerText); err != nil {
			return nil, fmt.Errorf("scan public question: %w", err)
		}
		position, ok := positions[questionID]
		if !ok {
			position = len(result)
			positions[questionID] = position
			result = append(result, PublicQuestionRecord{ID: questionID, Text: questionText, Answers: []PublicAnswerRecord{}})
		}
		if answerID.Valid {
			result[position].Answers = append(result[position].Answers, PublicAnswerRecord{ID: answerID.Int64, Text: answerText.String})
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate public questions: %w", err)
	}
	return result, nil
}

func (repository *TicketRepository) SelectedAnswerTexts(ctx context.Context, categoryID int64, answers []TicketAnswerRecord) ([]string, error) {
	texts := make([]string, 0, len(answers))
	for _, answer := range answers {
		var text string
		err := repository.db.QueryRowContext(ctx, `
			SELECT a.text
			FROM questions AS q
			JOIN answers AS a ON a.question_id = q.id
			WHERE q.id = $1 AND a.id = $2 AND q.cat_id = $3`, answer.QuestionID, answer.AnswerID, categoryID).Scan(&text)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidAnswer
		}
		if err != nil {
			return nil, fmt.Errorf("get selected answer text: %w", err)
		}
		texts = append(texts, text)
	}
	return texts, nil
}

func (repository *TicketRepository) CreatePublicTicket(ctx context.Context, input CreateTicketRecord) (CreatedTicketRecord, error) {
	tx, err := repository.db.BeginTx(ctx, nil)
	if err != nil {
		return CreatedTicketRecord{}, fmt.Errorf("begin create ticket transaction: %w", err)
	}
	defer tx.Rollback()

	var categoryExists bool
	if err := tx.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM categories WHERE id = $1)`, input.CategoryID).Scan(&categoryExists); err != nil {
		return CreatedTicketRecord{}, fmt.Errorf("check ticket category: %w", err)
	}
	if !categoryExists {
		return CreatedTicketRecord{}, ErrCategoryNotFound
	}

	for _, answer := range input.Answers {
		var valid bool
		if err := tx.QueryRowContext(ctx, `
			SELECT EXISTS(
				SELECT 1
				FROM questions AS q
				JOIN answers AS a ON a.question_id = q.id
				WHERE q.id = $1 AND a.id = $2 AND q.cat_id = $3
			)`, answer.QuestionID, answer.AnswerID, input.CategoryID).Scan(&valid); err != nil {
			return CreatedTicketRecord{}, fmt.Errorf("validate ticket answer: %w", err)
		}
		if !valid {
			return CreatedTicketRecord{}, ErrInvalidAnswer
		}
	}

	var ticketID int64
	var createdAt time.Time
	err = tx.QueryRowContext(ctx, `
		INSERT INTO tickets (track_id, priority, morda_type, category_id, status, return_count, created_at)
		VALUES ($1, $2, $3, $4, $5, 0, $6)
		RETURNING id, created_at`, input.TrackID, input.Priority, input.ApplicantType, input.CategoryID, input.Status, input.CreatedAt).Scan(&ticketID, &createdAt)
	if err != nil {
		var postgresError *pq.Error
		if errors.As(err, &postgresError) && postgresError.Code == "23505" {
			return CreatedTicketRecord{}, ErrTrackIDExists
		}
		return CreatedTicketRecord{}, fmt.Errorf("insert ticket: %w", err)
	}

	if input.Description != "" {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO messages (ticket_id, text, type, created_at)
			VALUES ($1, $2, $3, $4)`, ticketID, input.Description, models.MessageTypeApplicant, input.CreatedAt); err != nil {
			return CreatedTicketRecord{}, fmt.Errorf("insert initial ticket message: %w", err)
		}
	}

	for _, answer := range input.Answers {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO "QA" (question_id, answer_id, ticket_id)
			VALUES ($1, $2, $3)`, answer.QuestionID, answer.AnswerID, ticketID); err != nil {
			return CreatedTicketRecord{}, fmt.Errorf("insert ticket answer: %w", err)
		}
	}

	for _, attachment := range input.Attachments {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO attachments (
				ticket_id, storage_key, safe_name, mime_type,
				original_size_bytes, size_bytes, width, height, created_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
			ticketID, attachment.StorageKey, attachment.SafeName, attachment.MIMEType,
			attachment.OriginalSizeBytes, attachment.SizeBytes, attachment.Width, attachment.Height, input.CreatedAt,
		); err != nil {
			return CreatedTicketRecord{}, fmt.Errorf("insert initial ticket attachment: %w", err)
		}
	}

	if input.Crisis && input.CrisisContact != "" {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO crisis_contact (contact, ticket_id, created_at)
			VALUES ($1, $2, $3)`, input.CrisisContact, ticketID, input.CreatedAt); err != nil {
			return CreatedTicketRecord{}, fmt.Errorf("insert crisis contact: %w", err)
		}
	}
	if input.Crisis {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO ticket_events (ticket_id, event_type, to_status, reason_code, created_at)
			VALUES ($1, 'crisis_detected', $2, 'dictionary_match', $3)`, ticketID, input.Status, input.CreatedAt); err != nil {
			return CreatedTicketRecord{}, fmt.Errorf("insert crisis event: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return CreatedTicketRecord{}, fmt.Errorf("commit create ticket transaction: %w", err)
	}
	return CreatedTicketRecord{TrackID: input.TrackID, Status: input.Status, CreatedAt: createdAt.UTC()}, nil
}

func (repository *TicketRepository) CreateReview(ctx context.Context, trackID string, rating int, text string) (CreatedReviewRecord, error) {
	tx, err := repository.db.BeginTx(ctx, nil)
	if err != nil {
		return CreatedReviewRecord{}, fmt.Errorf("begin create review transaction: %w", err)
	}
	defer tx.Rollback()

	var ticketID int64
	var status models.TicketStatus
	err = tx.QueryRowContext(ctx, `SELECT id, status FROM tickets WHERE track_id = $1 FOR UPDATE`, trackID).Scan(&ticketID, &status)
	if errors.Is(err, sql.ErrNoRows) {
		return CreatedReviewRecord{}, ErrTicketNotFound
	}
	if err != nil {
		return CreatedReviewRecord{}, fmt.Errorf("lock ticket for review: %w", err)
	}
	if status != models.TicketStatusCompleted {
		return CreatedReviewRecord{}, ErrReviewNotAllowed
	}

	var exists bool
	if err := tx.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM reviews WHERE ticket_id = $1)`, ticketID).Scan(&exists); err != nil {
		return CreatedReviewRecord{}, fmt.Errorf("check existing review: %w", err)
	}
	if exists {
		return CreatedReviewRecord{}, ErrReviewAlreadyExists
	}

	result := CreatedReviewRecord{TicketID: ticketID, Rating: rating, Text: text}
	if err := tx.QueryRowContext(ctx, `
		INSERT INTO reviews (ticket_id, rating, text)
		VALUES ($1, $2, $3)
		RETURNING id`, ticketID, rating, text).Scan(&result.ID); err != nil {
		return CreatedReviewRecord{}, fmt.Errorf("insert review: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return CreatedReviewRecord{}, fmt.Errorf("commit review transaction: %w", err)
	}
	return result, nil
}
