// Package storage contains private attachment storage adapters.
package storage

import (
	"context"
	"io"
)

// Storage persists sanitized attachments. Keys are internal identifiers and
// must never be exposed as public URLs.
type Storage interface {
	Put(ctx context.Context, key string, source io.Reader) error
	Open(ctx context.Context, key string) (io.ReadCloser, error)
	Delete(ctx context.Context, key string) error
}
