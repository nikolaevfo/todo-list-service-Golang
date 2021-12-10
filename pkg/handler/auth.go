package handler

import (
	"net/http"
	todo "to-do-list"

	"github.com/gin-gonic/gin"
)

func (h *Handler) signUp(c *gin.Context) {
	var input todo.User

	// байндим входные данные
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// вызываем авторизацию если ошибка, пишем код 500, здесь в ответ должен придти id
	id, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// записываем в ответ статус код 200, если все ок
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type signInInput struct {
	Username string `json: "username" binding:"required"`
	Password string `json: "password" binding:"required"`
}

func (h *Handler) signIn(c *gin.Context) {
	var input signInInput

	// байндим входные данные
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// вызываем авторизацию если ошибка, пишем код 500, здесь в ответ должен придти id
	token, err := h.services.Authorization.GenerateToken(input.Username, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// записываем в ответ статус код 200, если все ок
	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
