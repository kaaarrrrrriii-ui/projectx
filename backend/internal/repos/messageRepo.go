package repos

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"example.com/german/backend/internal/models"
)

type MessageRepository struct {
	db *sql.DB
}

type NewAttachmentRecord struct {
	StorageKey        string
	SafeName          string
	MIMEType          string
	OriginalSizeBytes int64
	SizeBytes         int64
	Width             int
	Height            int
}

type CreatedAttachmentRecord struct {
	ID        int64
	SafeName  string
	MIMEType  string
	SizeBytes int64
	CreatedAt time.Time
}

type CreatedMessageRecord struct {
	ID           int64
	Text         string
	Type         models.MessageType
	CreatedAt    time.Time
	Attachments  []CreatedAttachmentRecord
	TicketStatus models.TicketStatus
}

func NewMessageRepository(db *sql.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (repository *MessageRepository) AddApplicantMessage(
	ctx context.Context,
	trackID string,
	text string,
	attachments []NewAttachmentRecord,
	maxTicketAttachments int,
	createdAt time.Time,
) (CreatedMessageRecord, error) {
	tx, err := repository.db.BeginTx(ctx, nil)
	if err != nil {
		return CreatedMessageRecord{}, fmt.Errorf("begin add message transaction: %w", err)
	}
	defer tx.Rollback()

	const ticketQuery = `
		SELECT id, status
		FROM tickets
		WHERE track_id = $1
		FOR UPDATE`
	var ticketID int64
	var status models.TicketStatus
	err = tx.QueryRowContext(ctx, ticketQuery, trackID).Scan(&ticketID, &status)
	if errors.Is(err, sql.ErrNoRows) {
		return CreatedMessageRecord{}, ErrTicketNotFound
	}
	if err != nil {
		return CreatedMessageRecord{}, fmt.Errorf("lock ticket for applicant message: %w", err)
	}

	nextStatus, ok := statusAfterApplicantMessage(status)
	if !ok {
		return CreatedMessageRecord{}, ErrInvalidTransition
	}

	if len(attachments) > 0 {
		var existingCount int
		if err := tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM attachments WHERE ticket_id = $1`, ticketID).Scan(&existingCount); err != nil {
			return CreatedMessageRecord{}, fmt.Errorf("count ticket attachments: %w", err)
		}
		if existingCount+len(attachments) > maxTicketAttachments {
			return CreatedMessageRecord{}, ErrAttachmentLimitReached
		}
	}

	const insertMessageQuery = `
		INSERT INTO messages (ticket_id, text, type, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`
	result := CreatedMessageRecord{
		Text:         text,
		Type:         models.MessageTypeApplicant,
		Attachments:  make([]CreatedAttachmentRecord, 0, len(attachments)),
		TicketStatus: nextStatus,
	}
	if err := tx.QueryRowContext(
		ctx,
		insertMessageQuery,
		ticketID,
		text,
		models.MessageTypeApplicant,
		createdAt,
	).Scan(&result.ID, &result.CreatedAt); err != nil {
		return CreatedMessageRecord{}, fmt.Errorf("insert applicant message: %w", err)
	}

	const insertAttachmentQuery = `
		INSERT INTO attachments (
			ticket_id, storage_key, safe_name, mime_type,
			original_size_bytes, size_bytes, width, height, created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at`
	for _, attachment := range attachments {
		created := CreatedAttachmentRecord{
			SafeName:  attachment.SafeName,
			MIMEType:  attachment.MIMEType,
			SizeBytes: attachment.SizeBytes,
		}
		if err := tx.QueryRowContext(
			ctx,
			insertAttachmentQuery,
			ticketID,
			attachment.StorageKey,
			attachment.SafeName,
			attachment.MIMEType,
			attachment.OriginalSizeBytes,
			attachment.SizeBytes,
			attachment.Width,
			attachment.Height,
			createdAt,
		).Scan(&created.ID, &created.CreatedAt); err != nil {
			return CreatedMessageRecord{}, fmt.Errorf("insert message attachment: %w", err)
		}
		result.Attachments = append(result.Attachments, created)
	}

	if nextStatus != status {
		if _, err := tx.ExecContext(ctx, `UPDATE tickets SET status = $2 WHERE id = $1`, ticketID, nextStatus); err != nil {
			return CreatedMessageRecord{}, fmt.Errorf("update ticket after applicant message: %w", err)
		}
	}
	if err := tx.Commit(); err != nil {
		return CreatedMessageRecord{}, fmt.Errorf("commit add message transaction: %w", err)
	}
	return result, nil
}

func statusAfterApplicantMessage(status models.TicketStatus) (models.TicketStatus, bool) {
	switch status {
	case models.TicketStatusNew, models.TicketStatusAssigned, models.TicketStatusInProgress:
		return status, true
	case models.TicketStatusNeedsClarification, models.TicketStatusAnswerReady:
		return models.TicketStatusInProgress, true
	default:
		return status, false
	}
}
