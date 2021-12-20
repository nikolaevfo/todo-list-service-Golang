package handler

import (
	"net/http"
	"strconv"
	todo "to-do-list"

	"github.com/gin-gonic/gin"
)

// обработчики построены по одному принципу в следующем порядке:
// приводим id пользователя из контекста в типу Int методом GetUserId
// байндим json в соответствующую структуру из todo.go
// вызываем метод сервиса, передаем в него полученную структуру с данными запроса
// записываем в ответ полученные данные из сервиса

// описываем данные для swagger
// @Summary      Create Todo List
// @Security ApiKeyAuth
// @Description  create todo list
// @Tags         lists
// ID create-list
// @Accept       json
// @Produce      json
// @Param        input body todo.TodoList true "list data"
// @Success      200  {integer}  integer "id"
// @Failure      400,404  {object}  errorResponse
// @Failure      500  {object}  errorResponse
// @Failure      default  {object}  errorResponse
// @Router       /api/lists [post]
func (h *Handler) createList(c *gin.Context) {
	userId, err := GetUserId(c)
	if err != nil {
		return
	}

	// структурные теги в todo.go позволят создать json,
	// так как в полях структур названия полей должны быть с заглавной буквы,
	// иначе будут неэкспортируемы
	var input todo.TodoList
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.TodoList.Create(userId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// дополнительная структура для ответа
type getAllListsResponse struct {
	Data []todo.TodoList `json:"data"`
}

// описываем данные для swagger
// @Summary      Get All Lists
// @Security ApiKeyAuth
// @Description  get all lists
// @Tags         lists
// ID get-all-lists
// @Accept       json
// @Produce      json
// @Success      200  {object}  getAllListsResponse
// @Failure      400,404  {object}  errorResponse
// @Failure      500  {object}  errorResponse
// @Failure      default  {object}  errorResponse
// @Router       /api/lists [get]
func (h *Handler) getAllLists(c *gin.Context) {
	userId, err := GetUserId(c)
	if err != nil {
		return
	}

	lists, err := h.services.TodoList.GetAll(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// добавляем ответ
	c.JSON(http.StatusOK, getAllListsResponse{
		Data: lists,
	})
}

// описываем данные для swagger
// @Summary      Get List By Id
// @Security ApiKeyAuth
// @Description  get list by id
// @Tags         lists
// ID get-list-by-id
// @Accept       json
// @Produce      json
// @Success      200  {object}  todo.ListsItem
// @Failure      400,404  {object}  errorResponse
// @Failure      500  {object}  errorResponse
// @Failure      default  {object}  errorResponse
// @Router       /api/lists/:id [get]
func (h *Handler) getListById(c *gin.Context) {
	userId, err := GetUserId(c)
	if err != nil {
		return
	}

	// получаем динамический id листа из строки запроса
	// с помощью метода Param
	// strconv.Atoi преобразует строку в число
	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	list, err := h.services.TodoList.GetById(userId, listId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// добавляем ответ
	c.JSON(http.StatusOK, list)
}

// описываем данные для swagger
// @Summary      Update List
// @Security ApiKeyAuth
// @Description  update list
// @Tags         lists
// ID update-list
// @Accept       json
// @Produce      json
// @Success      200  {string}  string
// @Failure      400,404  {object}  errorResponse
// @Failure      500  {object}  errorResponse
// @Failure      default  {object}  errorResponse
// @Router       /api/lists/:id [put]
func (h *Handler) updateList(c *gin.Context) {
	userId, err := GetUserId(c)
	if err != nil {
		return
	}

	// получаем id листа из строки запроса
	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	var input todo.UpdateListInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.TodoList.UpdateList(userId, listId, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// добавляем ответ
	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}

// описываем данные для swagger
// @Summary      Delete List
// @Security ApiKeyAuth
// @Description  delete list
// @Tags         lists
// ID delete-list
// @Accept       json
// @Produce      json
// @Success      200  {string}  string
// @Failure      400,404  {object}  errorResponse
// @Failure      500  {object}  errorResponse
// @Failure      default  {object}  errorResponse
// @Router       /api/lists/:id [delete]
func (h *Handler) deleteList(c *gin.Context) {
	userId, err := GetUserId(c)
	if err != nil {
		return
	}

	// получаем id листа из строки запроса
	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	err = h.services.TodoList.DeleteList(userId, listId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// добавляем ответ
	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}
