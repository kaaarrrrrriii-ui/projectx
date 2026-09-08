package service

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"hash/crc32"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"testing"
)

func TestImageProcessorSanitizesJPEGMetadata(t *testing.T) {
	t.Parallel()

	original := encodeTestJPEG(t, 3, 2)
	withMetadata := insertJPEGApplicationSegment(original, []byte("Exif\x00\x00GPSLatitude=test-location"))
	processor := newTestImageProcessor(t, DefaultImageProcessorConfig())

	var output bytes.Buffer
	result, err := processor.Sanitize(context.Background(), bytes.NewReader(withMetadata), &output)
	if err != nil {
		t.Fatalf("Sanitize() error = %v", err)
	}
	if result.MIMEType != "image/jpeg" || result.Extension != ".jpg" || result.Width != 3 || result.Height != 2 {
		t.Fatalf("Sanitize() result = %+v", result)
	}
	if bytes.Contains(output.Bytes(), []byte("Exif")) || bytes.Contains(output.Bytes(), []byte("GPSLatitude")) {
		t.Fatal("sanitized JPEG still contains injected metadata")
	}
	if _, err := jpeg.Decode(bytes.NewReader(output.Bytes())); err != nil {
		t.Fatalf("decode sanitized JPEG: %v", err)
	}
}

func TestImageProcessorSanitizesPNGMetadata(t *testing.T) {
	t.Parallel()

	original := encodeTestPNG(t, 2, 3)
	withMetadata := insertPNGChunk(t, original, "tEXt", []byte("Location\x00test-location"))
	processor := newTestImageProcessor(t, DefaultImageProcessorConfig())

	var output bytes.Buffer
	result, err := processor.Sanitize(context.Background(), bytes.NewReader(withMetadata), &output)
	if err != nil {
		t.Fatalf("Sanitize() error = %v", err)
	}
	if result.MIMEType != "image/png" || result.Extension != ".png" || result.Width != 2 || result.Height != 3 {
		t.Fatalf("Sanitize() result = %+v", result)
	}
	if bytes.Contains(output.Bytes(), []byte("Location")) || bytes.Contains(output.Bytes(), []byte("test-location")) {
		t.Fatal("sanitized PNG still contains injected metadata")
	}
	if _, err := png.Decode(bytes.NewReader(output.Bytes())); err != nil {
		t.Fatalf("decode sanitized PNG: %v", err)
	}
}

func TestImageProcessorDetectsActualFormat(t *testing.T) {
	t.Parallel()

	processor := newTestImageProcessor(t, DefaultImageProcessorConfig())
	var output bytes.Buffer
	result, err := processor.Sanitize(context.Background(), bytes.NewReader(encodeTestPNG(t, 1, 1)), &output)
	if err != nil {
		t.Fatalf("Sanitize() error = %v", err)
	}
	if result.MIMEType != "image/png" {
		t.Fatalf("Sanitize() MIME type = %q, want image/png", result.MIMEType)
	}
}

func TestImageProcessorRejectsUnsupportedAndDamagedFiles(t *testing.T) {
	t.Parallel()

	processor := newTestImageProcessor(t, DefaultImageProcessorConfig())
	tests := []struct {
		name  string
		input []byte
		want  error
	}{
		{name: "text", input: []byte("not an image"), want: ErrUnsupportedImage},
		{name: "damaged jpeg", input: []byte{0xff, 0xd8, 0xff, 0xe1, 0x00}, want: ErrInvalidImage},
		{name: "damaged png", input: []byte("\x89PNG\r\n\x1a\ninvalid"), want: ErrInvalidImage},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var output bytes.Buffer
			_, err := processor.Sanitize(context.Background(), bytes.NewReader(test.input), &output)
			if !errors.Is(err, test.want) {
				t.Fatalf("Sanitize() error = %v, want %v", err, test.want)
			}
			if output.Len() != 0 {
				t.Fatalf("Sanitize() wrote %d bytes for rejected input", output.Len())
			}
		})
	}
}

func TestImageProcessorRejectsResourceLimitViolations(t *testing.T) {
	t.Parallel()

	t.Run("input bytes", func(t *testing.T) {
		t.Parallel()
		config := DefaultImageProcessorConfig()
		config.MaxInputBytes = 4
		processor := newTestImageProcessor(t, config)
		var output bytes.Buffer
		_, err := processor.Sanitize(context.Background(), bytes.NewReader(encodeTestPNG(t, 1, 1)), &output)
		if !errors.Is(err, ErrImageInputTooLarge) {
			t.Fatalf("Sanitize() error = %v, want %v", err, ErrImageInputTooLarge)
		}
	})

	t.Run("pixel count", func(t *testing.T) {
		t.Parallel()
		config := DefaultImageProcessorConfig()
		config.MaxPixels = 3
		processor := newTestImageProcessor(t, config)
		var output bytes.Buffer
		_, err := processor.Sanitize(context.Background(), bytes.NewReader(encodeTestPNG(t, 2, 2)), &output)
		if !errors.Is(err, ErrImageTooManyPixels) {
			t.Fatalf("Sanitize() error = %v, want %v", err, ErrImageTooManyPixels)
		}
	})

	t.Run("output bytes", func(t *testing.T) {
		t.Parallel()
		config := DefaultImageProcessorConfig()
		config.MaxOutputBytes = 8
		processor := newTestImageProcessor(t, config)
		var output bytes.Buffer
		_, err := processor.Sanitize(context.Background(), bytes.NewReader(encodeTestPNG(t, 2, 2)), &output)
		if !errors.Is(err, ErrImageOutputTooLarge) {
			t.Fatalf("Sanitize() error = %v, want %v", err, ErrImageOutputTooLarge)
		}
		if output.Len() != 0 {
			t.Fatalf("Sanitize() wrote %d bytes for oversized output", output.Len())
		}
	})
}

func newTestImageProcessor(t *testing.T, config ImageProcessorConfig) *ImageProcessor {
	t.Helper()
	processor, err := NewImageProcessor(config)
	if err != nil {
		t.Fatalf("NewImageProcessor() error = %v", err)
	}
	return processor
}

func encodeTestJPEG(t *testing.T, width, height int) []byte {
	t.Helper()
	var output bytes.Buffer
	if err := jpeg.Encode(&output, testImage(width, height), &jpeg.Options{Quality: 90}); err != nil {
		t.Fatalf("encode JPEG fixture: %v", err)
	}
	return output.Bytes()
}

func encodeTestPNG(t *testing.T, width, height int) []byte {
	t.Helper()
	var output bytes.Buffer
	if err := png.Encode(&output, testImage(width, height)); err != nil {
		t.Fatalf("encode PNG fixture: %v", err)
	}
	return output.Bytes()
}

func testImage(width, height int) image.Image {
	result := image.NewNRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			result.Set(x, y, color.NRGBA{R: uint8(x * 31), G: uint8(y * 47), B: 120, A: 255})
		}
	}
	return result
}

func insertJPEGApplicationSegment(encoded, payload []byte) []byte {
	segmentLength := len(payload) + 2
	segment := []byte{0xff, 0xe1, byte(segmentLength >> 8), byte(segmentLength)}
	segment = append(segment, payload...)
	result := append([]byte{}, encoded[:2]...)
	result = append(result, segment...)
	return append(result, encoded[2:]...)
}

func insertPNGChunk(t *testing.T, encoded []byte, chunkType string, payload []byte) []byte {
	t.Helper()
	if len(chunkType) != 4 {
		t.Fatal("PNG chunk type must contain four bytes")
	}
	const signatureSize = 8
	const ihdrChunkSize = 4 + 4 + 13 + 4
	position := signatureSize + ihdrChunkSize
	chunk := make([]byte, 0, 12+len(payload))
	length := make([]byte, 4)
	binary.BigEndian.PutUint32(length, uint32(len(payload)))
	chunk = append(chunk, length...)
	chunk = append(chunk, chunkType...)
	chunk = append(chunk, payload...)
	checksum := make([]byte, 4)
	binary.BigEndian.PutUint32(checksum, crc32.ChecksumIEEE(chunk[4:]))
	chunk = append(chunk, checksum...)

	result := append([]byte{}, encoded[:position]...)
	result = append(result, chunk...)
	return append(result, encoded[position:]...)
}
