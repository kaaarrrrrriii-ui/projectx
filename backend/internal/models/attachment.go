package models

// ImageFormat is the normalized format of a sanitized attachment.
type ImageFormat string

const (
	ImageFormatJPEG ImageFormat = "jpeg"
	ImageFormatPNG  ImageFormat = "png"
)

// SanitizedImage describes an image after it has been decoded and re-encoded.
type SanitizedImage struct {
	Format    ImageFormat
	MIMEType  string
	Extension string
	Width     int
	Height    int
	SizeBytes int64
}

// AttachmentCandidate contains the resource accounting data used for a batch.
type AttachmentCandidate struct {
	OriginalSizeBytes int64
}
