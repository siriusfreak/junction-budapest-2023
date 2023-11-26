package main

import (
	"context"
	"orchestrator/internal/api"
	"orchestrator/internal/integrations/models"
	"orchestrator/internal/integrations/redis"
	"orchestrator/internal/interfaces"
	"orchestrator/internal/queue"
	tasksstorage "orchestrator/internal/storage/tasks_storage"
	videotempstorage "orchestrator/internal/storage/video_temp_storage"
	"orchestrator/internal/usecase"
	"orchestrator/internal/worker"
	"time"
)

func main() {
	ctx := context.Background()
	cfg := getConfig()

	videoStorage := videotempstorage.NewStorage(ctx, cfg.VideoStorageDir, cfg.VideoStorageTTL, cfg.VideoStorageCleanupInterval)

	OnePersonDetectClient := models.NewModelClientImpl(models.OnePersonDetectURL, models.OnePersonDetectProcess)

	redisClient := redis.CreateRedisClient(ctx, cfg.RedisAddr, cfg.RedisPass, 0)

	tasksStorage := tasksstorage.NewStorage(ctx, redisClient)

	queueOnePersonDetect := queue.NewQueue(ctx, "one-person-detect", redisClient)

	// queueAudioFakeDetection := queue.NewQueue(ctx, "audio-fake-detection", redisClient)
	// queueDeepfakeDetect := queue.NewQueue(ctx, "deepfake-detect", redisClient)
	// queueLipsMovementDetection := queue.NewQueue(ctx, "lips-movement-detection", redisClient)
	// queueOpenClosedEyeDetect := queue.NewQueue(ctx, "open-closed-eye-detect", redisClient)
	// queueWhisperLargeV3 := queue.NewQueue(ctx, "whisper-large-v3", redisClient)

	go worker.StartWorker(ctx, queueOnePersonDetect, videoStorage, tasksStorage, OnePersonDetectClient, "one-person-detect")

	// go worker.StartWorker(ctx, queueAudioFakeDetection, videoStorage, tasksStorage, nil, "audio-fake-detection")
	// go worker.StartWorker(ctx, queueDeepfakeDetect, videoStorage, tasksStorage, nil, "deepfake-detect")
	// go worker.StartWorker(ctx, queueLipsMovementDetection, videoStorage, tasksStorage, nil, "lips-movement-detection")
	// go worker.StartWorker(ctx, queueOpenClosedEyeDetect, videoStorage, tasksStorage, nil, "open-closed-eye-detect")
	// go worker.StartWorker(ctx, queueWhisperLargeV3, videoStorage, tasksStorage, nil, "whisper-large-v3")
	addTaskUseCase := usecase.NewAddTaskUseCase(
		videoStorage,
		tasksStorage,
		[]interfaces.Queue{
			queueOnePersonDetect,
		},
	)

	getTaskStatusUseCase := usecase.NewGetTaskStatusUseCase(tasksStorage)

	router := api.NewRouter(addTaskUseCase, getTaskStatusUseCase)

	router.Run(":8888")
}

type config struct {
	RedisAddr                   string
	RedisPass                   string
	VideoStorageDir             string
	VideoStorageTTL             time.Duration
	VideoStorageCleanupInterval time.Duration
}

func getConfig() config {
	return config{
		RedisAddr:                   "localhost:6379",
		RedisPass:                   "",
		VideoStorageDir:             "./videos",
		VideoStorageTTL:             12 * time.Hour,
		VideoStorageCleanupInterval: 10 * time.Minute,
	}
}
