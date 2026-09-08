package models

import (
	"time"

	"github.com/google/uuid"
)

// ImageFormat is the normalized format of a sanitized attachment.
type ImageFormat string

const (
	ImageFormatJPEG ImageFormat = "jpeg"
	ImageFormatPNG  ImageFormat = "png"
)

// SanitizedImage describes an image after it has been decoded and re-encoded.
type SanitizedImage struct {
	Format            ImageFormat
	MIMEType          string
	Extension         string
	Width             int
	Height            int
	OriginalSizeBytes int64
	SizeBytes         int64
}

// PreparedAttachment describes a sanitized file stored under a private key.
// These fields can later be written to the attachments table.
type PreparedAttachment struct {
	StorageKey        string
	SafeName          string
	Format            ImageFormat
	MIMEType          string
	OriginalSizeBytes int64
	SizeBytes         int64
	Width             int
	Height            int
}

// UploadPurpose describes where files from a batch may be attached.
// Keeping this value explicit prevents, for example, a private note file from
// being published as an initial ticket attachment.
type UploadPurpose string

const (
	UploadPurposeInitial UploadPurpose = "initial"
	UploadPurposeMessage UploadPurpose = "message"
	UploadPurposeNote    UploadPurpose = "note"
)

// AttachmentState describes the preparation progress of one file.
type AttachmentState string

const (
	AttachmentStateProcessing AttachmentState = "processing"
	AttachmentStateReady      AttachmentState = "ready"
	AttachmentStateFailed     AttachmentState = "failed"
)

// AttachmentFailureCode is a safe reason that may be shown to a client.
// Raw decoder and storage errors must not be stored here or returned by API.
type AttachmentFailureCode string

const (
	AttachmentFailureUnsupportedFormat  AttachmentFailureCode = "unsupported_format"
	AttachmentFailureInvalidImage       AttachmentFailureCode = "invalid_image"
	AttachmentFailureInputTooLarge      AttachmentFailureCode = "input_too_large"
	AttachmentFailureDimensionsExceeded AttachmentFailureCode = "dimensions_exceeded"
	AttachmentFailureOutputTooLarge     AttachmentFailureCode = "output_too_large"
	AttachmentFailureStorage            AttachmentFailureCode = "storage_error"
)

// UploadBatch is a temporary group of files prepared for one parent object.
// Initial batches use OwnerTokenHash. Message and note batches use SessionID
// and are additionally restricted to TicketID.
//
// This is an internal domain model, not an HTTP response. Secret ownership
// fields must not be serialized directly to a client.
type UploadBatch struct {
	ID               uuid.UUID
	Purpose          UploadPurpose
	OwnerTokenHash   []byte
	SessionID        *uuid.UUID
	TicketID         *uuid.UUID
	ExpiresAt        time.Time
	ClaimedAt        *time.Time
	CleanupStartedAt *time.Time
	CreatedAt        time.Time
}

// Attachment contains metadata about one sanitized file. The binary content
// lives in private storage under StorageKey, not in the database.
//
// Before attachment all parent IDs are nil. A published file belongs to an
// initial ticket, a message, or an internal note; the database migration will
// enforce the valid combinations.
type Attachment struct {
	ID                uuid.UUID
	BatchID           uuid.UUID
	StorageKey        string
	SafeName          string
	MIMEType          string
	OriginalSizeBytes int64
	SizeBytes         int64
	Width             int
	Height            int
	State             AttachmentState
	FailureCode       *AttachmentFailureCode
	Position          int
	TicketID          *uuid.UUID
	MessageID         *uuid.UUID
	NoteID            *uuid.UUID
	CreatedAt         time.Time
	ReadyAt           *time.Time
}

// AttachmentCandidate contains the resource accounting data used for a batch.
type AttachmentCandidate struct {
	OriginalSizeBytes int64
}
