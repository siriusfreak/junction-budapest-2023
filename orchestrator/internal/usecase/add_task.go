package usecase

import (
	"context"
	"orchestrator/internal/domain"
	"orchestrator/internal/interfaces"
)

type AddTaskUseCase struct {
	videoStorage interfaces.VideoStorage
	taskStorage  interfaces.TaskStorage
	queues       []interfaces.Queue
}

func NewAddTaskUseCase(videoStorage interfaces.VideoStorage, taskStorage interfaces.TaskStorage, queues []interfaces.Queue) *AddTaskUseCase {
	return &AddTaskUseCase{
		videoStorage: videoStorage,
		taskStorage:  taskStorage,
		queues:       queues,
	}
}

func (uc *AddTaskUseCase) AddTask(ctx context.Context, video []byte) (string, error) {
	uid, err := uc.videoStorage.SaveFile(video)
	if err != nil {
		return "", err
	}

	err = uc.taskStorage.AddOrUpdateTask(ctx, &domain.VideoFakeCandidat{
		UID: uid,
	})
	if err != nil {
		return "", err
	}

	for _, queue := range uc.queues {
		err = queue.Add(ctx, uid)
		if err != nil {
			return "", err
		}
	}

	return uid, nil
}
