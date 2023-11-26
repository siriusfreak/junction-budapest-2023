package api

import (
	"io"
	"fmt"
	"net/http"
	"orchestrator/internal/usecase"

	"github.com/gin-gonic/gin"
)

func uploadHandler(addTaskUseCase *usecase.AddTaskUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		fileHeader, err := c.FormFile("video")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Не удалось получить видео файл: %v", err)})
			return
		}

		file, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Не удалось открыть видео файл: %v", err)})
			return
		}
		defer file.Close()

		data, err := io.ReadAll(file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Не удалось прочитать видео файл: %v", err)})
			return
		}

		uid, err := addTaskUseCase.AddTask(c.Request.Context(), data)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Не удалось добавить задачу: %v", err)})
			return
		}

		c.JSON(http.StatusOK, gin.H{"uid": uid})
	}
}
