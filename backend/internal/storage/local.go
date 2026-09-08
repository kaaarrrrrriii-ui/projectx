package storage

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const storageKeyBytes = 32

var (
	ErrInvalidStorageKey = errors.New("invalid attachment storage key")
	ErrStorageKeyExists  = errors.New("attachment storage key already exists")
)

// NewKey returns a cryptographically random 256-bit storage key.
func NewKey() (string, error) {
	random := make([]byte, storageKeyBytes)
	if _, err := rand.Read(random); err != nil {
		return "", fmt.Errorf("generate attachment storage key: %w", err)
	}
	return hex.EncodeToString(random), nil
}

// Local stores attachments below a private filesystem root. The configured
// root must be mounted on a persistent volume by the application environment.
type Local struct {
	root string
}

func NewLocal(root string) (*Local, error) {
	if root == "" {
		return nil, errors.New("attachment storage root is required")
	}
	absoluteRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, fmt.Errorf("resolve attachment storage root: %w", err)
	}
	if err := os.MkdirAll(absoluteRoot, 0o700); err != nil {
		return nil, fmt.Errorf("create attachment storage root: %w", err)
	}
	return &Local{root: absoluteRoot}, nil
}

func (s *Local) Put(ctx context.Context, key string, source io.Reader) error {
	target, err := s.pathForKey(key)
	if err != nil {
		return err
	}
	directory := filepath.Dir(target)
	if err := os.MkdirAll(directory, 0o700); err != nil {
		return fmt.Errorf("create attachment storage directory: %w", err)
	}
	if _, err := os.Lstat(target); err == nil {
		return ErrStorageKeyExists
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("inspect attachment storage target: %w", err)
	}

	temporary, err := os.CreateTemp(directory, ".attachment-*")
	if err != nil {
		return fmt.Errorf("create attachment storage temporary file: %w", err)
	}
	temporaryName := temporary.Name()
	committed := false
	defer func() {
		_ = temporary.Close()
		if !committed {
			_ = os.Remove(temporaryName)
		}
	}()

	if err := temporary.Chmod(0o600); err != nil {
		return fmt.Errorf("secure attachment storage file: %w", err)
	}
	if _, err := copyWithContext(ctx, temporary, source); err != nil {
		return fmt.Errorf("write attachment storage file: %w", err)
	}
	if err := temporary.Sync(); err != nil {
		return fmt.Errorf("sync attachment storage file: %w", err)
	}
	if err := temporary.Close(); err != nil {
		return fmt.Errorf("close attachment storage file: %w", err)
	}
	if _, err := os.Lstat(target); err == nil {
		return ErrStorageKeyExists
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("inspect attachment storage target: %w", err)
	}
	if err := os.Rename(temporaryName, target); err != nil {
		return fmt.Errorf("publish attachment storage file: %w", err)
	}
	committed = true
	return nil
}

func (s *Local) Open(_ context.Context, key string) (io.ReadCloser, error) {
	target, err := s.pathForKey(key)
	if err != nil {
		return nil, err
	}
	file, err := os.Open(target)
	if err != nil {
		return nil, fmt.Errorf("open attachment storage file: %w", err)
	}
	return file, nil
}

// Delete is idempotent so cleanup can safely resume after a restart.
func (s *Local) Delete(_ context.Context, key string) error {
	target, err := s.pathForKey(key)
	if err != nil {
		return err
	}
	if err := os.Remove(target); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("delete attachment storage file: %w", err)
	}
	return nil
}

func (s *Local) pathForKey(key string) (string, error) {
	if len(key) != storageKeyBytes*2 {
		return "", ErrInvalidStorageKey
	}
	if _, err := hex.DecodeString(key); err != nil {
		return "", ErrInvalidStorageKey
	}
	return filepath.Join(s.root, key[:2], key[2:4], key), nil
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
