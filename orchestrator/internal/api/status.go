package api

import (
	"net/http"
	"fmt"
	"orchestrator/internal/usecase"

	"github.com/gin-gonic/gin"
)

func statusHandler(getTaskStatusUseCase *usecase.GetTaskStatusUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := c.Param("uid")

		taskStatusResponse, err := getTaskStatusUseCase.GetTaskStatus(c.Request.Context(), uid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Не удалось получить статус: %v", err)})
			return
		}

		c.JSON(http.StatusOK, taskStatusResponse)
	}
}
