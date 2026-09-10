package repos

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"example.com/german/backend/internal/models"
)

type OperatorRepository struct {
	db *sql.DB
}

type RecommendedGroupRecord struct {
	ID    int64
	Title string
}

type EligibleWorkerRecord struct {
	ID            int64
	FullName      string
	GroupID       int64
	GroupTitle    string
	ActiveTickets int
	MaxTickets    int
	Recommended   bool
}

type EligibleWorkersRecord struct {
	RecommendedGroup *RecommendedGroupRecord
	Workers          []EligibleWorkerRecord
}

type AssignmentRecord struct {
	WorkerID int64
	Status   models.TicketStatus
}

type RejectionRecord struct {
	Status   models.TicketStatus
	ClosedAt time.Time
	Message  string
}

type OperatorTicketFilter struct {
	Queue         string
	Search        string
	Priority      models.TicketPriority
	Status        models.TicketStatus
	CategoryID    int64
	ApplicantType models.ApplicantType
	Limit         int
	Offset        int
}

type OperatorTicketRecord struct {
	TrackID            string
	CategoryID         int64
	Category           string
	Status             models.TicketStatus
	ApplicantType      models.ApplicantType
	Priority           models.TicketPriority
	CreatedAt          time.Time
	WaitingSeconds     int64
	ResponsibleID      sql.NullInt64
	ResponsibleName    sql.NullString
	ReturnCount        int
	ReturnReason       sql.NullString
	ReturnedAt         sql.NullTime
	PreviousWorkerID   sql.NullInt64
	PreviousWorkerName sql.NullString
}

type OperatorTicketPageRecord struct {
	Items []OperatorTicketRecord
	Total int
}

type NamedCountRecord struct {
	ID    int64
	Name  string
	Count int
}

type NumberedCountRecord struct {
	Value int
	Count int
}

type AnalyticsRecord struct {
	Total                     int
	Categories                []NamedCountRecord
	ApplicantTypes            []NumberedCountRecord
	Statuses                  []NumberedCountRecord
	DurationSecondsToAccept   *float64
	DurationSecondsToResponse *float64
	DurationSecondsToClose    *float64
	ExpertLoadPercent         float64
	UrgentCount               int
	ReturnedCount             int
}

type ReportTicketRecord struct {
	CreatedAt              time.Time
	ClosedAt               sql.NullTime
	Category               string
	ApplicantType          models.ApplicantType
	Status                 models.TicketStatus
	Priority               models.TicketPriority
	ReturnCount            int
	SecondsToAccept        sql.NullFloat64
	SecondsToFirstResponse sql.NullFloat64
	SecondsToClose         sql.NullFloat64
}

func NewOperatorRepository(db *sql.DB) *OperatorRepository {
	return &OperatorRepository{db: db}
}

func (repository *OperatorRepository) EligibleWorkers(ctx context.Context, trackID, name string, groupID int64) (EligibleWorkersRecord, error) {
	var categoryID int64
	if err := repository.db.QueryRowContext(ctx, `SELECT category_id FROM tickets WHERE track_id = $1`, trackID).Scan(&categoryID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return EligibleWorkersRecord{}, ErrTicketNotFound
		}
		return EligibleWorkersRecord{}, fmt.Errorf("get ticket for eligible workers: %w", err)
	}

	result := EligibleWorkersRecord{Workers: make([]EligibleWorkerRecord, 0)}
	var recommended RecommendedGroupRecord
	err := repository.db.QueryRowContext(ctx, `
		SELECT eg.id, eg.title
		FROM cats_expert_groups AS ceg
		JOIN expert_groups AS eg ON eg.id = ceg.group_id
		WHERE ceg.cat_id = $1
		ORDER BY eg.id
		LIMIT 1`, categoryID).Scan(&recommended.ID, &recommended.Title)
	if err == nil {
		result.RecommendedGroup = &recommended
	} else if !errors.Is(err, sql.ErrNoRows) {
		return EligibleWorkersRecord{}, fmt.Errorf("get recommended expert group: %w", err)
	}

	rows, err := repository.db.QueryContext(ctx, `
		SELECT u.id, u.full_name, eg.id, eg.title,
		       COUNT(tw.ticket_id) FILTER (
		           WHERE tw.actual AND active_ticket.status IN ($4, $5, $6, $7)
		       ) AS active_tickets,
		       u.max_tickets,
		       EXISTS (
		           SELECT 1 FROM cats_expert_groups AS route
		           WHERE route.cat_id = $3 AND route.group_id = u.expert_group_id
		       ) AS recommended
		FROM users AS u
		JOIN expert_groups AS eg ON eg.id = u.expert_group_id
		LEFT JOIN tickets_workers AS tw ON tw.worker_id = u.id
		LEFT JOIN tickets AS active_ticket ON active_ticket.id = tw.ticket_id
		WHERE u.role = 'expert'
		  AND ($1 = '' OR u.full_name ILIKE '%%' || $1 || '%%')
		  AND ($2 = 0 OR u.expert_group_id = $2)
		GROUP BY u.id, u.full_name, eg.id, eg.title, u.max_tickets, u.expert_group_id
		HAVING COUNT(tw.ticket_id) FILTER (
		           WHERE tw.actual AND active_ticket.status IN ($4, $5, $6, $7)
		       ) < u.max_tickets
		ORDER BY recommended DESC, active_tickets ASC, u.full_name ASC`,
		name,
		groupID,
		categoryID,
		models.TicketStatusAssigned,
		models.TicketStatusInProgress,
		models.TicketStatusNeedsClarification,
		models.TicketStatusAnswerReady,
	)
	if err != nil {
		return EligibleWorkersRecord{}, fmt.Errorf("list eligible workers: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var worker EligibleWorkerRecord
		if err := rows.Scan(
			&worker.ID,
			&worker.FullName,
			&worker.GroupID,
			&worker.GroupTitle,
			&worker.ActiveTickets,
			&worker.MaxTickets,
			&worker.Recommended,
		); err != nil {
			return EligibleWorkersRecord{}, fmt.Errorf("scan eligible worker: %w", err)
		}
		result.Workers = append(result.Workers, worker)
	}
	if err := rows.Err(); err != nil {
		return EligibleWorkersRecord{}, fmt.Errorf("iterate eligible workers: %w", err)
	}
	return result, nil
}

func (repository *OperatorRepository) SetResponsibleWorker(ctx context.Context, trackID string, workerID, actorID int64, changedAt time.Time) (AssignmentRecord, error) {
	tx, err := repository.db.BeginTx(ctx, nil)
	if err != nil {
		return AssignmentRecord{}, fmt.Errorf("begin responsible assignment transaction: %w", err)
	}
	defer tx.Rollback()

	ticketID, status, err := lockAssignableTicket(ctx, tx, trackID)
	if err != nil {
		return AssignmentRecord{}, err
	}
	if err := lockAvailableExpert(ctx, tx, workerID, ticketID); err != nil {
		return AssignmentRecord{}, err
	}

	var previousWorkerID sql.NullInt64
	err = tx.QueryRowContext(ctx, `
		SELECT worker_id FROM tickets_workers
		WHERE ticket_id = $1 AND actual AND is_responsible
		FOR UPDATE`, ticketID).Scan(&previousWorkerID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return AssignmentRecord{}, fmt.Errorf("get previous responsible worker: %w", err)
	}
	if previousWorkerID.Valid && previousWorkerID.Int64 != workerID {
		if _, err := tx.ExecContext(ctx, `
			UPDATE tickets_workers SET actual = FALSE
			WHERE ticket_id = $1 AND actual AND is_responsible`, ticketID); err != nil {
			return AssignmentRecord{}, fmt.Errorf("deactivate previous responsible worker: %w", err)
		}
	}
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO tickets_workers (worker_id, ticket_id, is_responsible, created_at, actual)
		VALUES ($1, $2, TRUE, $3, TRUE)
		ON CONFLICT (worker_id, ticket_id) DO UPDATE
		SET is_responsible = TRUE, actual = TRUE,
		    created_at = LEAST(tickets_workers.created_at, EXCLUDED.created_at)`,
		workerID, ticketID, changedAt); err != nil {
		return AssignmentRecord{}, fmt.Errorf("save responsible worker: %w", err)
	}
	if _, err := tx.ExecContext(ctx, `UPDATE tickets SET status = $2, closed_at = NULL WHERE id = $1`, ticketID, models.TicketStatusAssigned); err != nil {
		return AssignmentRecord{}, fmt.Errorf("set assigned ticket status: %w", err)
	}
	if err := insertEvent(ctx, tx, ticketID, actorID, "responsible_assigned", status, models.TicketStatusAssigned, previousWorkerID, sql.NullInt64{Int64: workerID, Valid: true}, changedAt); err != nil {
		return AssignmentRecord{}, err
	}
	if err := tx.Commit(); err != nil {
		return AssignmentRecord{}, fmt.Errorf("commit responsible assignment: %w", err)
	}
	return AssignmentRecord{WorkerID: workerID, Status: models.TicketStatusAssigned}, nil
}

func (repository *OperatorRepository) AddWorker(ctx context.Context, trackID string, workerID, actorID int64, changedAt time.Time) (AssignmentRecord, error) {
	tx, err := repository.db.BeginTx(ctx, nil)
	if err != nil {
		return AssignmentRecord{}, fmt.Errorf("begin co-worker assignment transaction: %w", err)
	}
	defer tx.Rollback()
	ticketID, status, err := lockAssignableTicket(ctx, tx, trackID)
	if err != nil {
		return AssignmentRecord{}, err
	}
	var hasResponsible bool
	if err := tx.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM tickets_workers
			WHERE ticket_id = $1 AND actual AND is_responsible
		)`, ticketID).Scan(&hasResponsible); err != nil {
		return AssignmentRecord{}, fmt.Errorf("check responsible worker: %w", err)
	}
	if !hasResponsible {
		return AssignmentRecord{}, ErrResponsibleRequired
	}
	if err := lockAvailableExpert(ctx, tx, workerID, ticketID); err != nil {
		return AssignmentRecord{}, err
	}

	var responsible, actual bool
	err = tx.QueryRowContext(ctx, `
		SELECT is_responsible, actual FROM tickets_workers
		WHERE worker_id = $1 AND ticket_id = $2
		FOR UPDATE`, workerID, ticketID).Scan(&responsible, &actual)
	if err == nil && actual {
		return AssignmentRecord{}, ErrWorkerAlreadyAssigned
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return AssignmentRecord{}, fmt.Errorf("inspect existing co-worker: %w", err)
	}
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO tickets_workers (worker_id, ticket_id, is_responsible, created_at, actual)
		VALUES ($1, $2, FALSE, $3, TRUE)
		ON CONFLICT (worker_id, ticket_id) DO UPDATE
		SET is_responsible = FALSE, actual = TRUE,
		    created_at = LEAST(tickets_workers.created_at, EXCLUDED.created_at)`,
		workerID, ticketID, changedAt); err != nil {
		return AssignmentRecord{}, fmt.Errorf("save co-worker: %w", err)
	}
	if err := insertEvent(ctx, tx, ticketID, actorID, "worker_added", status, status, sql.NullInt64{}, sql.NullInt64{Int64: workerID, Valid: true}, changedAt); err != nil {
		return AssignmentRecord{}, err
	}
	if err := tx.Commit(); err != nil {
		return AssignmentRecord{}, fmt.Errorf("commit co-worker assignment: %w", err)
	}
	return AssignmentRecord{WorkerID: workerID, Status: status}, nil
}

func (repository *OperatorRepository) RemoveWorker(ctx context.Context, trackID string, workerID, actorID int64, changedAt time.Time) (AssignmentRecord, error) {
	tx, err := repository.db.BeginTx(ctx, nil)
	if err != nil {
		return AssignmentRecord{}, fmt.Errorf("begin co-worker removal transaction: %w", err)
	}
	defer tx.Rollback()
	ticketID, status, err := lockAssignableTicket(ctx, tx, trackID)
	if err != nil {
		return AssignmentRecord{}, err
	}
	var responsible bool
	err = tx.QueryRowContext(ctx, `
		SELECT is_responsible FROM tickets_workers
		WHERE worker_id = $1 AND ticket_id = $2 AND actual
		FOR UPDATE`, workerID, ticketID).Scan(&responsible)
	if errors.Is(err, sql.ErrNoRows) {
		return AssignmentRecord{}, ErrWorkerNotAssigned
	}
	if err != nil {
		return AssignmentRecord{}, fmt.Errorf("get co-worker for removal: %w", err)
	}
	if responsible {
		return AssignmentRecord{}, ErrResponsibleWorker
	}
	if _, err := tx.ExecContext(ctx, `
		UPDATE tickets_workers SET actual = FALSE
		WHERE worker_id = $1 AND ticket_id = $2`, workerID, ticketID); err != nil {
		return AssignmentRecord{}, fmt.Errorf("remove co-worker: %w", err)
	}
	if err := insertEvent(ctx, tx, ticketID, actorID, "worker_removed", status, status, sql.NullInt64{Int64: workerID, Valid: true}, sql.NullInt64{}, changedAt); err != nil {

		return AssignmentRecord{}, err
	}
	if err := tx.Commit(); err != nil {
		return AssignmentRecord{}, fmt.Errorf("commit co-worker removal: %w", err)
	}
	return AssignmentRecord{WorkerID: workerID, Status: status}, nil
}

func (repository *OperatorRepository) Reject(ctx context.Context, trackID string, actorID int64, reasonCode, message string, changedAt time.Time) (RejectionRecord, error) {
	tx, err := repository.db.BeginTx(ctx, nil)
	if err != nil {
		return RejectionRecord{}, fmt.Errorf("begin reject ticket transaction: %w", err)
	}
	defer tx.Rollback()
	ticketID, status, err := lockAssignableTicket(ctx, tx, trackID)
	if err != nil {
		return RejectionRecord{}, err
	}
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO messages (ticket_id, text, type, created_at)
		VALUES ($1, $2, $3, $4)`, ticketID, message, models.MessageTypeOperator, changedAt); err != nil {
		return RejectionRecord{}, fmt.Errorf("save operator rejection message: %w", err)
	}
	if _, err := tx.ExecContext(ctx, `
		UPDATE tickets SET status = $2, closed_at = $3 WHERE id = $1`,
		ticketID, models.TicketStatusRejected, changedAt); err != nil {
		return RejectionRecord{}, fmt.Errorf("reject ticket: %w", err)
	}
	if _, err := tx.ExecContext(ctx, `UPDATE tickets_workers SET actual = FALSE WHERE ticket_id = $1 AND actual`, ticketID); err != nil {
		return RejectionRecord{}, fmt.Errorf("deactivate rejected ticket workers: %w", err)
	}
	if err := insertEventWithReason(ctx, tx, ticketID, actorID, "rejected", status, models.TicketStatusRejected, reasonCode, changedAt); err != nil {
		return RejectionRecord{}, err
	}
	if err := tx.Commit(); err != nil {
		return RejectionRecord{}, fmt.Errorf("commit reject ticket: %w", err)
	}
	return RejectionRecord{Status: models.TicketStatusRejected, ClosedAt: changedAt.UTC(), Message: message}, nil
}

func lockAssignableTicket(ctx context.Context, tx *sql.Tx, trackID string) (int64, models.TicketStatus, error) {
	var ticketID int64
	var status models.TicketStatus
	err := tx.QueryRowContext(ctx, `SELECT id, status FROM tickets WHERE track_id = $1 FOR UPDATE`, trackID).Scan(&ticketID, &status)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, 0, ErrTicketNotFound
	}
	if err != nil {
		return 0, 0, fmt.Errorf("lock ticket: %w", err)
	}
	switch status {
	case models.TicketStatusNew, models.TicketStatusAssigned, models.TicketStatusInProgress,
		models.TicketStatusNeedsClarification, models.TicketStatusAnswerReady, models.TicketStatusReturned:
		return ticketID, status, nil
	default:
		return 0, 0, ErrInvalidTransition
	}
}

func lockAvailableExpert(ctx context.Context, tx *sql.Tx, workerID, ticketID int64) error {
	var role string
	var maxTickets int
	err := tx.QueryRowContext(ctx, `SELECT role, max_tickets FROM users WHERE id = $1 FOR UPDATE`, workerID).Scan(&role, &maxTickets)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrWorkerNotFound
	}
	if err != nil {
		return fmt.Errorf("lock expert: %w", err)
	}
	if role != "expert" {
		return ErrWorkerNotFound
	}
	var alreadyActive bool
	if err := tx.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM tickets_workers
			WHERE worker_id = $1 AND ticket_id = $2 AND actual
		)`, workerID, ticketID).Scan(&alreadyActive); err != nil {
		return fmt.Errorf("check current worker assignment: %w", err)
	}
	if alreadyActive {
		return nil
	}
	var activeTickets int
	if err := tx.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM tickets_workers AS tw
		JOIN tickets AS t ON t.id = tw.ticket_id
		WHERE tw.worker_id = $1 AND tw.actual
		  AND t.status IN ($2, $3, $4, $5)`,
		workerID,
		models.TicketStatusAssigned,
		models.TicketStatusInProgress,
		models.TicketStatusNeedsClarification,
		models.TicketStatusAnswerReady,
	).Scan(&activeTickets); err != nil {
		return fmt.Errorf("count active expert tickets: %w", err)
	}
	if activeTickets >= maxTickets {
		return ErrWorkerUnavailable
	}
	return nil
}

func insertEvent(ctx context.Context, tx *sql.Tx, ticketID, actorID int64, eventType string, fromStatus, toStatus models.TicketStatus, fromWorker, toWorker sql.NullInt64, createdAt time.Time) error {
	_, err := tx.ExecContext(ctx, `
		INSERT INTO ticket_events (
			ticket_id, actor_user_id, event_type, from_status, to_status,
			from_worker_id, to_worker_id, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		ticketID, actorID, eventType, fromStatus, toStatus, fromWorker, toWorker, createdAt)
	if err != nil {
		return fmt.Errorf("save ticket event: %w", err)
	}
	return nil
}

func insertEventWithReason(ctx context.Context, tx *sql.Tx, ticketID, actorID int64, eventType string, fromStatus, toStatus models.TicketStatus, reasonCode string, createdAt time.Time) error {
	_, err := tx.ExecContext(ctx, `
		INSERT INTO ticket_events (
			ticket_id, actor_user_id, event_type, from_status, to_status,
			reason_code, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		ticketID, actorID, eventType, fromStatus, toStatus, reasonCode, createdAt)
	if err != nil {
		return fmt.Errorf("save ticket event with reason: %w", err)
	}
	return nil
}

func (repository *OperatorRepository) ListTickets(ctx context.Context, filter OperatorTicketFilter) (OperatorTicketPageRecord, error) {
	countWhere, countArguments := operatorTicketWhere(filter, 1)
	var total int
	if err := repository.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM tickets AS t `+countWhere, countArguments...).Scan(&total); err != nil {
		return OperatorTicketPageRecord{}, fmt.Errorf("count operator tickets: %w", err)
	}

	dataWhere, filterArguments := operatorTicketWhere(filter, 2)
	arguments := []any{models.MessageTypeReturnReason}
	arguments = append(arguments, filterArguments...)
	arguments = append(arguments, filter.Limit, filter.Offset)
	limitPosition := len(arguments) - 1
	offsetPosition := len(arguments)
	orderBy := "t.priority DESC, t.created_at ASC, t.id ASC"
	if filter.Queue == "returned" {
		orderBy = "return_event.created_at ASC NULLS LAST, t.id ASC"
	}
	query := `
		SELECT t.track_id, c.id, c.name, t.status, t.morda_type, t.priority,
		       t.created_at,
		       GREATEST(0, EXTRACT(EPOCH FROM (CURRENT_TIMESTAMP - t.created_at)))::BIGINT,
		       responsible.worker_id, responsible.full_name,
		       t.return_count,
		       return_reason.text, return_event.created_at,
		       return_event.from_worker_id, previous_worker.full_name
		FROM tickets AS t
		JOIN categories AS c ON c.id = t.category_id
		LEFT JOIN LATERAL (
			SELECT u.id AS worker_id, u.full_name
			FROM tickets_workers AS tw
			JOIN users AS u ON u.id = tw.worker_id
			WHERE tw.ticket_id = t.id AND tw.actual AND tw.is_responsible
			LIMIT 1
		) AS responsible ON TRUE
		LEFT JOIN LATERAL (
			SELECT e.created_at, e.from_worker_id
			FROM ticket_events AS e
			WHERE e.ticket_id = t.id AND e.event_type = 'returned'
			ORDER BY e.created_at DESC, e.id DESC
			LIMIT 1
		) AS return_event ON TRUE
		LEFT JOIN users AS previous_worker ON previous_worker.id = return_event.from_worker_id
		LEFT JOIN LATERAL (
			SELECT m.text
			FROM messages AS m
			WHERE m.ticket_id = t.id AND m.type = $1
			ORDER BY m.created_at DESC, m.id DESC
			LIMIT 1
		) AS return_reason ON TRUE
		` + dataWhere + `
		ORDER BY ` + orderBy + `
		LIMIT $` + fmt.Sprint(limitPosition) + ` OFFSET $` + fmt.Sprint(offsetPosition)

	rows, err := repository.db.QueryContext(ctx, query, arguments...)
	if err != nil {
		return OperatorTicketPageRecord{}, fmt.Errorf("list operator tickets: %w", err)
	}
	defer rows.Close()
	result := OperatorTicketPageRecord{Items: make([]OperatorTicketRecord, 0), Total: total}
	for rows.Next() {
		var item OperatorTicketRecord
		if err := rows.Scan(
			&item.TrackID,
			&item.CategoryID,
			&item.Category,
			&item.Status,
			&item.ApplicantType,
			&item.Priority,
			&item.CreatedAt,
			&item.WaitingSeconds,
			&item.ResponsibleID,
			&item.ResponsibleName,
			&item.ReturnCount,
			&item.ReturnReason,
			&item.ReturnedAt,
			&item.PreviousWorkerID,
			&item.PreviousWorkerName,
		); err != nil {
			return OperatorTicketPageRecord{}, fmt.Errorf("scan operator ticket: %w", err)
		}
		result.Items = append(result.Items, item)
	}
	if err := rows.Err(); err != nil {
		return OperatorTicketPageRecord{}, fmt.Errorf("iterate operator tickets: %w", err)
	}
	return result, nil
}

func operatorTicketWhere(filter OperatorTicketFilter, firstPosition int) (string, []any) {
	conditions := make([]string, 0, 6)
	arguments := make([]any, 0, 5)
	add := func(template string, value any) {
		position := firstPosition + len(arguments)
		conditions = append(conditions, fmt.Sprintf(template, position))
		arguments = append(arguments, value)
	}
	switch filter.Queue {
	case "assigned":
		conditions = append(conditions, fmt.Sprintf("t.status IN (%d, %d, %d, %d)",
			models.TicketStatusAssigned,
			models.TicketStatusInProgress,
			models.TicketStatusNeedsClarification,
			models.TicketStatusAnswerReady,
		))
	case "returned":
		conditions = append(conditions, fmt.Sprintf("t.status = %d", models.TicketStatusReturned))
	}
	if filter.Search != "" {
		add("t.track_id ILIKE '%%' || $%d || '%%'", filter.Search)
	}
	if filter.Priority != 0 {
		add("t.priority = $%d", filter.Priority)
	}
	if filter.Status != 0 {
		add("t.status = $%d", filter.Status)
	}
	if filter.CategoryID != 0 {
		add("t.category_id = $%d", filter.CategoryID)
	}
	if filter.ApplicantType != 0 {
		add("t.morda_type = $%d", filter.ApplicantType)
	}
	if len(conditions) == 0 {
		return "", arguments
	}
	return "WHERE " + strings.Join(conditions, " AND "), arguments
}

func (repository *OperatorRepository) Analytics(ctx context.Context, start, end time.Time) (AnalyticsRecord, error) {
	result := AnalyticsRecord{
		Categories:     make([]NamedCountRecord, 0),
		ApplicantTypes: make([]NumberedCountRecord, 0),
		Statuses:       make([]NumberedCountRecord, 0),
	}
	if err := repository.db.QueryRowContext(ctx, `
		SELECT COUNT(*),
		       COUNT(*) FILTER (WHERE priority = $3),
		       COUNT(*) FILTER (WHERE return_count > 0)
		FROM tickets
		WHERE created_at >= $1 AND created_at < $2`,
		start, end, models.TicketPriorityUrgent,
	).Scan(&result.Total, &result.UrgentCount, &result.ReturnedCount); err != nil {
		return AnalyticsRecord{}, fmt.Errorf("get analytics totals: %w", err)
	}

	categoryRows, err := repository.db.QueryContext(ctx, `
		SELECT c.id, c.name, COUNT(*)
		FROM tickets AS t
		JOIN categories AS c ON c.id = t.category_id
		WHERE t.created_at >= $1 AND t.created_at < $2
		GROUP BY c.id, c.name
		ORDER BY COUNT(*) DESC, c.name`, start, end)
	if err != nil {
		return AnalyticsRecord{}, fmt.Errorf("get category analytics: %w", err)
	}
	for categoryRows.Next() {
		var item NamedCountRecord
		if err := categoryRows.Scan(&item.ID, &item.Name, &item.Count); err != nil {
			categoryRows.Close()
			return AnalyticsRecord{}, fmt.Errorf("scan category analytics: %w", err)
		}
		result.Categories = append(result.Categories, item)
	}
	if err := categoryRows.Close(); err != nil {
		return AnalyticsRecord{}, fmt.Errorf("close category analytics: %w", err)
	}
	if err := categoryRows.Err(); err != nil {
		return AnalyticsRecord{}, fmt.Errorf("iterate category analytics: %w", err)
	}

	result.ApplicantTypes, err = repository.numberedCounts(ctx, `
		SELECT morda_type, COUNT(*) FROM tickets
		WHERE created_at >= $1 AND created_at < $2
		GROUP BY morda_type ORDER BY morda_type`, start, end)
	if err != nil {
		return AnalyticsRecord{}, fmt.Errorf("get applicant analytics: %w", err)
	}
	result.Statuses, err = repository.numberedCounts(ctx, `
		SELECT status, COUNT(*) FROM tickets
		WHERE created_at >= $1 AND created_at < $2
		GROUP BY status ORDER BY status`, start, end)
	if err != nil {
		return AnalyticsRecord{}, fmt.Errorf("get status analytics: %w", err)
	}

	var accept, response, closeDuration sql.NullFloat64
	if err := repository.db.QueryRowContext(ctx, `
		SELECT
			AVG(EXTRACT(EPOCH FROM (accepted.created_at - t.created_at)))
				FILTER (WHERE accepted.created_at >= t.created_at),
			AVG(EXTRACT(EPOCH FROM (first_response.created_at - t.created_at)))
				FILTER (WHERE first_response.created_at >= t.created_at),
			AVG(EXTRACT(EPOCH FROM (t.closed_at - t.created_at)))
				FILTER (WHERE t.closed_at >= t.created_at)
		FROM tickets AS t
		LEFT JOIN LATERAL (
			SELECT MIN(tw.created_at) AS created_at
			FROM tickets_workers AS tw
			WHERE tw.ticket_id = t.id
		) AS accepted ON TRUE
		LEFT JOIN LATERAL (
			SELECT m.created_at
			FROM messages AS m
			WHERE m.ticket_id = t.id AND m.type = $3
			ORDER BY m.created_at, m.id LIMIT 1
		) AS first_response ON TRUE
		WHERE t.created_at >= $1 AND t.created_at < $2`,
		start, end, models.MessageTypeSpecialist,
	).Scan(&accept, &response, &closeDuration); err != nil {
		return AnalyticsRecord{}, fmt.Errorf("get duration analytics: %w", err)
	}
	result.DurationSecondsToAccept = nullableFloat(accept)
	result.DurationSecondsToResponse = nullableFloat(response)
	result.DurationSecondsToClose = nullableFloat(closeDuration)

	if err := repository.db.QueryRowContext(ctx, `
		SELECT COALESCE(
			100.0 * SUM(active_tickets) / NULLIF(SUM(max_tickets), 0),
			0
		)
		FROM (
			SELECT u.id, u.max_tickets,
			       COUNT(tw.ticket_id) FILTER (
			           WHERE tw.actual AND t.status IN ($1, $2, $3, $4)
			       ) AS active_tickets
			FROM users AS u
			LEFT JOIN tickets_workers AS tw ON tw.worker_id = u.id
			LEFT JOIN tickets AS t ON t.id = tw.ticket_id
			WHERE u.role = 'expert'
			GROUP BY u.id, u.max_tickets
		) AS expert_load`,
		models.TicketStatusAssigned,
		models.TicketStatusInProgress,
		models.TicketStatusNeedsClarification,
		models.TicketStatusAnswerReady,
	).Scan(&result.ExpertLoadPercent); err != nil {
		return AnalyticsRecord{}, fmt.Errorf("get expert load analytics: %w", err)
	}
	return result, nil
}

func (repository *OperatorRepository) numberedCounts(ctx context.Context, query string, start, end time.Time) ([]NumberedCountRecord, error) {
	rows, err := repository.db.QueryContext(ctx, query, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]NumberedCountRecord, 0)
	for rows.Next() {
		var item NumberedCountRecord
		if err := rows.Scan(&item.Value, &item.Count); err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, rows.Err()
}

func nullableFloat(value sql.NullFloat64) *float64 {
	if !value.Valid {
		return nil
	}
	return &value.Float64
}

func (repository *OperatorRepository) ReportTickets(ctx context.Context, start, end time.Time) ([]ReportTicketRecord, error) {
	rows, err := repository.db.QueryContext(ctx, `
		SELECT t.created_at, t.closed_at, c.name, t.morda_type, t.status,
		       t.priority, t.return_count,
		       EXTRACT(EPOCH FROM (accepted.created_at - t.created_at)),
		       EXTRACT(EPOCH FROM (first_response.created_at - t.created_at)),
		       EXTRACT(EPOCH FROM (t.closed_at - t.created_at))
		FROM tickets AS t
		JOIN categories AS c ON c.id = t.category_id
		LEFT JOIN LATERAL (
			SELECT MIN(tw.created_at) AS created_at
			FROM tickets_workers AS tw
			WHERE tw.ticket_id = t.id
		) AS accepted ON TRUE
		LEFT JOIN LATERAL (
			SELECT m.created_at
			FROM messages AS m
			WHERE m.ticket_id = t.id AND m.type = $3
			ORDER BY m.created_at, m.id LIMIT 1
		) AS first_response ON TRUE
		WHERE t.created_at >= $1 AND t.created_at < $2
		ORDER BY t.created_at, t.id`, start, end, models.MessageTypeSpecialist)
	if err != nil {
		return nil, fmt.Errorf("get report tickets: %w", err)
	}
	defer rows.Close()
	result := make([]ReportTicketRecord, 0)
	for rows.Next() {
		var item ReportTicketRecord
		if err := rows.Scan(
			&item.CreatedAt,
			&item.ClosedAt,
			&item.Category,
			&item.ApplicantType,
			&item.Status,
			&item.Priority,
			&item.ReturnCount,
			&item.SecondsToAccept,
			&item.SecondsToFirstResponse,
			&item.SecondsToClose,
		); err != nil {
			return nil, fmt.Errorf("scan report ticket: %w", err)
		}
		result = append(result, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate report tickets: %w", err)
	}
	return result, nil
}
