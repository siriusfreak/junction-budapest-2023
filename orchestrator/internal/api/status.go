package api

import (
	"net/http"
	"orchestrator/internal/usecase"

	"github.com/gin-gonic/gin"
)

func statusHandler(getTaskStatusUseCase *usecase.GetTaskStatusUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := c.Param("uid")

		taskStatusResponse, err := getTaskStatusUseCase.GetTaskStatus(c.Request.Context(), uid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить статус"})
			return
		}

		c.JSON(http.StatusOK, taskStatusResponse)
	}
}
