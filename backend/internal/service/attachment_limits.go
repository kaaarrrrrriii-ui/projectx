package service

import (
	"errors"
	"fmt"
)

const (
	DefaultMaxAttachmentFiles = 5
	DefaultMaxBatchBytes      = int64(10 * 1024 * 1024)
)

var (
	ErrTooManyAttachments      = errors.New("too many attachments")
	ErrAttachmentBatchTooLarge = errors.New("attachment batch is too large")
	ErrInvalidAttachmentSize   = errors.New("invalid attachment size")
)

// AttachmentLimits are applied to the original files in one upload batch.
type AttachmentLimits struct {
	MaxFiles      int
	MaxTotalBytes int64
}

type AttachmentCandidate struct {
	OriginalSizeBytes int64
}

func DefaultAttachmentLimits() AttachmentLimits {
	return AttachmentLimits{
		MaxFiles:      DefaultMaxAttachmentFiles,
		MaxTotalBytes: DefaultMaxBatchBytes,
	}
}

// ValidateAttachmentBatch validates both existing and newly uploaded files.
// The database-backed caller must repeat this check while holding the batch
// lock so concurrent uploads cannot exceed the limits together.
func ValidateAttachmentBatch(limits AttachmentLimits, existing, incoming []AttachmentCandidate) error {
	if limits.MaxFiles <= 0 || limits.MaxTotalBytes <= 0 {
		return fmt.Errorf("invalid attachment limits")
	}

	if len(existing)+len(incoming) > limits.MaxFiles {
		return fmt.Errorf("%w: maximum is %d", ErrTooManyAttachments, limits.MaxFiles)
	}

	var total int64
	for _, candidate := range appendCandidates(existing, incoming) {
		if candidate.OriginalSizeBytes < 0 {
			return ErrInvalidAttachmentSize
		}
		if candidate.OriginalSizeBytes > limits.MaxTotalBytes-total {
			return fmt.Errorf("%w: maximum is %d bytes", ErrAttachmentBatchTooLarge, limits.MaxTotalBytes)
		}
		total += candidate.OriginalSizeBytes
	}

	return nil
}

func appendCandidates(existing, incoming []AttachmentCandidate) []AttachmentCandidate {
	all := make([]AttachmentCandidate, 0, len(existing)+len(incoming))
	all = append(all, existing...)
	return append(all, incoming...)
}
