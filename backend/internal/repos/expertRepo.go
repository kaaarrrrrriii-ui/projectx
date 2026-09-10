package repos

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"example.com/german/backend/internal/models"
)

type ExpertRepository struct {
	db *sql.DB
}

func NewExpertRepository(db *sql.DB) *ExpertRepository {
	return &ExpertRepository{db: db}
}

type ExpertProfileRecord struct {
	ID            int64
	FullName      string
	GroupID       int64
	GroupTitle    string
	ActiveTickets int
	MaxTickets    int
	AverageRating float64
}

type ExpertTicketFilter struct {
	WorkerID      int64
	Queue         string
	Search        string
	Priority      models.TicketPriority
	CategoryID    int64
	ApplicantType models.ApplicantType
	Limit         int
	Offset        int
}

type ExpertTicketRecord struct {
	TrackID        string
	CategoryID     int64
	Category       string
	Status         models.TicketStatus
	ApplicantType  models.ApplicantType
	Priority       models.TicketPriority
	CreatedAt      time.Time
	WaitingSeconds int64
	ReturnCount    int
	IsResponsible  bool
}

type ExpertTicketPageRecord struct {
	Items []ExpertTicketRecord
	Total int
}

type ExpertClarificationRecord struct {
	QuestionID int64
	Question   string
	AnswerID   int64
	Answer     string
}

type ExpertWorkerRecord struct {
	ID            int64
	FullName      string
	GroupTitle    string
	IsResponsible bool
}

type ExpertNoteRecord struct {
	ID         int64
	Text       string
	AuthorID   int64
	AuthorName string
	CreatedAt  time.Time
}

type ExpertTicketDetailRecord struct {
	TrackID        string
	Status         models.TicketStatus
	Priority       models.TicketPriority
	CreatedAt      time.Time
	CategoryID     int64
	Category       string
	ApplicantType  models.ApplicantType
	Description    string
	ReturnCount    int
	ReturnReason   sql.NullString
	Clarifications []ExpertClarificationRecord
	Attachments    []ChatAttachmentRecord
	Messages       []ChatMessageRecord
	Workers        []ExpertWorkerRecord
	Notes          []ExpertNoteRecord
}

type ExpertMessageRecord struct {
	ID           int64
	Text         string
	Type         models.MessageType
	CreatedAt    time.Time
	TicketStatus models.TicketStatus
}

type WorkerRequestPayload struct {
	Kind        string `json:"kind"`
	Version     int    `json:"version"`
	RequestType string `json:"request_type"`
	Status      string `json:"status"`
	Reason      string `json:"reason"`
	CreatedBy   int64  `json:"created_by"`
}

type WorkerRequestStatusPayload struct {
	Kind        string `json:"kind"`
	Version     int    `json:"version"`
	RequestID   int64  `json:"request_id"`
	Status      string `json:"status"`
	CompletedBy int64  `json:"completed_by"`
}

type WorkerRequestRecord struct {
	ID            int64
	TicketID      int64
	TrackID       string
	CategoryID    int64
	Category      string
	ApplicantType models.ApplicantType
	Priority      models.TicketPriority
	RequestType   string
	Reason        string
	Status        string
	CreatedAt     time.Time
	CreatedBy     int64
}

type WorkerRequestPageRecord struct {
	Items []WorkerRequestRecord
	Total int
}

type ExpertDashboardRecord struct {
	QueueCount      int
	InProgressCount int
	ReturnedCount   int
	UrgentCount     int
	ProcessedCount  int
}

func (repository *ExpertRepository) Profile(ctx context.Context, workerID int64) (ExpertProfileRecord, error) {
	var record ExpertProfileRecord
	err := repository.db.QueryRowContext(ctx, `
		SELECT u.id, u.full_name, eg.id, eg.title,
		       COUNT(tw.ticket_id) FILTER (
		           WHERE tw.actual AND t.status IN ($2, $3, $4, $5)
		       ), u.max_tickets, u.avg_rating
		FROM users AS u
		JOIN expert_groups AS eg ON eg.id = u.expert_group_id
		LEFT JOIN tickets_workers AS tw ON tw.worker_id = u.id
		LEFT JOIN tickets AS t ON t.id = tw.ticket_id
		WHERE u.id = $1 AND u.role = 'expert'
		GROUP BY u.id, u.full_name, eg.id, eg.title, u.max_tickets, u.avg_rating`,
		workerID,
		models.TicketStatusAssigned,
		models.TicketStatusInProgress,
		models.TicketStatusNeedsClarification,
		models.TicketStatusAnswerReady,
	).Scan(
		&record.ID, &record.FullName, &record.GroupID, &record.GroupTitle,
		&record.ActiveTickets, &record.MaxTickets, &record.AverageRating,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return ExpertProfileRecord{}, ErrWorkerNotFound
	}
	if err != nil {
		return ExpertProfileRecord{}, fmt.Errorf("get expert profile: %w", err)
	}
	return record, nil
}

func (repository *ExpertRepository) Dashboard(ctx context.Context, workerID int64) (ExpertDashboardRecord, error) {
	var record ExpertDashboardRecord
	err := repository.db.QueryRowContext(ctx, `
		SELECT
			COUNT(*) FILTER (WHERE tw.actual AND t.status = $2 AND t.return_count = 0),
			COUNT(*) FILTER (WHERE tw.actual AND t.status IN ($3, $4, $5)),
			COUNT(*) FILTER (WHERE tw.actual AND t.status = $2 AND t.return_count > 0),
			COUNT(*) FILTER (WHERE tw.actual AND t.priority = $6 AND t.status IN ($2, $3, $4, $5)),
			COUNT(DISTINCT t.id) FILTER (WHERE EXISTS (
				SELECT 1 FROM ticket_events AS e
				WHERE e.ticket_id = t.id AND e.actor_user_id = $1
				  AND e.event_type = 'specialist_answered'
			))
		FROM tickets_workers AS tw
		JOIN tickets AS t ON t.id = tw.ticket_id
		WHERE tw.worker_id = $1`,
		workerID,
		models.TicketStatusAssigned,
		models.TicketStatusInProgress,
		models.TicketStatusNeedsClarification,
		models.TicketStatusAnswerReady,
		models.TicketPriorityUrgent,
	).Scan(&record.QueueCount, &record.InProgressCount, &record.ReturnedCount, &record.UrgentCount, &record.ProcessedCount)
	if err != nil {
		return ExpertDashboardRecord{}, fmt.Errorf("get expert dashboard: %w", err)
	}
	return record, nil
}

func (repository *ExpertRepository) ListTickets(ctx context.Context, filter ExpertTicketFilter) (ExpertTicketPageRecord, error) {
	where, args := expertTicketWhere(filter)
	var total int
	if err := repository.db.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM tickets AS t
		JOIN tickets_workers AS tw ON tw.ticket_id = t.id
		`+where, args...).Scan(&total); err != nil {
		return ExpertTicketPageRecord{}, fmt.Errorf("count expert tickets: %w", err)
	}

	args = append(args, filter.Limit, filter.Offset)
	limitPosition := len(args) - 1
	offsetPosition := len(args)
	rows, err := repository.db.QueryContext(ctx, `
		SELECT t.track_id, c.id, c.name, t.status, t.morda_type, t.priority,
		       t.created_at,
		       GREATEST(0, EXTRACT(EPOCH FROM (CURRENT_TIMESTAMP - t.created_at)))::BIGINT,
		       t.return_count, tw.is_responsible
		FROM tickets AS t
		JOIN tickets_workers AS tw ON tw.ticket_id = t.id
		JOIN categories AS c ON c.id = t.category_id
		`+where+`
		ORDER BY CASE t.priority WHEN 2 THEN 0 WHEN 1 THEN 1 ELSE 2 END,
		         t.created_at ASC, t.id ASC
		LIMIT $`+fmt.Sprint(limitPosition)+` OFFSET $`+fmt.Sprint(offsetPosition), args...)
	if err != nil {
		return ExpertTicketPageRecord{}, fmt.Errorf("list expert tickets: %w", err)
	}
	defer rows.Close()
	result := ExpertTicketPageRecord{Items: make([]ExpertTicketRecord, 0), Total: total}
	for rows.Next() {
		var item ExpertTicketRecord
		if err := rows.Scan(
			&item.TrackID, &item.CategoryID, &item.Category, &item.Status,
			&item.ApplicantType, &item.Priority, &item.CreatedAt,
			&item.WaitingSeconds, &item.ReturnCount, &item.IsResponsible,
		); err != nil {
			return ExpertTicketPageRecord{}, fmt.Errorf("scan expert ticket: %w", err)
		}
		result.Items = append(result.Items, item)
	}
	if err := rows.Err(); err != nil {
		return ExpertTicketPageRecord{}, fmt.Errorf("iterate expert tickets: %w", err)
	}
	return result, nil
}

func expertTicketWhere(filter ExpertTicketFilter) (string, []any) {
	conditions := []string{"tw.worker_id = $1", "tw.actual"}
	args := []any{filter.WorkerID}
	add := func(template string, value any) {
		args = append(args, value)
		conditions = append(conditions, fmt.Sprintf(template, len(args)))
	}
	switch filter.Queue {
	case "queue":
		conditions = append(conditions, fmt.Sprintf("t.status = %d AND t.return_count = 0", models.TicketStatusAssigned))
	case "assigned":
		conditions = append(conditions, fmt.Sprintf("t.status IN (%d, %d, %d)", models.TicketStatusInProgress, models.TicketStatusNeedsClarification, models.TicketStatusAnswerReady))
	case "returned":
		conditions = append(conditions, fmt.Sprintf("t.status = %d AND t.return_count > 0", models.TicketStatusAssigned))
	}
	if filter.Search != "" {
		add("t.track_id ILIKE '%%' || $%d || '%%'", filter.Search)
	}
	if filter.Priority != 0 {
		add("t.priority = $%d", filter.Priority)
	}
	if filter.CategoryID != 0 {
		add("t.category_id = $%d", filter.CategoryID)
	}
	if filter.ApplicantType != 0 {
		add("t.morda_type = $%d", filter.ApplicantType)
	}
	return "WHERE " + strings.Join(conditions, " AND "), args
}

func (repository *ExpertRepository) OpenTicket(ctx context.Context, trackID string, workerID int64, openedAt time.Time) (models.TicketStatus, error) {
	tx, err := repository.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("begin expert open transaction: %w", err)
	}
	defer tx.Rollback()
	ticketID, status, _, err := lockExpertTicket(ctx, tx, trackID, workerID)
	if err != nil {
		return 0, err
	}
	if status == models.TicketStatusAssigned {
		if _, err := tx.ExecContext(ctx, `UPDATE tickets SET status = $2 WHERE id = $1`, ticketID, models.TicketStatusInProgress); err != nil {
			return 0, fmt.Errorf("start expert ticket: %w", err)
		}
		if err := insertEvent(ctx, tx, ticketID, workerID, "expert_opened", status, models.TicketStatusInProgress, sql.NullInt64{}, sql.NullInt64{}, openedAt); err != nil {
			return 0, err
		}
		status = models.TicketStatusInProgress
	}
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("commit expert open transaction: %w", err)
	}
	return status, nil
}

func (repository *ExpertRepository) GetTicket(ctx context.Context, trackID string, workerID int64) (ExpertTicketDetailRecord, error) {
	var ticketID int64
	var record ExpertTicketDetailRecord
	err := repository.db.QueryRowContext(ctx, `
		SELECT t.id, t.track_id, t.status, t.priority, t.created_at,
		       c.id, c.name, t.morda_type, t.return_count,
		       COALESCE((
		           SELECT m.text FROM messages AS m
		           WHERE m.ticket_id = t.id AND m.type = $3
		           ORDER BY m.created_at, m.id LIMIT 1
		       ), ''),
		       (
		           SELECT m.text FROM messages AS m
		           WHERE m.ticket_id = t.id AND m.type = $4
		           ORDER BY m.created_at DESC, m.id DESC LIMIT 1
		       )
		FROM tickets AS t
		JOIN categories AS c ON c.id = t.category_id
		JOIN tickets_workers AS tw ON tw.ticket_id = t.id
		WHERE t.track_id = $1 AND tw.worker_id = $2 AND tw.actual`,
		trackID, workerID, models.MessageTypeApplicant, models.MessageTypeReturnReason,
	).Scan(
		&ticketID, &record.TrackID, &record.Status, &record.Priority, &record.CreatedAt,
		&record.CategoryID, &record.Category, &record.ApplicantType, &record.ReturnCount,
		&record.Description, &record.ReturnReason,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return ExpertTicketDetailRecord{}, ErrExpertAccessDenied
	}
	if err != nil {
		return ExpertTicketDetailRecord{}, fmt.Errorf("get expert ticket: %w", err)
	}
	if err := repository.loadTicketRelations(ctx, ticketID, &record); err != nil {
		return ExpertTicketDetailRecord{}, err
	}
	return record, nil
}

func (repository *ExpertRepository) loadTicketRelations(ctx context.Context, ticketID int64, record *ExpertTicketDetailRecord) error {
	record.Clarifications = make([]ExpertClarificationRecord, 0)
	rows, err := repository.db.QueryContext(ctx, `
		SELECT q.id, q.text, a.id, a.text
		FROM "QA" AS qa
		JOIN questions AS q ON q.id = qa.question_id
		JOIN answers AS a ON a.id = qa.answer_id AND a.question_id = qa.question_id
		WHERE qa.ticket_id = $1
		ORDER BY q.id, a.id`, ticketID)
	if err != nil {
		return fmt.Errorf("list expert clarifications: %w", err)
	}
	for rows.Next() {
		var item ExpertClarificationRecord
		if err := rows.Scan(&item.QuestionID, &item.Question, &item.AnswerID, &item.Answer); err != nil {
			rows.Close()
			return fmt.Errorf("scan expert clarification: %w", err)
		}
		record.Clarifications = append(record.Clarifications, item)
	}
	if err := rows.Close(); err != nil {
		return fmt.Errorf("close expert clarifications: %w", err)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate expert clarifications: %w", err)
	}

	record.Attachments = make([]ChatAttachmentRecord, 0)
	attachmentRows, err := repository.db.QueryContext(ctx, `
		SELECT id, safe_name, mime_type, size_bytes, created_at
		FROM attachments WHERE ticket_id = $1 ORDER BY created_at, id`, ticketID)
	if err != nil {
		return fmt.Errorf("list expert attachments: %w", err)
	}
	for attachmentRows.Next() {
		var item ChatAttachmentRecord
		if err := attachmentRows.Scan(&item.ID, &item.SafeName, &item.MIMEType, &item.SizeBytes, &item.CreatedAt); err != nil {
			attachmentRows.Close()
			return fmt.Errorf("scan expert attachment: %w", err)
		}
		record.Attachments = append(record.Attachments, item)
	}
	if err := attachmentRows.Close(); err != nil {
		return fmt.Errorf("close expert attachments: %w", err)
	}
	if err := attachmentRows.Err(); err != nil {
		return fmt.Errorf("iterate expert attachments: %w", err)
	}

	record.Messages = make([]ChatMessageRecord, 0)
	record.Notes = make([]ExpertNoteRecord, 0)
	messageRows, err := repository.db.QueryContext(ctx, `
		SELECT id, text, type, created_at FROM messages
		WHERE ticket_id = $1 AND type IN ($2, $3, $4)
		ORDER BY created_at, id`, ticketID, models.MessageTypeApplicant, models.MessageTypeSpecialist, models.MessageTypeInternalNote)
	if err != nil {
		return fmt.Errorf("list expert messages: %w", err)
	}
	for messageRows.Next() {
		var item ChatMessageRecord
		if err := messageRows.Scan(&item.ID, &item.Text, &item.Type, &item.CreatedAt); err != nil {
			messageRows.Close()
			return fmt.Errorf("scan expert message: %w", err)
		}
		if item.Type == models.MessageTypeInternalNote {
			var payload expertNotePayload
			if json.Unmarshal([]byte(item.Text), &payload) == nil && payload.Kind == "expert_note" {
				record.Notes = append(record.Notes, ExpertNoteRecord{
					ID: item.ID, Text: payload.Text, AuthorID: payload.AuthorID,
					AuthorName: payload.AuthorName, CreatedAt: item.CreatedAt,
				})
			}
			continue
		}
		record.Messages = append(record.Messages, item)
	}
	if err := messageRows.Close(); err != nil {
		return fmt.Errorf("close expert messages: %w", err)
	}
	if err := messageRows.Err(); err != nil {
		return fmt.Errorf("iterate expert messages: %w", err)
	}

	record.Workers = make([]ExpertWorkerRecord, 0)
	workerRows, err := repository.db.QueryContext(ctx, `
		SELECT u.id, u.full_name, eg.title, tw.is_responsible
		FROM tickets_workers AS tw
		JOIN users AS u ON u.id = tw.worker_id
		JOIN expert_groups AS eg ON eg.id = u.expert_group_id
		WHERE tw.ticket_id = $1 AND tw.actual
		ORDER BY tw.is_responsible DESC, u.full_name`, ticketID)
	if err != nil {
		return fmt.Errorf("list expert ticket workers: %w", err)
	}
	for workerRows.Next() {
		var item ExpertWorkerRecord
		if err := workerRows.Scan(&item.ID, &item.FullName, &item.GroupTitle, &item.IsResponsible); err != nil {
			workerRows.Close()
			return fmt.Errorf("scan expert ticket worker: %w", err)
		}
		record.Workers = append(record.Workers, item)
	}
	if err := workerRows.Close(); err != nil {
		return fmt.Errorf("close expert ticket workers: %w", err)
	}
	return workerRows.Err()
}

type expertNotePayload struct {
	Kind       string `json:"kind"`
	Version    int    `json:"version"`
	Text       string `json:"text"`
	AuthorID   int64  `json:"author_id"`
	AuthorName string `json:"author_name"`
}

func (repository *ExpertRepository) AddNote(ctx context.Context, trackID string, workerID int64, authorName, text string, createdAt time.Time) (ExpertNoteRecord, error) {
	payload := expertNotePayload{Kind: "expert_note", Version: 1, Text: text, AuthorID: workerID, AuthorName: authorName}
	encoded, err := json.Marshal(payload)
	if err != nil {
		return ExpertNoteRecord{}, fmt.Errorf("encode expert note: %w", err)
	}
	tx, err := repository.db.BeginTx(ctx, nil)
	if err != nil {
		return ExpertNoteRecord{}, fmt.Errorf("begin expert note transaction: %w", err)
	}
	defer tx.Rollback()
	ticketID, status, _, err := lockExpertTicket(ctx, tx, trackID, workerID)
	if err != nil {
		return ExpertNoteRecord{}, err
	}
	if !isExpertActiveStatus(status) {
		return ExpertNoteRecord{}, ErrInvalidTransition
	}
	var id int64
	if err := tx.QueryRowContext(ctx, `
		INSERT INTO messages (ticket_id, text, type, created_at)
		VALUES ($1, $2, $3, $4) RETURNING id`,
		ticketID, string(encoded), models.MessageTypeInternalNote, createdAt,
	).Scan(&id); err != nil {
		return ExpertNoteRecord{}, fmt.Errorf("save expert note: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return ExpertNoteRecord{}, fmt.Errorf("commit expert note: %w", err)
	}
	return ExpertNoteRecord{ID: id, Text: text, AuthorID: workerID, AuthorName: authorName, CreatedAt: createdAt}, nil
}

func (repository *ExpertRepository) AddAnswer(ctx context.Context, trackID string, workerID int64, text string, createdAt time.Time) (ExpertMessageRecord, error) {
	tx, err := repository.db.BeginTx(ctx, nil)
	if err != nil {
		return ExpertMessageRecord{}, fmt.Errorf("begin expert answer transaction: %w", err)
	}
	defer tx.Rollback()
	ticketID, status, _, err := lockExpertTicket(ctx, tx, trackID, workerID)
	if err != nil {
		return ExpertMessageRecord{}, err
	}
	switch status {
	case models.TicketStatusAssigned, models.TicketStatusInProgress, models.TicketStatusNeedsClarification:
	default:
		return ExpertMessageRecord{}, ErrInvalidTransition
	}
	var id int64
	if err := tx.QueryRowContext(ctx, `
		INSERT INTO messages (ticket_id, text, type, created_at)
		VALUES ($1, $2, $3, $4) RETURNING id`,
		ticketID, text, models.MessageTypeSpecialist, createdAt,
	).Scan(&id); err != nil {
		return ExpertMessageRecord{}, fmt.Errorf("save expert answer: %w", err)
	}
	if _, err := tx.ExecContext(ctx, `UPDATE tickets SET status = $2 WHERE id = $1`, ticketID, models.TicketStatusAnswerReady); err != nil {
		return ExpertMessageRecord{}, fmt.Errorf("set expert answer status: %w", err)
	}
	if err := insertEvent(ctx, tx, ticketID, workerID, "specialist_answered", status, models.TicketStatusAnswerReady, sql.NullInt64{}, sql.NullInt64{}, createdAt); err != nil {
		return ExpertMessageRecord{}, err
	}
	if err := tx.Commit(); err != nil {
		return ExpertMessageRecord{}, fmt.Errorf("commit expert answer: %w", err)
	}
	return ExpertMessageRecord{ID: id, Text: text, Type: models.MessageTypeSpecialist, CreatedAt: createdAt, TicketStatus: models.TicketStatusAnswerReady}, nil
}

func (repository *ExpertRepository) CreateWorkerRequest(ctx context.Context, trackID string, workerID int64, requestType, reason string, createdAt time.Time) (WorkerRequestRecord, error) {
	tx, err := repository.db.BeginTx(ctx, nil)
	if err != nil {
		return WorkerRequestRecord{}, fmt.Errorf("begin worker request transaction: %w", err)
	}
	defer tx.Rollback()
	ticketID, status, _, err := lockExpertTicket(ctx, tx, trackID, workerID)
	if err != nil {
		return WorkerRequestRecord{}, err
	}
	if !isExpertActiveStatus(status) {
		return WorkerRequestRecord{}, ErrInvalidTransition
	}
	payload := WorkerRequestPayload{Kind: "worker_request", Version: 1, RequestType: requestType, Status: "sent", Reason: reason, CreatedBy: workerID}
	encoded, err := json.Marshal(payload)
	if err != nil {
		return WorkerRequestRecord{}, fmt.Errorf("encode worker request: %w", err)
	}
	var id int64
	if err := tx.QueryRowContext(ctx, `
		INSERT INTO messages (ticket_id, text, type, created_at)
		VALUES ($1, $2, $3, $4) RETURNING id`,
		ticketID, string(encoded), models.MessageTypeInternalNote, createdAt,
	).Scan(&id); err != nil {
		return WorkerRequestRecord{}, fmt.Errorf("save worker request: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return WorkerRequestRecord{}, fmt.Errorf("commit worker request: %w", err)
	}
	return WorkerRequestRecord{ID: id, TicketID: ticketID, TrackID: trackID, RequestType: requestType, Reason: reason, Status: "sent", CreatedAt: createdAt, CreatedBy: workerID}, nil
}

func (repository *ExpertRepository) ListWorkerRequests(ctx context.Context, workerID int64, status string, limit, offset int) (WorkerRequestPageRecord, error) {
	rows, err := repository.db.QueryContext(ctx, `
		SELECT m.id, m.ticket_id, t.track_id, c.id, c.name, t.morda_type,
		       t.priority, m.text, m.created_at
		FROM messages AS m
		JOIN tickets AS t ON t.id = m.ticket_id
		JOIN categories AS c ON c.id = t.category_id
		WHERE m.type = $1 AND m.text LIKE '%"kind":"worker_request"%'
		ORDER BY m.created_at DESC, m.id DESC`, models.MessageTypeInternalNote)
	if err != nil {
		return WorkerRequestPageRecord{}, fmt.Errorf("list worker request messages: %w", err)
	}
	defer rows.Close()
	all := make([]WorkerRequestRecord, 0)
	for rows.Next() {
		var item WorkerRequestRecord
		var encoded string
		if err := rows.Scan(
			&item.ID, &item.TicketID, &item.TrackID, &item.CategoryID, &item.Category,
			&item.ApplicantType, &item.Priority, &encoded, &item.CreatedAt,
		); err != nil {
			return WorkerRequestPageRecord{}, fmt.Errorf("scan worker request message: %w", err)
		}
		var payload WorkerRequestPayload
		if json.Unmarshal([]byte(encoded), &payload) != nil || payload.Kind != "worker_request" || payload.CreatedBy != workerID {
			continue
		}
		item.RequestType, item.Reason = payload.RequestType, payload.Reason
		item.CreatedBy = payload.CreatedBy
		item.Status, err = repository.workerRequestStatus(ctx, item.TicketID, item.ID, payload.Status)
		if err != nil {
			return WorkerRequestPageRecord{}, err
		}
		if status == "" || item.Status == status {
			all = append(all, item)
		}
	}
	if err := rows.Err(); err != nil {
		return WorkerRequestPageRecord{}, fmt.Errorf("iterate worker request messages: %w", err)
	}
	total := len(all)
	if offset >= total {
		return WorkerRequestPageRecord{Items: make([]WorkerRequestRecord, 0), Total: total}, nil
	}
	end := offset + limit
	if end > total {
		end = total
	}
	return WorkerRequestPageRecord{Items: all[offset:end], Total: total}, nil
}

func (repository *ExpertRepository) ListOperatorWorkerRequests(ctx context.Context, status string, limit, offset int) (WorkerRequestPageRecord, error) {
	rows, err := repository.db.QueryContext(ctx, `
		SELECT m.id, m.ticket_id, t.track_id, c.id, c.name, t.morda_type,
		       t.priority, m.text, m.created_at
		FROM messages AS m
		JOIN tickets AS t ON t.id = m.ticket_id
		JOIN categories AS c ON c.id = t.category_id
		WHERE m.type = $1 AND m.text LIKE '%"kind":"worker_request"%'
		ORDER BY m.created_at ASC, m.id ASC`, models.MessageTypeInternalNote)
	if err != nil {
		return WorkerRequestPageRecord{}, fmt.Errorf("list operator worker requests: %w", err)
	}
	defer rows.Close()
	all := make([]WorkerRequestRecord, 0)
	for rows.Next() {
		var item WorkerRequestRecord
		var encoded string
		if err := rows.Scan(&item.ID, &item.TicketID, &item.TrackID, &item.CategoryID, &item.Category, &item.ApplicantType, &item.Priority, &encoded, &item.CreatedAt); err != nil {
			return WorkerRequestPageRecord{}, fmt.Errorf("scan operator worker request: %w", err)
		}
		var payload WorkerRequestPayload
		if json.Unmarshal([]byte(encoded), &payload) != nil || payload.Kind != "worker_request" {
			continue
		}
		item.RequestType, item.Reason, item.CreatedBy = payload.RequestType, payload.Reason, payload.CreatedBy
		item.Status, err = repository.workerRequestStatus(ctx, item.TicketID, item.ID, payload.Status)
		if err != nil {
			return WorkerRequestPageRecord{}, err
		}
		if status == "" || item.Status == status {
			all = append(all, item)
		}
	}
	if err := rows.Err(); err != nil {
		return WorkerRequestPageRecord{}, fmt.Errorf("iterate operator worker requests: %w", err)
	}
	total := len(all)
	if offset >= total {
		return WorkerRequestPageRecord{Items: []WorkerRequestRecord{}, Total: total}, nil
	}
	end := offset + limit
	if end > total {
		end = total
	}
	return WorkerRequestPageRecord{Items: all[offset:end], Total: total}, nil
}

func (repository *ExpertRepository) workerRequestStatus(ctx context.Context, ticketID, requestID int64, fallback string) (string, error) {
	rows, err := repository.db.QueryContext(ctx, `
		SELECT text FROM messages
		WHERE ticket_id = $1 AND type = $2 AND text LIKE '%"kind":"worker_request_status"%'
		ORDER BY created_at DESC, id DESC`, ticketID, models.MessageTypeInternalNote)
	if err != nil {
		return "", fmt.Errorf("list worker request statuses: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var encoded string
		if err := rows.Scan(&encoded); err != nil {
			return "", fmt.Errorf("scan worker request status: %w", err)
		}
		var payload WorkerRequestStatusPayload
		if json.Unmarshal([]byte(encoded), &payload) == nil && payload.Kind == "worker_request_status" && payload.RequestID == requestID {
			return payload.Status, nil
		}
	}
	if err := rows.Err(); err != nil {
		return "", fmt.Errorf("iterate worker request statuses: %w", err)
	}
	return fallback, nil
}

func (repository *ExpertRepository) CanAccessTicket(ctx context.Context, trackID string, workerID int64) error {
	var allowed bool
	if err := repository.db.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM tickets AS t
			JOIN tickets_workers AS tw ON tw.ticket_id = t.id
			WHERE t.track_id = $1 AND tw.worker_id = $2 AND tw.actual
		)`, trackID, workerID).Scan(&allowed); err != nil {
		return fmt.Errorf("check expert ticket access: %w", err)
	}
	if !allowed {
		return ErrExpertAccessDenied
	}
	return nil
}

func lockExpertTicket(ctx context.Context, tx *sql.Tx, trackID string, workerID int64) (int64, models.TicketStatus, int, error) {
	var ticketID int64
	var status models.TicketStatus
	var returnCount int
	err := tx.QueryRowContext(ctx, `
		SELECT t.id, t.status, t.return_count
		FROM tickets AS t
		JOIN tickets_workers AS tw ON tw.ticket_id = t.id
		WHERE t.track_id = $1 AND tw.worker_id = $2 AND tw.actual
		FOR UPDATE OF t`, trackID, workerID).Scan(&ticketID, &status, &returnCount)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, 0, 0, ErrExpertAccessDenied
	}
	if err != nil {
		return 0, 0, 0, fmt.Errorf("lock expert ticket: %w", err)
	}
	return ticketID, status, returnCount, nil
}

func isExpertActiveStatus(status models.TicketStatus) bool {
	switch status {
	case models.TicketStatusAssigned, models.TicketStatusInProgress,
		models.TicketStatusNeedsClarification, models.TicketStatusAnswerReady:
		return true
	default:
		return false
	}
}

func (repository *ExpertRepository) Analytics(ctx context.Context, workerID int64, start, end time.Time) (AnalyticsRecord, error) {
	return repository.analytics(ctx, workerID, start, end)
}

func (repository *ExpertRepository) analytics(ctx context.Context, workerID int64, start, end time.Time) (AnalyticsRecord, error) {
	result := AnalyticsRecord{Categories: []NamedCountRecord{}, ApplicantTypes: []NumberedCountRecord{}, Statuses: []NumberedCountRecord{}}
	base := `EXISTS (SELECT 1 FROM tickets_workers tw WHERE tw.ticket_id = t.id AND tw.worker_id = $3)`
	if err := repository.db.QueryRowContext(ctx, `
		SELECT COUNT(*), COUNT(*) FILTER (WHERE t.priority = $4), COUNT(*) FILTER (WHERE t.return_count > 0)
		FROM tickets t WHERE t.created_at >= $1 AND t.created_at < $2 AND `+base,
		start, end, workerID, models.TicketPriorityUrgent).Scan(&result.Total, &result.UrgentCount, &result.ReturnedCount); err != nil {
		return AnalyticsRecord{}, fmt.Errorf("get expert analytics totals: %w", err)
	}
	categoryRows, err := repository.db.QueryContext(ctx, `
		SELECT c.id, c.name, COUNT(*) FROM tickets t
		JOIN categories c ON c.id = t.category_id
		WHERE t.created_at >= $1 AND t.created_at < $2 AND `+base+`
		GROUP BY c.id, c.name ORDER BY COUNT(*) DESC, c.name`, start, end, workerID)
	if err != nil {
		return AnalyticsRecord{}, fmt.Errorf("get expert category analytics: %w", err)
	}
	for categoryRows.Next() {
		var item NamedCountRecord
		if err := categoryRows.Scan(&item.ID, &item.Name, &item.Count); err != nil {
			categoryRows.Close()
			return AnalyticsRecord{}, fmt.Errorf("scan expert category analytics: %w", err)
		}
		result.Categories = append(result.Categories, item)
	}
	categoryRows.Close()
	result.ApplicantTypes, err = repository.expertNumberedCounts(ctx, `SELECT t.morda_type, COUNT(*) FROM tickets t WHERE t.created_at >= $1 AND t.created_at < $2 AND `+base+` GROUP BY t.morda_type ORDER BY t.morda_type`, start, end, workerID)
	if err != nil {
		return AnalyticsRecord{}, fmt.Errorf("get expert applicant analytics: %w", err)
	}
	result.Statuses, err = repository.expertNumberedCounts(ctx, `SELECT t.status, COUNT(*) FROM tickets t WHERE t.created_at >= $1 AND t.created_at < $2 AND `+base+` GROUP BY t.status ORDER BY t.status`, start, end, workerID)
	if err != nil {
		return AnalyticsRecord{}, fmt.Errorf("get expert status analytics: %w", err)
	}
	var accept, response, closeDuration sql.NullFloat64
	if err := repository.db.QueryRowContext(ctx, `
		SELECT
			AVG(EXTRACT(EPOCH FROM (assigned.created_at - t.created_at))) FILTER (WHERE assigned.created_at >= t.created_at),
			AVG(EXTRACT(EPOCH FROM (answered.created_at - t.created_at))) FILTER (WHERE answered.created_at >= t.created_at),
			AVG(EXTRACT(EPOCH FROM (t.closed_at - t.created_at))) FILTER (WHERE t.closed_at >= t.created_at)
		FROM tickets t
		LEFT JOIN LATERAL (SELECT MIN(tw.created_at) created_at FROM tickets_workers tw WHERE tw.ticket_id = t.id AND tw.worker_id = $3) assigned ON TRUE
		LEFT JOIN LATERAL (SELECT MIN(e.created_at) created_at FROM ticket_events e WHERE e.ticket_id = t.id AND e.actor_user_id = $3 AND e.event_type = 'specialist_answered') answered ON TRUE
		WHERE t.created_at >= $1 AND t.created_at < $2 AND `+base, start, end, workerID).Scan(&accept, &response, &closeDuration); err != nil {
		return AnalyticsRecord{}, fmt.Errorf("get expert duration analytics: %w", err)
	}
	result.DurationSecondsToAccept = nullableFloat(accept)
	result.DurationSecondsToResponse = nullableFloat(response)
	result.DurationSecondsToClose = nullableFloat(closeDuration)
	var active, maximum int
	if err := repository.db.QueryRowContext(ctx, `
		SELECT COUNT(tw.ticket_id) FILTER (WHERE tw.actual AND t.status IN ($2, $3, $4, $5)), u.max_tickets
		FROM users u LEFT JOIN tickets_workers tw ON tw.worker_id = u.id LEFT JOIN tickets t ON t.id = tw.ticket_id
		WHERE u.id = $1 GROUP BY u.max_tickets`, workerID, models.TicketStatusAssigned, models.TicketStatusInProgress, models.TicketStatusNeedsClarification, models.TicketStatusAnswerReady).Scan(&active, &maximum); err != nil {
		return AnalyticsRecord{}, fmt.Errorf("get expert load analytics: %w", err)
	}
	if maximum > 0 {
		result.ExpertLoadPercent = 100 * float64(active) / float64(maximum)
	}
	return result, nil
}

func (repository *ExpertRepository) expertNumberedCounts(ctx context.Context, query string, start, end time.Time, workerID int64) ([]NumberedCountRecord, error) {
	rows, err := repository.db.QueryContext(ctx, query, start, end, workerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []NumberedCountRecord{}
	for rows.Next() {
		var item NumberedCountRecord
		if err := rows.Scan(&item.Value, &item.Count); err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, rows.Err()
}

func (repository *ExpertRepository) ReportTickets(ctx context.Context, workerID int64, start, end time.Time) ([]ReportTicketRecord, error) {
	rows, err := repository.db.QueryContext(ctx, `
		SELECT t.created_at, t.closed_at, c.name, t.morda_type, t.status, t.priority, t.return_count,
		       EXTRACT(EPOCH FROM (assigned.created_at - t.created_at)),
		       EXTRACT(EPOCH FROM (answered.created_at - t.created_at)),
		       EXTRACT(EPOCH FROM (t.closed_at - t.created_at))
		FROM tickets t JOIN categories c ON c.id = t.category_id
		LEFT JOIN LATERAL (SELECT MIN(tw.created_at) created_at FROM tickets_workers tw WHERE tw.ticket_id = t.id AND tw.worker_id = $3) assigned ON TRUE
		LEFT JOIN LATERAL (SELECT MIN(e.created_at) created_at FROM ticket_events e WHERE e.ticket_id = t.id AND e.actor_user_id = $3 AND e.event_type = 'specialist_answered') answered ON TRUE
		WHERE t.created_at >= $1 AND t.created_at < $2
		  AND EXISTS (SELECT 1 FROM tickets_workers tw WHERE tw.ticket_id = t.id AND tw.worker_id = $3)
		ORDER BY t.created_at, t.id`, start, end, workerID)
	if err != nil {
		return nil, fmt.Errorf("get expert report tickets: %w", err)
	}
	defer rows.Close()
	result := []ReportTicketRecord{}
	for rows.Next() {
		var item ReportTicketRecord
		if err := rows.Scan(&item.CreatedAt, &item.ClosedAt, &item.Category, &item.ApplicantType, &item.Status, &item.Priority, &item.ReturnCount, &item.SecondsToAccept, &item.SecondsToFirstResponse, &item.SecondsToClose); err != nil {
			return nil, fmt.Errorf("scan expert report ticket: %w", err)
		}
		result = append(result, item)
	}
	return result, rows.Err()
}
