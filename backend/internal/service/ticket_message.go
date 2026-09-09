package service

import (
	"context"
	"errors"
	"io"
	"strings"
	"time"

	"example.com/german/backend/internal/repos"
	"example.com/german/backend/internal/storage"
)

type UploadedAttachment struct {
	Source io.Reader
	Size   int64
}

type messageRepository interface {
	AddApplicantMessage(
		context.Context,
		string,
		string,
		[]repos.NewAttachmentRecord,
		int,
		time.Time,
	) (repos.CreatedMessageRecord, error)
}

type TicketMessageService struct {
	repository messageRepository
	preparer   *AttachmentPreparer
	storage    storage.Storage
	limits     AttachmentLimits
	now        func() time.Time
}

type CreateMessageResponse struct {
	Message      MessageResponse `json:"message"`
	TicketStatus string          `json:"ticket_status"`
}

func NewTicketMessageService(
	repository messageRepository,
	preparer *AttachmentPreparer,
	attachmentStorage storage.Storage,
) (*TicketMessageService, error) {
	if repository == nil {
		return nil, errors.New("message repository is required")
	}
	if preparer == nil {
		return nil, errors.New("attachment preparer is required")
	}
	if attachmentStorage == nil {
		return nil, errors.New("attachment storage is required")
	}
	return &TicketMessageService{
		repository: repository,
		preparer:   preparer,
		storage:    attachmentStorage,
		limits:     DefaultAttachmentLimits(),
		now:        time.Now,
	}, nil
}

func (service *TicketMessageService) AddApplicantMessage(
	ctx context.Context,
	trackID string,
	text string,
	uploads []UploadedAttachment,
) (CreateMessageResponse, error) {
	normalized, err := normalizeTrackID(trackID)
	if err != nil {
		return CreateMessageResponse{}, err
	}
	text = strings.TrimSpace(text)
	if text == "" {
		return CreateMessageResponse{}, ErrMessageTextRequired
	}

	candidates := make([]AttachmentCandidate, 0, len(uploads))
	for _, upload := range uploads {
		candidates = append(candidates, AttachmentCandidate{OriginalSizeBytes: upload.Size})
	}
	if err := ValidateAttachmentBatch(service.limits, nil, candidates); err != nil {
		return CreateMessageResponse{}, err
	}

	prepared := make([]PreparedAttachment, 0, len(uploads))
	cleanup := func() {
		for _, attachment := range prepared {
			_ = service.storage.Delete(context.Background(), attachment.StorageKey)
		}
	}
	for _, upload := range uploads {
		if upload.Source == nil {
			cleanup()
			return CreateMessageResponse{}, errors.New("attachment source is required")
		}
		attachment, err := service.preparer.Prepare(ctx, upload.Source)
		if err != nil {
			cleanup()
			return CreateMessageResponse{}, err
		}
		prepared = append(prepared, attachment)
	}

	newAttachments := make([]repos.NewAttachmentRecord, 0, len(prepared))
	for _, attachment := range prepared {
		newAttachments = append(newAttachments, repos.NewAttachmentRecord{
			StorageKey:        attachment.StorageKey,
			SafeName:          attachment.SafeName,
			MIMEType:          attachment.MIMEType,
			OriginalSizeBytes: attachment.OriginalSizeBytes,
			SizeBytes:         attachment.SizeBytes,
			Width:             attachment.Width,
			Height:            attachment.Height,
		})
	}

	record, err := service.repository.AddApplicantMessage(
		ctx,
		normalized,
		text,
		newAttachments,
		service.limits.MaxFiles,
		service.now().UTC(),
	)
	if err != nil {
		cleanup()
		return CreateMessageResponse{}, err
	}

	response := CreateMessageResponse{
		Message: MessageResponse{
			ID:          record.ID,
			Text:        record.Text,
			Type:        record.Type.String(),
			CreatedAt:   record.CreatedAt.UTC(),
			Attachments: make([]AttachmentResponse, 0, len(record.Attachments)),
		},
		TicketStatus: record.TicketStatus.String(),
	}
	for _, attachment := range record.Attachments {
		response.Message.Attachments = append(response.Message.Attachments, AttachmentResponse{
			ID:        attachment.ID,
			Name:      attachment.SafeName,
			MIMEType:  attachment.MIMEType,
			SizeBytes: attachment.SizeBytes,
		})
	}
	return response, nil
}
