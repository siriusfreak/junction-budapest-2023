package api

import (
	"io"
	"net/http"
	"orchestrator/internal/usecase"

	"github.com/gin-gonic/gin"
)

func uploadHandler(addTaskUseCase *usecase.AddTaskUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		fileHeader, err := c.FormFile("video")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Не удалось получить видео файл"})
			return
		}

		file, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось открыть видео файл"})
			return
		}
		defer file.Close()

		data, err := io.ReadAll(file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось прочитать видео файл"})
			return
		}

		uid, err := addTaskUseCase.AddTask(c.Request.Context(), data)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось добавить задачу"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"uid": uid})
	}
}
