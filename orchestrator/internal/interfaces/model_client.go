package interfaces

import "orchestrator/internal/domain"

type ModelClient interface {
	Process(data []byte) (*domain.VideoFakeCandidat, error)
}
