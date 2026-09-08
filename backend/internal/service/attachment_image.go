package service

import (
	"context"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"

	"example.com/german/backend/internal/models"
)

const (
	DefaultMaxImagePixels         = uint64(25_000_000)
	DefaultMaxSanitizedImageBytes = int64(20 * 1024 * 1024)
	DefaultJPEGQuality            = 90
	DefaultMaxConcurrentDecodes   = 2
)

var (
	ErrImageInputTooLarge  = errors.New("image input is too large")
	ErrUnsupportedImage    = errors.New("unsupported image format")
	ErrInvalidImage        = errors.New("invalid image")
	ErrImageTooManyPixels  = errors.New("image dimensions exceed the limit")
	ErrImageOutputTooLarge = errors.New("sanitized image is too large")
)

type ImageProcessorConfig struct {
	MaxInputBytes  int64
	MaxOutputBytes int64
	MaxPixels      uint64
	JPEGQuality    int
	MaxConcurrent  int
}

func DefaultImageProcessorConfig() ImageProcessorConfig {
	return ImageProcessorConfig{
		MaxInputBytes:  DefaultMaxBatchBytes,
		MaxOutputBytes: DefaultMaxSanitizedImageBytes,
		MaxPixels:      DefaultMaxImagePixels,
		JPEGQuality:    DefaultJPEGQuality,
		MaxConcurrent:  DefaultMaxConcurrentDecodes,
	}
}

// ImageProcessor validates, fully decodes, and re-encodes supported images.
// Re-encoding deliberately does not copy EXIF, GPS, comments, or PNG chunks.
type ImageProcessor struct {
	config    ImageProcessorConfig
	semaphore chan struct{}
}

func NewImageProcessor(config ImageProcessorConfig) (*ImageProcessor, error) {
	if config.MaxInputBytes <= 0 || config.MaxOutputBytes <= 0 || config.MaxPixels == 0 {
		return nil, errors.New("image size limits must be positive")
	}
	if config.JPEGQuality < 1 || config.JPEGQuality > 100 {
		return nil, errors.New("JPEG quality must be between 1 and 100")
	}
	if config.MaxConcurrent <= 0 {
		return nil, errors.New("maximum concurrent decodes must be positive")
	}

	return &ImageProcessor{
		config:    config,
		semaphore: make(chan struct{}, config.MaxConcurrent),
	}, nil
}

// Sanitize writes to dst only after a complete sanitized image has been
// produced. Callers that persist dst must still use an atomic storage write.
func (p *ImageProcessor) Sanitize(ctx context.Context, src io.Reader, dst io.Writer) (models.SanitizedImage, error) {
	input, err := secureTempFile("attachment-source-*")
	if err != nil {
		return models.SanitizedImage{}, fmt.Errorf("create image input buffer: %w", err)
	}
	defer removeTempFile(input)

	written, err := copyWithContext(ctx, input, io.LimitReader(src, p.config.MaxInputBytes+1))
	if err != nil {
		return models.SanitizedImage{}, fmt.Errorf("buffer image input: %w", err)
	}
	if written > p.config.MaxInputBytes {
		return models.SanitizedImage{}, ErrImageInputTooLarge
	}
	if _, err := input.Seek(0, io.SeekStart); err != nil {
		return models.SanitizedImage{}, fmt.Errorf("rewind image input: %w", err)
	}

	if err := p.acquire(ctx); err != nil {
		return models.SanitizedImage{}, err
	}
	defer p.release()

	format, err := detectImageFormat(input)
	if err != nil {
		return models.SanitizedImage{}, err
	}
	if _, err := input.Seek(0, io.SeekStart); err != nil {
		return models.SanitizedImage{}, fmt.Errorf("rewind image input: %w", err)
	}

	width, height, err := decodeImageConfig(input, format)
	if err != nil {
		return models.SanitizedImage{}, ErrInvalidImage
	}
	if width <= 0 || height <= 0 || uint64(width)*uint64(height) > p.config.MaxPixels {
		return models.SanitizedImage{}, ErrImageTooManyPixels
	}
	if _, err := input.Seek(0, io.SeekStart); err != nil {
		return models.SanitizedImage{}, fmt.Errorf("rewind image input: %w", err)
	}

	decoded, err := decodeImage(input, format)
	if err != nil {
		return models.SanitizedImage{}, ErrInvalidImage
	}

	output, err := secureTempFile("attachment-sanitized-*")
	if err != nil {
		return models.SanitizedImage{}, fmt.Errorf("create image output buffer: %w", err)
	}
	defer removeTempFile(output)

	limited := &limitWriter{writer: output, remaining: p.config.MaxOutputBytes}
	if err := encodeImage(limited, decoded, format, p.config.JPEGQuality); err != nil {
		if errors.Is(err, ErrImageOutputTooLarge) {
			return models.SanitizedImage{}, err
		}
		return models.SanitizedImage{}, fmt.Errorf("encode sanitized image: %w", err)
	}
	size, err := output.Seek(0, io.SeekCurrent)
	if err != nil {
		return models.SanitizedImage{}, fmt.Errorf("measure sanitized image: %w", err)
	}
	if _, err := output.Seek(0, io.SeekStart); err != nil {
		return models.SanitizedImage{}, fmt.Errorf("rewind sanitized image: %w", err)
	}
	if _, err := copyWithContext(ctx, dst, output); err != nil {
		return models.SanitizedImage{}, fmt.Errorf("write sanitized image: %w", err)
	}

	result := models.SanitizedImage{
		Format:    format,
		Width:     width,
		Height:    height,
		SizeBytes: size,
	}
	if format == models.ImageFormatJPEG {
		result.MIMEType = "image/jpeg"
		result.Extension = ".jpg"
	} else {
		result.MIMEType = "image/png"
		result.Extension = ".png"
	}

	return result, nil
}

func (p *ImageProcessor) acquire(ctx context.Context) error {
	select {
	case p.semaphore <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (p *ImageProcessor) release() {
	<-p.semaphore
}

func detectImageFormat(input io.Reader) (models.ImageFormat, error) {
	header := make([]byte, 8)
	n, err := io.ReadFull(input, header)
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
		return "", ErrInvalidImage
	}
	header = header[:n]

	if len(header) >= 3 && header[0] == 0xff && header[1] == 0xd8 && header[2] == 0xff {
		return models.ImageFormatJPEG, nil
	}
	if len(header) == 8 && string(header) == "\x89PNG\r\n\x1a\n" {
		return models.ImageFormatPNG, nil
	}
	return "", ErrUnsupportedImage
}

func decodeImageConfig(input io.Reader, format models.ImageFormat) (int, int, error) {
	if format == models.ImageFormatJPEG {
		config, err := jpeg.DecodeConfig(input)
		return config.Width, config.Height, err
	}
	config, err := png.DecodeConfig(input)
	return config.Width, config.Height, err
}

func decodeImage(input io.Reader, format models.ImageFormat) (image.Image, error) {
	if format == models.ImageFormatJPEG {
		return jpeg.Decode(input)
	}
	return png.Decode(input)
}

func encodeImage(dst io.Writer, src image.Image, format models.ImageFormat, jpegQuality int) error {
	if format == models.ImageFormatJPEG {
		return jpeg.Encode(dst, src, &jpeg.Options{Quality: jpegQuality})
	}
	return png.Encode(dst, src)
}

func secureTempFile(pattern string) (*os.File, error) {
	file, err := os.CreateTemp("", pattern)
	if err != nil {
		return nil, err
	}
	if err := file.Chmod(0o600); err != nil {
		file.Close()
		os.Remove(file.Name())
		return nil, err
	}
	return file, nil
}

func removeTempFile(file *os.File) {
	name := file.Name()
	_ = file.Close()
	_ = os.Remove(name)
}

type limitWriter struct {
	writer    io.Writer
	remaining int64
}

func (w *limitWriter) Write(data []byte) (int, error) {
	if int64(len(data)) > w.remaining {
		var written int
		if w.remaining > 0 {
			n, err := w.writer.Write(data[:int(w.remaining)])
			written = n
			w.remaining -= int64(n)
			if err != nil {
				return n, err
			}
		}
		return written, ErrImageOutputTooLarge
	}
	n, err := w.writer.Write(data)
	w.remaining -= int64(n)
	return n, err
}

func copyWithContext(ctx context.Context, dst io.Writer, src io.Reader) (int64, error) {
	buffer := make([]byte, 32*1024)
	var total int64
	for {
		if err := ctx.Err(); err != nil {
			return total, err
		}
		n, readErr := src.Read(buffer)
		if n > 0 {
			written, writeErr := dst.Write(buffer[:n])
			total += int64(written)
			if writeErr != nil {
				return total, writeErr
			}
			if written != n {
				return total, io.ErrShortWrite
			}
		}
		if errors.Is(readErr, io.EOF) {
			return total, nil
		}
		if readErr != nil {
			return total, readErr
		}
	}
}
