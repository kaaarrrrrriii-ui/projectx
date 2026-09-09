package service

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"io"

	"example.com/german/backend/internal/storage"
)

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

// AttachmentPreparer connects image sanitizing with private storage.
// It does not create database records or attach a file to a ticket.
type AttachmentPreparer struct {
	processor   *ImageProcessor
	storage     storage.Storage
	generateKey func() (string, error)
}

func NewAttachmentPreparer(processor *ImageProcessor, attachmentStorage storage.Storage) (*AttachmentPreparer, error) {
	if processor == nil {
		return nil, errors.New("image processor is required")
	}
	if attachmentStorage == nil {
		return nil, errors.New("attachment storage is required")
	}

	return &AttachmentPreparer{
		processor:   processor,
		storage:     attachmentStorage,
		generateKey: storage.NewKey,
	}, nil
}

// Prepare sanitizes an image and persists only the sanitized result.
// The returned data is ready to be saved in the attachments table later.
func (p *AttachmentPreparer) Prepare(ctx context.Context, source io.Reader) (PreparedAttachment, error) {
	if source == nil {
		return PreparedAttachment{}, errors.New("attachment source is required")
	}

	preparedFile, err := secureTempFile("attachment-prepared-*")
	if err != nil {
		return PreparedAttachment{}, fmt.Errorf("create prepared attachment buffer: %w", err)
	}
	defer removeTempFile(preparedFile)

	imageInfo, err := p.processor.Sanitize(ctx, source, preparedFile)
	if err != nil {
		return PreparedAttachment{}, err
	}
	if _, err := preparedFile.Seek(0, io.SeekStart); err != nil {
		return PreparedAttachment{}, fmt.Errorf("rewind prepared attachment: %w", err)
	}

	storageKey, err := p.generateKey()
	if err != nil {
		return PreparedAttachment{}, fmt.Errorf("generate attachment storage key: %w", err)
	}
	safeName, err := safeAttachmentName(storageKey, imageInfo.Extension)
	if err != nil {
		return PreparedAttachment{}, err
	}

	if err := p.storage.Put(ctx, storageKey, preparedFile); err != nil {
		return PreparedAttachment{}, fmt.Errorf("store prepared attachment: %w", err)
	}

	return PreparedAttachment{
		StorageKey:        storageKey,
		SafeName:          safeName,
		Format:            imageInfo.Format,
		MIMEType:          imageInfo.MIMEType,
		OriginalSizeBytes: imageInfo.OriginalSizeBytes,
		SizeBytes:         imageInfo.SizeBytes,
		Width:             imageInfo.Width,
		Height:            imageInfo.Height,
	}, nil
}

func safeAttachmentName(storageKey, extension string) (string, error) {
	if len(storageKey) != 64 {
		return "", errors.New("invalid attachment storage key")
	}
	if _, err := hex.DecodeString(storageKey); err != nil {
		return "", errors.New("invalid attachment storage key")
	}
	if extension != ".jpg" && extension != ".png" {
		return "", errors.New("unsupported sanitized image extension")
	}
	return "image-" + storageKey[:12] + extension, nil
}
