package repos

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"example.com/german/backend/internal/models"
)

type OperatorDashboardRecord struct {
	NewCount      int
	AssignedCount int
	ReturnedCount int
	CrisisCount   int
	OverdueCount  int
}

type OperatorTicketEventRecord struct {
	ID         int64
	EventType  string
	ActorID    sql.NullInt64
	ActorName  sql.NullString
	FromStatus sql.NullInt64
	ToStatus   sql.NullInt64
	Reason     sql.NullString
	CreatedAt  time.Time
}

type OperatorTicketDetailRecord struct {
	TrackID           string
	Status            models.TicketStatus
	Priority          models.TicketPriority
	CreatedAt         time.Time
	ClosedAt          sql.NullTime
	CategoryID        int64
	Category          string
	ApplicantType     models.ApplicantType
	Description       string
	ReturnCount       int
	ReturnReason      sql.NullString
	CrisisDetected    bool
	CrisisContact     sql.NullString
	Clarifications    []ExpertClarificationRecord
	Attachments       []ChatAttachmentRecord
	Workers           []ExpertWorkerRecord
	Notes             []ExpertNoteRecord
	RecommendedGroups []RecommendedGroupRecord
	Events            []OperatorTicketEventRecord
}

type OperatorTicketUpdate struct {
	CategoryID *int64
	Priority   *models.TicketPriority
	Status     *models.TicketStatus
	Reason     string
}

type OperatorCloseRecord struct {
	Status   models.TicketStatus
	ClosedAt time.Time
	Message  string
}

func (repository *OperatorRepository) Dashboard(ctx context.Context, newOverdueBefore, responseOverdueBefore time.Time) (OperatorDashboardRecord, error) {
	var result OperatorDashboardRecord
	err := repository.db.QueryRowContext(ctx, `
		SELECT
		  COUNT(*) FILTER (WHERE t.status=$1),
		  COUNT(*) FILTER (WHERE t.status IN ($2,$3,$4,$5)),
		  COUNT(*) FILTER (WHERE t.status=$6),
		  COUNT(*) FILTER (WHERE t.status IN ($1,$2,$3,$4,$5,$6) AND EXISTS (
		    SELECT 1 FROM ticket_events ce WHERE ce.ticket_id=t.id AND ce.event_type='crisis_detected'
		  )),
		  COUNT(*) FILTER (WHERE
		    (t.status=$1 AND t.created_at < $7) OR
		    (t.status IN ($2,$3,$4,$5) AND COALESCE(assigned.created_at,t.created_at) < $8 AND NOT EXISTS (
		      SELECT 1 FROM messages sm WHERE sm.ticket_id=t.id AND sm.type=$9
		    ))
		  )
		FROM tickets t
		LEFT JOIN LATERAL (
		  SELECT MIN(tw.created_at) AS created_at FROM tickets_workers tw
		  WHERE tw.ticket_id=t.id AND tw.actual AND tw.is_responsible
		) assigned ON TRUE`,
		models.TicketStatusNew, models.TicketStatusAssigned, models.TicketStatusInProgress,
		models.TicketStatusNeedsClarification, models.TicketStatusAnswerReady, models.TicketStatusReturned,
		newOverdueBefore, responseOverdueBefore, models.MessageTypeSpecialist,
	).Scan(&result.NewCount, &result.AssignedCount, &result.ReturnedCount, &result.CrisisCount, &result.OverdueCount)
	if err != nil {
		return OperatorDashboardRecord{}, fmt.Errorf("get operator dashboard: %w", err)
	}
	return result, nil
}

func (repository *OperatorRepository) GetOperatorTicket(ctx context.Context, trackID string) (OperatorTicketDetailRecord, error) {
	var ticketID int64
	var result OperatorTicketDetailRecord
	err := repository.db.QueryRowContext(ctx, `
		SELECT t.id,t.track_id,t.status,t.priority,t.created_at,t.closed_at,
		       c.id,c.name,t.morda_type,t.return_count,
		       COALESCE((SELECT m.text FROM messages m WHERE m.ticket_id=t.id AND m.type=$2 ORDER BY m.created_at,m.id LIMIT 1),''),
		       (SELECT m.text FROM messages m WHERE m.ticket_id=t.id AND m.type=$3 ORDER BY m.created_at DESC,m.id DESC LIMIT 1),
		       EXISTS(SELECT 1 FROM ticket_events e WHERE e.ticket_id=t.id AND e.event_type='crisis_detected'),
		       cc.contact
		FROM tickets t
		JOIN categories c ON c.id=t.category_id
		LEFT JOIN crisis_contact cc ON cc.ticket_id=t.id
		WHERE t.track_id=$1`, trackID, models.MessageTypeApplicant, models.MessageTypeReturnReason).Scan(
		&ticketID, &result.TrackID, &result.Status, &result.Priority, &result.CreatedAt, &result.ClosedAt,
		&result.CategoryID, &result.Category, &result.ApplicantType, &result.ReturnCount,
		&result.Description, &result.ReturnReason, &result.CrisisDetected, &result.CrisisContact,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return OperatorTicketDetailRecord{}, ErrTicketNotFound
	}
	if err != nil {
		return OperatorTicketDetailRecord{}, fmt.Errorf("get operator ticket: %w", err)
	}
	if !result.CrisisDetected {
		result.CrisisContact = sql.NullString{}
	}
	if err := repository.loadOperatorTicketRelations(ctx, ticketID, &result); err != nil {
		return OperatorTicketDetailRecord{}, err
	}
	return result, nil
}

func (repository *OperatorRepository) loadOperatorTicketRelations(ctx context.Context, ticketID int64, result *OperatorTicketDetailRecord) error {
	result.Clarifications = make([]ExpertClarificationRecord, 0)
	rows, err := repository.db.QueryContext(ctx, `SELECT q.id,q.text,a.id,a.text FROM "QA" qa JOIN questions q ON q.id=qa.question_id JOIN answers a ON a.id=qa.answer_id AND a.question_id=qa.question_id WHERE qa.ticket_id=$1 ORDER BY q.id,a.id`, ticketID)
	if err != nil {
		return fmt.Errorf("list operator clarifications: %w", err)
	}
	for rows.Next() {
		var item ExpertClarificationRecord
		if err := rows.Scan(&item.QuestionID, &item.Question, &item.AnswerID, &item.Answer); err != nil {
			rows.Close()
			return fmt.Errorf("scan operator clarification: %w", err)
		}
		result.Clarifications = append(result.Clarifications, item)
	}
	if err := rows.Close(); err != nil {
		return err
	}
	if err := rows.Err(); err != nil {
		return err
	}

	result.Attachments = make([]ChatAttachmentRecord, 0)
	rows, err = repository.db.QueryContext(ctx, `
		SELECT a.id,a.safe_name,a.mime_type,a.size_bytes,a.created_at
		FROM attachments a
		WHERE a.ticket_id=$1
		  AND a.created_at <= COALESCE(
		    (SELECT MIN(m.created_at) FROM messages m WHERE m.ticket_id=$1 AND m.type=$2),
		    (SELECT t.created_at FROM tickets t WHERE t.id=$1)
		  )
		ORDER BY a.created_at,a.id`, ticketID, models.MessageTypeApplicant)
	if err != nil {
		return fmt.Errorf("list operator attachments: %w", err)
	}
	for rows.Next() {
		var item ChatAttachmentRecord
		if err := rows.Scan(&item.ID, &item.SafeName, &item.MIMEType, &item.SizeBytes, &item.CreatedAt); err != nil {
			rows.Close()
			return err
		}
		result.Attachments = append(result.Attachments, item)
	}
	if err := rows.Close(); err != nil {
		return err
	}
	if err := rows.Err(); err != nil {
		return err
	}

	result.Workers = make([]ExpertWorkerRecord, 0)
	rows, err = repository.db.QueryContext(ctx, `SELECT u.id,u.full_name,eg.title,tw.is_responsible FROM tickets_workers tw JOIN users u ON u.id=tw.worker_id JOIN expert_groups eg ON eg.id=u.expert_group_id WHERE tw.ticket_id=$1 AND tw.actual ORDER BY tw.is_responsible DESC,u.full_name`, ticketID)
	if err != nil {
		return fmt.Errorf("list operator workers: %w", err)
	}
	for rows.Next() {
		var item ExpertWorkerRecord
		if err := rows.Scan(&item.ID, &item.FullName, &item.GroupTitle, &item.IsResponsible); err != nil {
			rows.Close()
			return err
		}
		result.Workers = append(result.Workers, item)
	}
	if err := rows.Close(); err != nil {
		return err
	}
	if err := rows.Err(); err != nil {
		return err
	}

	result.Notes = make([]ExpertNoteRecord, 0)
	rows, err = repository.db.QueryContext(ctx, `SELECT id,text,created_at FROM messages WHERE ticket_id=$1 AND type=$2 ORDER BY created_at,id`, ticketID, models.MessageTypeInternalNote)
	if err != nil {
		return fmt.Errorf("list operator notes: %w", err)
	}
	for rows.Next() {
		var id int64
		var raw string
		var createdAt time.Time
		if err := rows.Scan(&id, &raw, &createdAt); err != nil {
			rows.Close()
			return err
		}
		var payload expertNotePayload
		if json.Unmarshal([]byte(raw), &payload) == nil && payload.Kind == "expert_note" {
			result.Notes = append(result.Notes, ExpertNoteRecord{ID: id, Text: payload.Text, AuthorID: payload.AuthorID, AuthorName: payload.AuthorName, CreatedAt: createdAt})
		}
	}
	if err := rows.Close(); err != nil {
		return err
	}
	if err := rows.Err(); err != nil {
		return err
	}

	result.RecommendedGroups = make([]RecommendedGroupRecord, 0)
	rows, err = repository.db.QueryContext(ctx, `SELECT eg.id,eg.title FROM cats_expert_groups ceg JOIN expert_groups eg ON eg.id=ceg.group_id JOIN tickets t ON t.category_id=ceg.cat_id WHERE t.id=$1 ORDER BY eg.title,eg.id`, ticketID)
	if err != nil {
		return fmt.Errorf("list operator recommended groups: %w", err)
	}
	for rows.Next() {
		var item RecommendedGroupRecord
		if err := rows.Scan(&item.ID, &item.Title); err != nil {
			rows.Close()
			return err
		}
		result.RecommendedGroups = append(result.RecommendedGroups, item)
	}
	if err := rows.Close(); err != nil {
		return err
	}
	if err := rows.Err(); err != nil {
		return err
	}

	result.Events = make([]OperatorTicketEventRecord, 0)
	rows, err = repository.db.QueryContext(ctx, `SELECT e.id,e.event_type,e.actor_user_id,u.full_name,e.from_status,e.to_status,e.reason_code,e.created_at FROM ticket_events e LEFT JOIN users u ON u.id=e.actor_user_id WHERE e.ticket_id=$1 ORDER BY e.created_at,e.id`, ticketID)
	if err != nil {
		return fmt.Errorf("list operator events: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var item OperatorTicketEventRecord
		if err := rows.Scan(&item.ID, &item.EventType, &item.ActorID, &item.ActorName, &item.FromStatus, &item.ToStatus, &item.Reason, &item.CreatedAt); err != nil {
			return err
		}
		result.Events = append(result.Events, item)
	}
	return rows.Err()
}

func (repository *OperatorRepository) CanAccessOperatorAttachment(ctx context.Context, trackID string, attachmentID int64) error {
	var allowed bool
	err := repository.db.QueryRowContext(ctx, `
		SELECT EXISTS(
		  SELECT 1
		  FROM attachments a
		  JOIN tickets t ON t.id=a.ticket_id
		  WHERE t.track_id=$1 AND a.id=$2
		    AND a.created_at <= COALESCE(
		      (SELECT MIN(m.created_at) FROM messages m WHERE m.ticket_id=t.id AND m.type=$3),
		      t.created_at
		    )
		)`, trackID, attachmentID, models.MessageTypeApplicant).Scan(&allowed)
	if err != nil {
		return fmt.Errorf("check operator attachment access: %w", err)
	}
	if !allowed {
		return ErrAttachmentNotFound
	}
	return nil
}

func (repository *OperatorRepository) UpdateOperatorTicket(ctx context.Context, trackID string, actorID int64, update OperatorTicketUpdate, changedAt time.Time) error {
	tx, err := repository.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin operator ticket update: %w", err)
	}
	defer tx.Rollback()
	var ticketID int64
	var status models.TicketStatus
	var categoryID int64
	var priority models.TicketPriority
	if err = tx.QueryRowContext(ctx, `SELECT id,status,category_id,priority FROM tickets WHERE track_id=$1 FOR UPDATE`, trackID).Scan(&ticketID, &status, &categoryID, &priority); errors.Is(err, sql.ErrNoRows) {
		return ErrTicketNotFound
	} else if err != nil {
		return err
	}
	if (update.CategoryID != nil && *update.CategoryID != categoryID || update.Priority != nil && *update.Priority != priority) && status != models.TicketStatusNew && status != models.TicketStatusReturned {
		return ErrInvalidTransition
	}
	if update.CategoryID != nil && *update.CategoryID != categoryID {
		var exists bool
		if err = tx.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM categories WHERE id=$1)`, *update.CategoryID).Scan(&exists); err != nil {
			return err
		}
		if !exists {
			return ErrAdminResourceNotFound
		}
		if _, err = tx.ExecContext(ctx, `UPDATE tickets SET category_id=$2 WHERE id=$1`, ticketID, *update.CategoryID); err != nil {
			return err
		}
		if err = insertEventWithReason(ctx, tx, ticketID, actorID, "operator_category_changed", status, status, update.Reason, changedAt); err != nil {
			return err
		}
	}
	if update.Priority != nil && *update.Priority != priority {
		if _, err = tx.ExecContext(ctx, `UPDATE tickets SET priority=$2 WHERE id=$1`, ticketID, *update.Priority); err != nil {
			return err
		}
		if err = insertEventWithReason(ctx, tx, ticketID, actorID, "operator_priority_changed", status, status, update.Reason, changedAt); err != nil {
			return err
		}
	}
	if update.Status != nil && *update.Status != status {
		if !operatorTransitionAllowed(status, *update.Status) {
			return ErrInvalidTransition
		}
		if *update.Status == models.TicketStatusAssigned {
			var responsible bool
			if err = tx.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM tickets_workers WHERE ticket_id=$1 AND actual AND is_responsible)`, ticketID).Scan(&responsible); err != nil {
				return err
			}
			if !responsible {
				return ErrResponsibleRequired
			}
		}
		closedAt := any(nil)
		if *update.Status == models.TicketStatusRejected || *update.Status == models.TicketStatusCompleted {
			closedAt = changedAt
		}
		if _, err = tx.ExecContext(ctx, `UPDATE tickets SET status=$2,closed_at=$3 WHERE id=$1`, ticketID, *update.Status, closedAt); err != nil {
			return err
		}
		if *update.Status == models.TicketStatusRejected || *update.Status == models.TicketStatusCompleted {
			if _, err = tx.ExecContext(ctx, `UPDATE tickets_workers SET actual=FALSE WHERE ticket_id=$1 AND actual`, ticketID); err != nil {
				return err
			}
		}
		if *update.Status == models.TicketStatusRejected && update.Reason != "" {
			if _, err = tx.ExecContext(ctx, `INSERT INTO messages(ticket_id,text,type,created_at) VALUES($1,$2,$3,$4)`, ticketID, update.Reason, models.MessageTypeOperator, changedAt); err != nil {
				return err
			}
		}
		if err = insertEventWithReason(ctx, tx, ticketID, actorID, "operator_status_changed", status, *update.Status, update.Reason, changedAt); err != nil {
			return err
		}
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit operator ticket update: %w", err)
	}
	return nil
}

func operatorTransitionAllowed(from, to models.TicketStatus) bool {
	switch from {
	case models.TicketStatusNew:
		return to == models.TicketStatusAssigned || to == models.TicketStatusRejected
	case models.TicketStatusReturned:
		return to == models.TicketStatusAssigned || to == models.TicketStatusRejected || to == models.TicketStatusCompleted
	default:
		return false
	}
}

func (repository *OperatorRepository) CloseOperatorTicket(ctx context.Context, trackID string, actorID int64, message string, closedAt time.Time) (OperatorCloseRecord, error) {
	tx, err := repository.db.BeginTx(ctx, nil)
	if err != nil {
		return OperatorCloseRecord{}, err
	}
	defer tx.Rollback()
	var ticketID int64
	var status models.TicketStatus
	if err = tx.QueryRowContext(ctx, `SELECT id,status FROM tickets WHERE track_id=$1 FOR UPDATE`, trackID).Scan(&ticketID, &status); errors.Is(err, sql.ErrNoRows) {
		return OperatorCloseRecord{}, ErrTicketNotFound
	} else if err != nil {
		return OperatorCloseRecord{}, err
	}
	if status == models.TicketStatusCompleted {
		var existing time.Time
		if err = tx.QueryRowContext(ctx, `SELECT closed_at FROM tickets WHERE id=$1`, ticketID).Scan(&existing); err != nil {
			return OperatorCloseRecord{}, err
		}
		return OperatorCloseRecord{Status: status, ClosedAt: existing, Message: message}, tx.Commit()
	}
	if status != models.TicketStatusNew && status != models.TicketStatusReturned {
		return OperatorCloseRecord{}, ErrInvalidTransition
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO messages(ticket_id,text,type,created_at) VALUES($1,$2,$3,$4)`, ticketID, message, models.MessageTypeOperator, closedAt); err != nil {
		return OperatorCloseRecord{}, err
	}
	if _, err = tx.ExecContext(ctx, `UPDATE tickets SET status=$2,closed_at=$3 WHERE id=$1`, ticketID, models.TicketStatusCompleted, closedAt); err != nil {
		return OperatorCloseRecord{}, err
	}
	if _, err = tx.ExecContext(ctx, `UPDATE tickets_workers SET actual=FALSE WHERE ticket_id=$1 AND actual`, ticketID); err != nil {
		return OperatorCloseRecord{}, err
	}
	if err = insertEvent(ctx, tx, ticketID, actorID, "operator_closed", status, models.TicketStatusCompleted, sql.NullInt64{}, sql.NullInt64{}, closedAt); err != nil {
		return OperatorCloseRecord{}, err
	}
	if err = tx.Commit(); err != nil {
		return OperatorCloseRecord{}, err
	}
	return OperatorCloseRecord{Status: models.TicketStatusCompleted, ClosedAt: closedAt, Message: message}, nil
}
