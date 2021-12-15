package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// объявляем структуру ошибки
type errorResponse struct {
	Message string `json:"message"`
}

// структура статуса ответа
type statusResponse struct {
	Status string `json:"status"`
}

// функция для обработки ошибок
func newErrorResponse(c *gin.Context, statusCode int, message string) {
	logrus.Error(message)

	// метод блокирует выполнение следующих обработчик и записывает в ответ статус и сообщение
	c.AbortWithStatusJSON(statusCode, errorResponse{message})
}
