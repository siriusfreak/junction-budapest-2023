package api

import (
	"orchestrator/internal/usecase"

	"github.com/gin-gonic/gin"
)

func NewRouter(addTaskUseCase *usecase.AddTaskUseCase, getTaskStatusUseCase *usecase.GetTaskStatusUseCase) *gin.Engine {
	router := gin.Default()

	router.POST("/upload", uploadHandler(addTaskUseCase))

	router.GET("/status/:uid", statusHandler(getTaskStatusUseCase))

	return router
}
