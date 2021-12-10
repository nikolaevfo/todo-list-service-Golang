package handler

import (
	"net/http"
	"strconv"
	todo "to-do-list"

	"github.com/gin-gonic/gin"
)

func (h *Handler) createItem(c *gin.Context) {
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

	var input todo.TodoItem
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.TodoItem.CreateItem(userId, listId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) getAllItems(c *gin.Context) {
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

	items, err := h.services.TodoItem.GetAllItems(userId, listId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, items)
}

func (h *Handler) getItemById(c *gin.Context) {
	userId, err := GetUserId(c)
	if err != nil {
		return
	}

	// получаем id листа из строки запроса
	itemId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	item, err := h.services.TodoItem.GetItemById(userId, itemId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *Handler) updateItem(c *gin.Context) {
	userId, err := GetUserId(c)
	if err != nil {
		return
	}

	// получаем id листа из строки запроса
	itemId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	var input todo.UpdateItemInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.TodoItem.UpdateItem(userId, itemId, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// добавляем ответ
	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}

func (h *Handler) deleteItem(c *gin.Context) {
	userId, err := GetUserId(c)
	if err != nil {
		return
	}

	// получаем id листа из строки запроса
	itemId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	err = h.services.TodoItem.DeleteItem(userId, itemId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
}
