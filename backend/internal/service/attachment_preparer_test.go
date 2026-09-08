package service

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"example.com/german/backend/internal/storage"
)

func TestAttachmentPreparerSanitizesAndStoresImage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		input        func(*testing.T) []byte
		metadataText string
		wantMIME     string
		wantName     string
	}{
		{
			name: "jpeg",
			input: func(t *testing.T) []byte {
				return insertJPEGApplicationSegment(
					encodeTestJPEG(t, 3, 2),
					[]byte("Exif\x00\x00GPSLatitude=test-location"),
				)
			},
			metadataText: "GPSLatitude",
			wantMIME:     "image/jpeg",
			wantName:     "image-aaaaaaaaaaaa.jpg",
		},
		{
			name: "png",
			input: func(t *testing.T) []byte {
				return insertPNGChunk(t, encodeTestPNG(t, 3, 2), "tEXt", []byte("Location\x00test-location"))
			},
			metadataText: "Location",
			wantMIME:     "image/png",
			wantName:     "image-aaaaaaaaaaaa.png",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			input := test.input(t)
			localStorage, err := storage.NewLocal(t.TempDir())
			if err != nil {
				t.Fatalf("NewLocal() error = %v", err)
			}
			processor := newTestImageProcessor(t, DefaultImageProcessorConfig())
			preparer, err := NewAttachmentPreparer(processor, localStorage)
			if err != nil {
				t.Fatalf("NewAttachmentPreparer() error = %v", err)
			}
			preparer.generateKey = func() (string, error) {
				return strings.Repeat("a", 64), nil
			}

			prepared, err := preparer.Prepare(context.Background(), bytes.NewReader(input))
			if err != nil {
				t.Fatalf("Prepare() error = %v", err)
			}
			if prepared.SafeName != test.wantName || prepared.MIMEType != test.wantMIME {
				t.Fatalf("Prepare() result = %+v", prepared)
			}
			if prepared.OriginalSizeBytes != int64(len(input)) {
				t.Fatalf("Prepare() original size = %d, want %d", prepared.OriginalSizeBytes, len(input))
			}
			if prepared.SizeBytes <= 0 || prepared.Width != 3 || prepared.Height != 2 {
				t.Fatalf("Prepare() image data = %+v", prepared)
			}

			stored, err := localStorage.Open(context.Background(), prepared.StorageKey)
			if err != nil {
				t.Fatalf("Open() prepared attachment error = %v", err)
			}
			storedBytes, err := io.ReadAll(stored)
			closeErr := stored.Close()
			if err != nil {
				t.Fatalf("read prepared attachment: %v", err)
			}
			if closeErr != nil {
				t.Fatalf("close prepared attachment: %v", closeErr)
			}
			if bytes.Contains(storedBytes, []byte(test.metadataText)) || bytes.Contains(storedBytes, []byte("test-location")) {
				t.Fatal("stored attachment still contains source metadata")
			}
			if int64(len(storedBytes)) != prepared.SizeBytes {
				t.Fatalf("stored size = %d, result size = %d", len(storedBytes), prepared.SizeBytes)
			}
		})
	}
}

func TestAttachmentPreparerDoesNotStoreRejectedImage(t *testing.T) {
	t.Parallel()

	attachmentStorage := &recordingStorage{}
	processor := newTestImageProcessor(t, DefaultImageProcessorConfig())
	preparer, err := NewAttachmentPreparer(processor, attachmentStorage)
	if err != nil {
		t.Fatalf("NewAttachmentPreparer() error = %v", err)
	}

	_, err = preparer.Prepare(context.Background(), strings.NewReader("not an image"))
	if !errors.Is(err, ErrUnsupportedImage) {
		t.Fatalf("Prepare() error = %v, want %v", err, ErrUnsupportedImage)
	}
	if attachmentStorage.putCalls != 0 {
		t.Fatalf("storage Put() calls = %d, want 0", attachmentStorage.putCalls)
	}
}

func TestAttachmentPreparerReturnsStorageError(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("storage unavailable")
	attachmentStorage := &recordingStorage{putErr: wantErr}
	processor := newTestImageProcessor(t, DefaultImageProcessorConfig())
	preparer, err := NewAttachmentPreparer(processor, attachmentStorage)
	if err != nil {
		t.Fatalf("NewAttachmentPreparer() error = %v", err)
	}

	_, err = preparer.Prepare(context.Background(), bytes.NewReader(encodeTestPNG(t, 1, 1)))
	if !errors.Is(err, wantErr) {
		t.Fatalf("Prepare() error = %v, want wrapped %v", err, wantErr)
	}
	if attachmentStorage.putCalls != 1 {
		t.Fatalf("storage Put() calls = %d, want 1", attachmentStorage.putCalls)
	}
}

type recordingStorage struct {
	putCalls int
	putErr   error
}

func (s *recordingStorage) Put(_ context.Context, _ string, _ io.Reader) error {
	s.putCalls++
	return s.putErr
}

func (*recordingStorage) Open(context.Context, string) (io.ReadCloser, error) {
	return nil, errors.New("not implemented")
}

func (*recordingStorage) Delete(context.Context, string) error {
	return errors.New("not implemented")
}
