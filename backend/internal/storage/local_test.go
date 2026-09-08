package storage

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewKeyReturnsIndependentSafeKeys(t *testing.T) {
	t.Parallel()

	first, err := NewKey()
	if err != nil {
		t.Fatalf("NewKey() error = %v", err)
	}
	second, err := NewKey()
	if err != nil {
		t.Fatalf("NewKey() error = %v", err)
	}
	if len(first) != storageKeyBytes*2 || len(second) != storageKeyBytes*2 {
		t.Fatalf("NewKey() lengths = %d, %d", len(first), len(second))
	}
	if first == second {
		t.Fatal("NewKey() returned the same key twice")
	}
}

func TestLocalPersistsAcrossInstances(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	key, err := NewKey()
	if err != nil {
		t.Fatalf("NewKey() error = %v", err)
	}
	first, err := NewLocal(root)
	if err != nil {
		t.Fatalf("NewLocal() error = %v", err)
	}
	want := []byte("sanitized image bytes")
	if err := first.Put(context.Background(), key, bytes.NewReader(want)); err != nil {
		t.Fatalf("Put() error = %v", err)
	}

	second, err := NewLocal(root)
	if err != nil {
		t.Fatalf("NewLocal() second instance error = %v", err)
	}
	reader, err := second.Open(context.Background(), key)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer reader.Close()
	got, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("read stored attachment: %v", err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("stored attachment = %q, want %q", got, want)
	}
}

func TestLocalRejectsUnsafeKeys(t *testing.T) {
	t.Parallel()

	store, err := NewLocal(t.TempDir())
	if err != nil {
		t.Fatalf("NewLocal() error = %v", err)
	}
	unsafeKeys := []string{"", "../secret", strings.Repeat("g", storageKeyBytes*2), strings.Repeat("a", storageKeyBytes*2-1)}
	for _, key := range unsafeKeys {
		if err := store.Put(context.Background(), key, bytes.NewReader(nil)); !errors.Is(err, ErrInvalidStorageKey) {
			t.Fatalf("Put(%q) error = %v, want %v", key, err, ErrInvalidStorageKey)
		}
	}
}

func TestLocalDoesNotOverwriteExistingKey(t *testing.T) {
	t.Parallel()

	store, err := NewLocal(t.TempDir())
	if err != nil {
		t.Fatalf("NewLocal() error = %v", err)
	}
	key, err := NewKey()
	if err != nil {
		t.Fatalf("NewKey() error = %v", err)
	}
	if err := store.Put(context.Background(), key, strings.NewReader("first")); err != nil {
		t.Fatalf("Put() first error = %v", err)
	}
	if err := store.Put(context.Background(), key, strings.NewReader("second")); !errors.Is(err, ErrStorageKeyExists) {
		t.Fatalf("Put() second error = %v, want %v", err, ErrStorageKeyExists)
	}
}

func TestLocalDeleteIsIdempotent(t *testing.T) {
	t.Parallel()

	store, err := NewLocal(t.TempDir())
	if err != nil {
		t.Fatalf("NewLocal() error = %v", err)
	}
	key, err := NewKey()
	if err != nil {
		t.Fatalf("NewKey() error = %v", err)
	}
	if err := store.Put(context.Background(), key, strings.NewReader("data")); err != nil {
		t.Fatalf("Put() error = %v", err)
	}
	if err := store.Delete(context.Background(), key); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	if err := store.Delete(context.Background(), key); err != nil {
		t.Fatalf("Delete() repeated error = %v", err)
	}
}

func TestLocalRemovesTemporaryFileAfterCanceledWrite(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	store, err := NewLocal(root)
	if err != nil {
		t.Fatalf("NewLocal() error = %v", err)
	}
	key, err := NewKey()
	if err != nil {
		t.Fatalf("NewKey() error = %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := store.Put(ctx, key, strings.NewReader("data")); !errors.Is(err, context.Canceled) {
		t.Fatalf("Put() error = %v, want %v", err, context.Canceled)
	}

	directory := filepath.Join(root, key[:2], key[2:4])
	entries, err := os.ReadDir(directory)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("ReadDir() error = %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("temporary storage directory contains %d files after failed write", len(entries))
	}
}
