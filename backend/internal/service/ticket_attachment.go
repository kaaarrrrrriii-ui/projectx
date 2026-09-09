package service

import (
	"context"
	"errors"
	"io"

	"example.com/german/backend/internal/repos"
	"example.com/german/backend/internal/storage"
)

type attachmentRepository interface {
	GetForTicket(context.Context, string, int64) (repos.AttachmentFileRecord, error)
}

type TicketAttachmentService struct {
	repository attachmentRepository
	storage    storage.Storage
}

type AttachmentFile struct {
	Reader    io.ReadCloser
	Name      string
	MIMEType  string
	SizeBytes int64
}

func NewTicketAttachmentService(
	repository attachmentRepository,
	attachmentStorage storage.Storage,
) (*TicketAttachmentService, error) {
	if repository == nil {
		return nil, errors.New("attachment repository is required")
	}
	if attachmentStorage == nil {
		return nil, errors.New("attachment storage is required")
	}
	return &TicketAttachmentService{repository: repository, storage: attachmentStorage}, nil
}

func (service *TicketAttachmentService) Open(
	ctx context.Context,
	trackID string,
	attachmentID int64,
) (AttachmentFile, error) {
	normalized, err := normalizeTrackID(trackID)
	if err != nil {
		return AttachmentFile{}, err
	}
	record, err := service.repository.GetForTicket(ctx, normalized, attachmentID)
	if err != nil {
		return AttachmentFile{}, err
	}
	reader, err := service.storage.Open(ctx, record.StorageKey)
	if err != nil {
		return AttachmentFile{}, err
	}
	return AttachmentFile{
		Reader:    reader,
		Name:      record.SafeName,
		MIMEType:  record.MIMEType,
		SizeBytes: record.SizeBytes,
	}, nil
}
