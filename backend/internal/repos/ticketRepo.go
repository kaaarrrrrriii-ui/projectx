package repos

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"example.com/german/backend/internal/models"
)

const maxTicketReturns = 2

type TicketRepository struct {
	db *sql.DB
}

type TicketStatusRecord struct {
	ID          int64
	TrackID     string
	Status      models.TicketStatus
	CreatedAt   time.Time
	CategoryID  int64
	Category    string
	ReturnCount int
	Resolution  sql.NullString
}

type ChatMessageRecord struct {
	ID        int64
	Text      string
	Type      models.MessageType
	CreatedAt time.Time
}

type ChatAttachmentRecord struct {
	ID        int64
	SafeName  string
	MIMEType  string
	SizeBytes int64
	CreatedAt time.Time
}

type ChatSpecialistRecord struct {
	Label       string
	ExpertGroup string
}

type TicketChatRecord struct {
	TrackID     string
	Status      models.TicketStatus
	Specialist  *ChatSpecialistRecord
	Messages    []ChatMessageRecord
	Attachments []ChatAttachmentRecord
}

type CompleteTicketRecord struct {
	Status   models.TicketStatus
	ClosedAt time.Time
}

type ReturnTicketRecord struct {
	Status      models.TicketStatus
	ReturnCount int
}

func NewTicketRepository(db *sql.DB) *TicketRepository {
	return &TicketRepository{db: db}
}

func (repository *TicketRepository) GetStatusByTrack(ctx context.Context, trackID string) (TicketStatusRecord, error) {
	const query = `
		SELECT t.id, t.track_id, t.status, t.created_at,
		       c.id, c.name, t.return_count,
		       (
		           SELECT m.text FROM messages AS m
		           WHERE m.ticket_id = t.id AND m.type = $2
		           ORDER BY m.created_at DESC, m.id DESC LIMIT 1
		       )
		FROM tickets AS t
		JOIN categories AS c ON c.id = t.category_id
	WHERE t.track_id = $1`

	var record TicketStatusRecord
	err := repository.db.QueryRowContext(ctx, query, trackID, models.MessageTypeOperator).Scan(
		&record.ID,
		&record.TrackID,
		&record.Status,
		&record.CreatedAt,
		&record.CategoryID,
		&record.Category,
		&record.ReturnCount,
		&record.Resolution,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return TicketStatusRecord{}, ErrTicketNotFound
	}
	if err != nil {
		return TicketStatusRecord{}, fmt.Errorf("get ticket status: %w", err)
	}
	return record, nil
}

func (repository *TicketRepository) GetChatByTrack(ctx context.Context, trackID string) (TicketChatRecord, error) {
	const ticketQuery = `
		SELECT id, track_id, status
		FROM tickets
		WHERE track_id = $1`

	var ticketID int64
	var record TicketChatRecord
	err := repository.db.QueryRowContext(ctx, ticketQuery, trackID).Scan(&ticketID, &record.TrackID, &record.Status)
	if errors.Is(err, sql.ErrNoRows) {
		return TicketChatRecord{}, ErrTicketNotFound
	}
	if err != nil {
		return TicketChatRecord{}, fmt.Errorf("get ticket for chat: %w", err)
	}

	const specialistQuery = `
		SELECT eg.title
		FROM tickets_workers AS tw
		JOIN users AS u ON u.id = tw.worker_id
		JOIN expert_groups AS eg ON eg.id = u.expert_group_id
		WHERE tw.ticket_id = $1
		  AND tw.actual
		  AND tw.is_responsible
		LIMIT 1`
	var expertGroup string
	err = repository.db.QueryRowContext(ctx, specialistQuery, ticketID).Scan(&expertGroup)
	if err == nil {
		record.Specialist = &ChatSpecialistRecord{Label: "Специалист", ExpertGroup: expertGroup}
	} else if !errors.Is(err, sql.ErrNoRows) {
		return TicketChatRecord{}, fmt.Errorf("get chat specialist: %w", err)
	}

	const messagesQuery = `
		SELECT id, text, type, created_at
		FROM messages
		WHERE ticket_id = $1
		  AND type IN ($2, $3)
		ORDER BY created_at, id`
	rows, err := repository.db.QueryContext(
		ctx,
		messagesQuery,
		ticketID,
		models.MessageTypeApplicant,
		models.MessageTypeSpecialist,
	)
	if err != nil {
		return TicketChatRecord{}, fmt.Errorf("list chat messages: %w", err)
	}
	defer rows.Close()

	record.Messages = make([]ChatMessageRecord, 0)
	for rows.Next() {
		var message ChatMessageRecord
		if err := rows.Scan(&message.ID, &message.Text, &message.Type, &message.CreatedAt); err != nil {
			return TicketChatRecord{}, fmt.Errorf("scan chat message: %w", err)
		}
		record.Messages = append(record.Messages, message)
	}
	if err := rows.Err(); err != nil {
		return TicketChatRecord{}, fmt.Errorf("iterate chat messages: %w", err)
	}

	const attachmentsQuery = `
		SELECT id, safe_name, mime_type, size_bytes, created_at
		FROM attachments
		WHERE ticket_id = $1
		ORDER BY created_at, id`
	attachmentRows, err := repository.db.QueryContext(ctx, attachmentsQuery, ticketID)
	if err != nil {
		return TicketChatRecord{}, fmt.Errorf("list chat attachments: %w", err)
	}
	defer attachmentRows.Close()

	record.Attachments = make([]ChatAttachmentRecord, 0)
	for attachmentRows.Next() {
		var attachment ChatAttachmentRecord
		if err := attachmentRows.Scan(
			&attachment.ID,
			&attachment.SafeName,
			&attachment.MIMEType,
			&attachment.SizeBytes,
			&attachment.CreatedAt,
		); err != nil {
			return TicketChatRecord{}, fmt.Errorf("scan chat attachment: %w", err)
		}
		record.Attachments = append(record.Attachments, attachment)
	}
	if err := attachmentRows.Err(); err != nil {
		return TicketChatRecord{}, fmt.Errorf("iterate chat attachments: %w", err)
	}

	return record, nil
}

func (repository *TicketRepository) Complete(ctx context.Context, trackID string, closedAt time.Time) (CompleteTicketRecord, error) {
	tx, err := repository.db.BeginTx(ctx, nil)
	if err != nil {
		return CompleteTicketRecord{}, fmt.Errorf("begin complete ticket transaction: %w", err)
	}
	defer tx.Rollback()

	const selectQuery = `
		SELECT status, closed_at
		FROM tickets
		WHERE track_id = $1
		FOR UPDATE`
	var status models.TicketStatus
	var existingClosedAt sql.NullTime
	err = tx.QueryRowContext(ctx, selectQuery, trackID).Scan(&status, &existingClosedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return CompleteTicketRecord{}, ErrTicketNotFound
	}
	if err != nil {
		return CompleteTicketRecord{}, fmt.Errorf("lock ticket for completion: %w", err)
	}

	if status == models.TicketStatusCompleted && existingClosedAt.Valid {
		if err := tx.Commit(); err != nil {
			return CompleteTicketRecord{}, fmt.Errorf("commit completed ticket read: %w", err)
		}
		return CompleteTicketRecord{Status: status, ClosedAt: existingClosedAt.Time.UTC()}, nil
	}
	if !status.IsApplicantCompletable() {
		return CompleteTicketRecord{}, ErrInvalidTransition
	}

	const updateQuery = `
		UPDATE tickets
		SET status = $2, closed_at = $3
		WHERE track_id = $1`
	if _, err := tx.ExecContext(ctx, updateQuery, trackID, models.TicketStatusCompleted, closedAt); err != nil {
		return CompleteTicketRecord{}, fmt.Errorf("complete ticket: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return CompleteTicketRecord{}, fmt.Errorf("commit complete ticket transaction: %w", err)
	}
	return CompleteTicketRecord{Status: models.TicketStatusCompleted, ClosedAt: closedAt.UTC()}, nil
}

func (repository *TicketRepository) Return(ctx context.Context, trackID, reason string, returnedAt time.Time) (ReturnTicketRecord, error) {
	tx, err := repository.db.BeginTx(ctx, nil)
	if err != nil {
		return ReturnTicketRecord{}, fmt.Errorf("begin return ticket transaction: %w", err)
	}
	defer tx.Rollback()

	const selectQuery = `
		SELECT id, status, return_count
		FROM tickets
		WHERE track_id = $1
		FOR UPDATE`
	var ticketID int64
	var status models.TicketStatus
	var returnCount int
	err = tx.QueryRowContext(ctx, selectQuery, trackID).Scan(&ticketID, &status, &returnCount)
	if errors.Is(err, sql.ErrNoRows) {
		return ReturnTicketRecord{}, ErrTicketNotFound
	}
	if err != nil {
		return ReturnTicketRecord{}, fmt.Errorf("lock ticket for return: %w", err)
	}
	if status != models.TicketStatusAnswerReady {
		return ReturnTicketRecord{}, ErrInvalidTransition
	}
	if returnCount >= maxTicketReturns {
		return ReturnTicketRecord{}, ErrReturnLimitReached
	}

	if reason != "" {
		const messageQuery = `
			INSERT INTO messages (ticket_id, text, type, created_at)
			VALUES ($1, $2, $3, $4)`
		if _, err := tx.ExecContext(ctx, messageQuery, ticketID, reason, models.MessageTypeReturnReason, returnedAt); err != nil {
			return ReturnTicketRecord{}, fmt.Errorf("save ticket return reason: %w", err)
		}
	}

	returnCount++
	const updateQuery = `
		UPDATE tickets
		SET status = $2, return_count = $3, closed_at = NULL
		WHERE id = $1`
	if _, err := tx.ExecContext(ctx, updateQuery, ticketID, models.TicketStatusReturned, returnCount); err != nil {
		return ReturnTicketRecord{}, fmt.Errorf("return ticket: %w", err)
	}
	var previousWorkerID sql.NullInt64
	if err := tx.QueryRowContext(ctx, `
		SELECT worker_id FROM tickets_workers
		WHERE ticket_id = $1 AND actual AND is_responsible
		LIMIT 1`, ticketID).Scan(&previousWorkerID); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return ReturnTicketRecord{}, fmt.Errorf("get previous worker for return: %w", err)
	}
	if _, err := tx.ExecContext(ctx, `UPDATE tickets_workers SET actual = FALSE WHERE ticket_id = $1 AND actual`, ticketID); err != nil {
		return ReturnTicketRecord{}, fmt.Errorf("deactivate returned ticket workers: %w", err)
	}
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO ticket_events (
			ticket_id, actor_user_id, event_type, from_status, to_status,
			from_worker_id, created_at
		) VALUES ($1, NULL, 'returned', $2, $3, $4, $5)`,
		ticketID, status, models.TicketStatusReturned, previousWorkerID, returnedAt); err != nil {
		return ReturnTicketRecord{}, fmt.Errorf("save ticket return event: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return ReturnTicketRecord{}, fmt.Errorf("commit return ticket transaction: %w", err)
	}
	return ReturnTicketRecord{Status: models.TicketStatusReturned, ReturnCount: returnCount}, nil
}
