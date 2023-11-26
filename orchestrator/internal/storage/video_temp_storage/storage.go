package videotempstorage

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

type Storage struct {
	directory       string
	ttl             time.Duration
	cleanupInterval time.Duration
}

func NewStorage(ctx context.Context, dir string, ttl, cleanupInterval time.Duration) *Storage {
	storage := &Storage{
		directory:       dir,
		ttl:             ttl,
		cleanupInterval: cleanupInterval,
	}
	go storage.startCleanupRoutine(ctx)
	return storage
}

func (s *Storage) SaveFile(data []byte) (string, error) {
	uid := uuid.New().String()
	filePath := filepath.Join(s.directory, uid)
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return "", err
	}
	return uid, nil
}

func (s *Storage) GetFile(uid string) ([]byte, error) {
	filePath := filepath.Join(s.directory, uid)
	return os.ReadFile(filePath)
}
