package repos

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type AttachmentRepository struct {
	db *sql.DB
}

type AttachmentFileRecord struct {
	StorageKey string
	SafeName   string
	MIMEType   string
	SizeBytes  int64
}

func NewAttachmentRepository(db *sql.DB) *AttachmentRepository {
	return &AttachmentRepository{db: db}
}

func (repository *AttachmentRepository) GetForTicket(
	ctx context.Context,
	trackID string,
	attachmentID int64,
) (AttachmentFileRecord, error) {
	const query = `
		SELECT a.storage_key, a.safe_name, a.mime_type, a.size_bytes
		FROM attachments AS a
		JOIN tickets AS t ON t.id = a.ticket_id
		WHERE a.id = $1
		  AND t.track_id = $2`

	var record AttachmentFileRecord
	err := repository.db.QueryRowContext(ctx, query, attachmentID, trackID).Scan(
		&record.StorageKey,
		&record.SafeName,
		&record.MIMEType,
		&record.SizeBytes,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return AttachmentFileRecord{}, ErrAttachmentNotFound
	}
	if err != nil {
		return AttachmentFileRecord{}, fmt.Errorf("get attachment for ticket: %w", err)
	}
	return record, nil
}
