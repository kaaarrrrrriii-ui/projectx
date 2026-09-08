package service

import (
	"errors"
	"testing"

	"example.com/german/backend/internal/models"
)

func TestValidateAttachmentBatchAcceptsExactLimits(t *testing.T) {
	t.Parallel()

	limits := DefaultAttachmentLimits()
	existing := []models.AttachmentCandidate{
		{OriginalSizeBytes: 2 * 1024 * 1024},
		{OriginalSizeBytes: 2 * 1024 * 1024},
	}
	incoming := []models.AttachmentCandidate{
		{OriginalSizeBytes: 2 * 1024 * 1024},
		{OriginalSizeBytes: 2 * 1024 * 1024},
		{OriginalSizeBytes: 2 * 1024 * 1024},
	}

	if err := ValidateAttachmentBatch(limits, existing, incoming); err != nil {
		t.Fatalf("ValidateAttachmentBatch() error = %v", err)
	}
}

func TestValidateAttachmentBatchRejectsSixFiles(t *testing.T) {
	t.Parallel()

	files := make([]models.AttachmentCandidate, 6)
	if err := ValidateAttachmentBatch(DefaultAttachmentLimits(), nil, files); !errors.Is(err, ErrTooManyAttachments) {
		t.Fatalf("ValidateAttachmentBatch() error = %v, want %v", err, ErrTooManyAttachments)
	}
}

func TestValidateAttachmentBatchRejectsTotalOverLimit(t *testing.T) {
	t.Parallel()

	files := []models.AttachmentCandidate{{OriginalSizeBytes: DefaultMaxBatchBytes + 1}}
	if err := ValidateAttachmentBatch(DefaultAttachmentLimits(), nil, files); !errors.Is(err, ErrAttachmentBatchTooLarge) {
		t.Fatalf("ValidateAttachmentBatch() error = %v, want %v", err, ErrAttachmentBatchTooLarge)
	}
}

func TestValidateAttachmentBatchRejectsNegativeSize(t *testing.T) {
	t.Parallel()

	files := []models.AttachmentCandidate{{OriginalSizeBytes: -1}}
	if err := ValidateAttachmentBatch(DefaultAttachmentLimits(), nil, files); !errors.Is(err, ErrInvalidAttachmentSize) {
		t.Fatalf("ValidateAttachmentBatch() error = %v, want %v", err, ErrInvalidAttachmentSize)
	}
}
