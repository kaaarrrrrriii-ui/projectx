package service

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"example.com/german/backend/internal/models"
	"example.com/german/backend/internal/repos"
	"example.com/german/backend/internal/storage"
)

type stubMessageRepository struct {
	attachments []repos.NewAttachmentRecord
	text        string
	result      repos.CreatedMessageRecord
	err         error
}

func (repository *stubMessageRepository) AddApplicantMessage(
	_ context.Context,
	_ string,
	text string,
	attachments []repos.NewAttachmentRecord,
	_ int,
	_ time.Time,
) (repos.CreatedMessageRecord, error) {
	repository.text = text
	repository.attachments = attachments
	return repository.result, repository.err
}

func newTicketMessageTestService(t *testing.T, repository messageRepository) (*TicketMessageService, *storage.Local) {
	t.Helper()
	attachmentStorage, err := storage.NewLocal(t.TempDir())
	if err != nil {
		t.Fatalf("NewLocal() error = %v", err)
	}
	processor, err := NewImageProcessor(DefaultImageProcessorConfig())
	if err != nil {
		t.Fatalf("NewImageProcessor() error = %v", err)
	}
	preparer, err := NewAttachmentPreparer(processor, attachmentStorage)
	if err != nil {
		t.Fatalf("NewAttachmentPreparer() error = %v", err)
	}
	messageService, err := NewTicketMessageService(repository, preparer, attachmentStorage)
	if err != nil {
		t.Fatalf("NewTicketMessageService() error = %v", err)
	}
	return messageService, attachmentStorage
}

func TestTicketMessageServiceRequiresTextEvenWithAttachment(t *testing.T) {
	t.Parallel()

	repository := &stubMessageRepository{}
	messageService, _ := newTicketMessageTestService(t, repository)
	image := encodeTestPNG(t, 1, 1)
	_, err := messageService.AddApplicantMessage(context.Background(), "ОТК-ABCD-2345", "   ", []UploadedAttachment{{
		Source: bytes.NewReader(image),
		Size:   int64(len(image)),
	}})
	if !errors.Is(err, ErrMessageTextRequired) {
		t.Fatalf("AddApplicantMessage() error = %v, want %v", err, ErrMessageTextRequired)
	}
	if len(repository.attachments) != 0 {
		t.Fatal("repository called for a message without text")
	}
}

func TestTicketMessageServiceLimitsTextLength(t *testing.T) {
	t.Parallel()
	repository := &stubMessageRepository{}
	messageService, _ := newTicketMessageTestService(t, repository)

	_, err := messageService.AddApplicantMessage(context.Background(), "ОТК-ABCD-2345", strings.Repeat("я", MaxApplicantMessageCharacters+1), nil)
	if !errors.Is(err, ErrTextTooLong) {
		t.Fatalf("AddApplicantMessage() error = %v, want %v", err, ErrTextTooLong)
	}
}

func TestTicketMessageServiceSanitizesAndPersistsImageMetadata(t *testing.T) {
	t.Parallel()

	createdAt := time.Date(2026, 9, 9, 10, 30, 0, 0, time.UTC)
	repository := &stubMessageRepository{result: repos.CreatedMessageRecord{
		ID:           15,
		Text:         "Дополнительная информация",
		Type:         models.MessageTypeApplicant,
		CreatedAt:    createdAt,
		TicketStatus: models.TicketStatusInProgress,
		Attachments: []repos.CreatedAttachmentRecord{{
			ID:        8,
			SafeName:  "image-safe.png",
			MIMEType:  "image/png",
			SizeBytes: 70,
			CreatedAt: createdAt,
		}},
	}}
	messageService, attachmentStorage := newTicketMessageTestService(t, repository)
	image := encodeTestPNG(t, 2, 3)
	response, err := messageService.AddApplicantMessage(
		context.Background(),
		"ОТК-ABCD-2345",
		"  Дополнительная информация  ",
		[]UploadedAttachment{{Source: bytes.NewReader(image), Size: int64(len(image))}},
	)
	if err != nil {
		t.Fatalf("AddApplicantMessage() error = %v", err)
	}
	if repository.text != "Дополнительная информация" || len(repository.attachments) != 1 {
		t.Fatalf("repository text = %q, attachments = %+v", repository.text, repository.attachments)
	}
	prepared := repository.attachments[0]
	if prepared.MIMEType != "image/png" || prepared.Width != 2 || prepared.Height != 3 || prepared.StorageKey == "" {
		t.Fatalf("prepared attachment = %+v", prepared)
	}
	stored, err := attachmentStorage.Open(context.Background(), prepared.StorageKey)
	if err != nil {
		t.Fatalf("stored attachment is unavailable: %v", err)
	}
	_ = stored.Close()
	if response.TicketStatus != "in_progress" || len(response.Message.Attachments) != 1 || response.Message.Attachments[0].ID != 8 {
		t.Fatalf("response = %+v", response)
	}
}

func TestTicketMessageServiceDeletesPreparedFilesWhenDatabaseFails(t *testing.T) {
	t.Parallel()

	repository := &stubMessageRepository{err: errors.New("database unavailable")}
	messageService, attachmentStorage := newTicketMessageTestService(t, repository)
	image := encodeTestPNG(t, 1, 1)
	_, err := messageService.AddApplicantMessage(
		context.Background(),
		"ОТК-ABCD-2345",
		"Текст",
		[]UploadedAttachment{{Source: bytes.NewReader(image), Size: int64(len(image))}},
	)
	if err == nil {
		t.Fatal("AddApplicantMessage() error = nil")
	}
	if len(repository.attachments) != 1 {
		t.Fatalf("repository attachments = %+v", repository.attachments)
	}
	if _, openErr := attachmentStorage.Open(context.Background(), repository.attachments[0].StorageKey); openErr == nil {
		t.Fatal("prepared file still exists after database failure")
	}
}
