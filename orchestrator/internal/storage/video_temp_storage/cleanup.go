package videotempstorage

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"
)

func (s *Storage) cleanup() error {
	files, err := os.ReadDir(s.directory)
	if err != nil {
		return err
	}

	now := time.Now()
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			log.Print("file info err: %w", err)
		}
		if now.Sub(info.ModTime()) > s.ttl {
			os.Remove(filepath.Join(s.directory, file.Name()))
		}
	}
	return nil
}

func (s *Storage) startCleanupRoutine(ctx context.Context) {
	ticker := time.NewTicker(s.cleanupInterval)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.cleanup()
		}
	}
}
