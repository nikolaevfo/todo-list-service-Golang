package handler

import (
	"net/http"
	todo "to-do-list"

	"github.com/gin-gonic/gin"
)

// описываем данные для swagger
// @Summary      signUp
// @Description  create account
// @Tags         auth
// ID create-account
// @Accept       json
// @Produce      json
// @Param        input body todo.User true "account info"
// @Success      200  {integer}  integer "id"
// @Failure      400,404  {object}  errorResponse
// @Failure      500  {object}  errorResponse
// @Failure      default  {object}  errorResponse
// @Router       /auth/sign-up [post]
func (h *Handler) signUp(c *gin.Context) {
	var input todo.User

	// байндим входные данные в соответствии со структурой todo.User
	if err := c.BindJSON(&input); err != nil {
		// вызываем созданную нами функцию обработки ошибок (код статуса 400)
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// передаем данные на слой ниже, в сервис
	// вызываем сервис авторизации
	// передаем тело запрсоа. в ответ должен придти id
	id, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		// обрабатываем ошибку, возвращаем код 500.
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// записываем в ответ статус код 200, если все ок
	// и тело json со значением id пользователя
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type signInInput struct {
	Username string `json: "username" binding:"required"`
	Password string `json: "password" binding:"required"`
}

// @Summary      signIn
// @Description  login
// @Tags         auth
// ID login
// @Accept       json
// @Produce      json
// @Param        input body signInInput true "credential"
// @Success      200  {string}  string "token"
// @Failure      400,404  {object}  errorResponse
// @Failure      500  {object}  errorResponse
// @Failure      default  {object}  errorResponse
// @Router       /auth/sign-in [post]
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
