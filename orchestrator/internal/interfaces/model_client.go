package interfaces

import "orchestrator/internal/domain"

type ModelClient interface {
	Process(data []byte, format string) (*domain.VideoFakeCandidat, error)
}
