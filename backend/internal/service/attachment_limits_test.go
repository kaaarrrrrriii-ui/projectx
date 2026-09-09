package service

import (
	"errors"
	"testing"
)

func TestValidateAttachmentBatchAcceptsFiveFilesUpToLimitEach(t *testing.T) {
	t.Parallel()

	limits := DefaultAttachmentLimits()
	existing := []AttachmentCandidate{
		{OriginalSizeBytes: DefaultMaxAttachmentBytes},
		{OriginalSizeBytes: DefaultMaxAttachmentBytes},
	}
	incoming := []AttachmentCandidate{
		{OriginalSizeBytes: DefaultMaxAttachmentBytes},
		{OriginalSizeBytes: DefaultMaxAttachmentBytes},
		{OriginalSizeBytes: DefaultMaxAttachmentBytes},
	}

	if err := ValidateAttachmentBatch(limits, existing, incoming); err != nil {
		t.Fatalf("ValidateAttachmentBatch() error = %v", err)
	}
}

func TestValidateAttachmentBatchRejectsSixFiles(t *testing.T) {
	t.Parallel()

	files := make([]AttachmentCandidate, 6)
	if err := ValidateAttachmentBatch(DefaultAttachmentLimits(), nil, files); !errors.Is(err, ErrTooManyAttachments) {
		t.Fatalf("ValidateAttachmentBatch() error = %v, want %v", err, ErrTooManyAttachments)
	}
}

func TestValidateAttachmentBatchRejectsFileOverLimit(t *testing.T) {
	t.Parallel()

	files := []AttachmentCandidate{{OriginalSizeBytes: DefaultMaxAttachmentBytes + 1}}
	if err := ValidateAttachmentBatch(DefaultAttachmentLimits(), nil, files); !errors.Is(err, ErrAttachmentTooLarge) {
		t.Fatalf("ValidateAttachmentBatch() error = %v, want %v", err, ErrAttachmentTooLarge)
	}
}

func TestValidateAttachmentBatchRejectsNegativeSize(t *testing.T) {
	t.Parallel()

	files := []AttachmentCandidate{{OriginalSizeBytes: -1}}
	if err := ValidateAttachmentBatch(DefaultAttachmentLimits(), nil, files); !errors.Is(err, ErrInvalidAttachmentSize) {
		t.Fatalf("ValidateAttachmentBatch() error = %v, want %v", err, ErrInvalidAttachmentSize)
	}
}
