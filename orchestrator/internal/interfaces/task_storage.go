package interfaces

import (
	"context"
	"orchestrator/internal/domain"
)

type TaskStorage interface {
	AddOrUpdateTask(ctx context.Context, video *domain.VideoFakeCandidat) error
	GetTask(ctx context.Context, uid string) (*domain.VideoFakeCandidat, error)
}
