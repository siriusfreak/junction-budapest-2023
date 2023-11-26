package worker

import (
	"context"
	"fmt"
	"log"
	"orchestrator/internal/interfaces"
)

type worker struct {
	name         string
	queue        interfaces.Queue
	videoStorage interfaces.VideoStorage
	tasksStorage interfaces.TaskStorage
	modelClient  interfaces.ModelClient
}

func StartWorker(ctx context.Context, queue interfaces.Queue, videoStorage interfaces.VideoStorage, tasksStorage interfaces.TaskStorage, modelClient interfaces.ModelClient, name string) {
	worker := &worker{
		name:         name,
		queue:        queue,
		tasksStorage: tasksStorage,
		videoStorage: videoStorage,
		modelClient:  modelClient,
	}
	log.Printf("Worker %s started", worker.name)

	for {
		select {
		case <-ctx.Done():
			log.Printf("Worker %s stopped", worker.name)
			return
		default:
			if err := worker.processTask(ctx); err != nil {
				log.Printf("Worker %s encountered an error: %s", worker.name, err)
			}
		}
	}
}

func (w *worker) processTask(ctx context.Context) error {
	uid, err := w.queue.Pop(ctx)
	if err != nil {
		return fmt.Errorf("error worker %v with Pop: %w", w.name, err)
	}
	if uid == "" {
		return nil
	}

	bytes, err := w.videoStorage.GetFile(uid)
	if err != nil {
		return fmt.Errorf("error worker %v with GetFile: %w", w.name, err)
	}

	video, err := w.modelClient.Process(bytes)
	if err != nil {
		return fmt.Errorf("error worker %v with Process: %w", w.name, err)
	}
	
	video.UID = uid

	err = w.tasksStorage.AddOrUpdateTask(ctx, video)
	if err != nil {
		return fmt.Errorf("error worker %v with AddOrUpdateTask: %w", w.name, err)
	}

	return nil
}
