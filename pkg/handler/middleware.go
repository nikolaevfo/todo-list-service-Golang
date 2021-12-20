package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

// метод мидлвары для идентификация user и записи его id в контекст
func (h *Handler) userIdentity(c *gin.Context) {
	// получаем значение из заголовка авторизации
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		// возвращаем статус 401, пользователь не авторизован
		newErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		return
	}

	// разделяем строку по пробелам, должен вернуться массив из 2 элементов
	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		// возвращаем статус 401, пользователь не авторизован
		newErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
		return
	}

	// используем функцию сервиса для парсинга токена и получения UserId
	UserId, err := h.services.Authorization.ParseToken(headerParts[1])
	if err != nil {
		// возвращаем статус 401, пользователь не авторизован
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	// записываем UserId в контекст, для проверки в последующих обработчиках
	c.Set(userCtx, UserId)
}

// метод для приведения id пользователя к типу Int
func GetUserId(c *gin.Context) (int, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "user id not found")
		return 0, errors.New("user id not found")
	}

	idInt, ok := id.(int)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "user id id of invalid type")
		return 0, errors.New("user id id of invalid type")
	}
	return idInt, nil
}
