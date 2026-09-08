package models

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

// AttachmentCandidate contains the resource accounting data used for a batch.
type AttachmentCandidate struct {
	OriginalSizeBytes int64
}
