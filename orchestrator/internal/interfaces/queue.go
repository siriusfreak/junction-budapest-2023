package interfaces

import "context"

type Queue interface {
	Add(ctx context.Context, uid string) error
	Pop(ctx context.Context) (string, error)
}
