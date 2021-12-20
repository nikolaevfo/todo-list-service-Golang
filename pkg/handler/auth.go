package handler

import (
	"net/http"
	todo "to-do-list"

	"github.com/gin-gonic/gin"
)

// обработчики построены по одному принципу в следующем порядке:
// приводим id пользователя из контекста в типу Int методом GetUserId
// байндим json в соответствующую структуру из todo.go
// вызываем метод сервиса, передаем в него полученную структуру с данными запроса
// записываем в ответ полученные данные из сервиса

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
	// структура инпута для парсинга из json
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

// структура для парсинга тела запроса из json
type signInInput struct {
	Username string `json: "username" binding:"required"`
	Password string `json: "password" binding:"required"`
}

// метод signIn будет возвращать токен, если пользоватьель найден в БД
//
// описываем данные для swagger
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

	// вызываем метод создания токена, если ошибка, пишем код 500
	token, err := h.services.Authorization.GenerateToken(input.Username, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// записываем в ответ статус код 200, если все ок, и токен
	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
