package handler

import (
	"net/http"
	"strconv"
	todo "to-do-list"

	"github.com/gin-gonic/gin"
)

func (h *Handler) createList(c *gin.Context) {
	userId, err := GetUserId(c)
	if err != nil {
		return
	}

	// структурные теги в todo.go позволят создать json, так как в полях структур названия полей должны быть с заглавной буквы, иначе будут неэкспортируемы
	var input todo.TodoList
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.TodoList.Create(userId, input)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type getAllListsResponse struct {
	Data []todo.TodoList `json:"data"`
}

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

func (h *Handler) getListById(c *gin.Context) {
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

	list, err := h.services.TodoList.GetById(userId, listId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// добавляем ответ
	c.JSON(http.StatusOK, list)
}

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
