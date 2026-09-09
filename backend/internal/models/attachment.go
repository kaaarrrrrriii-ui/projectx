package models

import "time"

// Attachment maps to the attachments table. The UML leaves file-specific
// fields open; these columns contain the metadata required by file storage.
type Attachment struct {
	ID                int64     `json:"id"`
	TicketID          int64     `json:"ticket_id"`
	StorageKey        string    `json:"-"`
	SafeName          string    `json:"safe_name"`
	MIMEType          string    `json:"mime_type"`
	OriginalSizeBytes int64     `json:"original_size_bytes"`
	SizeBytes         int64     `json:"size_bytes"`
	Width             int       `json:"width"`
	Height            int       `json:"height"`
	CreatedAt         time.Time `json:"created_at"`
}
