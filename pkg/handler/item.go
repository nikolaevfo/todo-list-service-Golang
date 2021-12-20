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
// @Summary      Create Item
// @Security ApiKeyAuth
// @Description  create item
// @Tags         lists
// ID create-item
// @Accept       json
// @Produce      json
// @Param        input body todo.TodoItem true "item data"
// @Success      200  {integer}  integer "id"
// @Failure      400,404  {object}  errorResponse
// @Failure      500  {object}  errorResponse
// @Failure      default  {object}  errorResponse
// @Router       /api/lists/:id/items [post]
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

// описываем данные для swagger
// @Summary      Get All Items
// @Security ApiKeyAuth
// @Description  get all items
// @Tags         lists
// ID get-all-items
// @Accept       json
// @Produce      json
// @Success      200  {object}  []todo.TodoItem
// @Failure      400,404  {object}  errorResponse
// @Failure      500  {object}  errorResponse
// @Failure      default  {object}  errorResponse
// @Router       /api/lists/:id/items [get]
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

// описываем данные для swagger
// @Summary      Get Item By Id
// @Security ApiKeyAuth
// @Description  get item by id
// @Tags         items
// ID get-item-by-id
// @Accept       json
// @Produce      json
// @Success      200  {object}  todo.TodoItem
// @Failure      400,404  {object}  errorResponse
// @Failure      500  {object}  errorResponse
// @Failure      default  {object}  errorResponse
// @Router       /api/items/:id [get]
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

// описываем данные для swagger
// @Summary      Update Item
// @Security ApiKeyAuth
// @Description  update item
// @Tags         items
// ID update-item
// @Accept       json
// @Produce      json
// @Success      200  {string}  string
// @Failure      400,404  {object}  errorResponse
// @Failure      500  {object}  errorResponse
// @Failure      default  {object}  errorResponse
// @Router       /api/items/:id [put]
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

// описываем данные для swagger
// @Summary      Delete Item
// @Security ApiKeyAuth
// @Description  delete item
// @Tags         items
// ID delete-item
// @Accept       json
// @Produce      json
// @Success      200  {string}  string
// @Failure      400,404  {object}  errorResponse
// @Failure      500  {object}  errorResponse
// @Failure      default  {object}  errorResponse
// @Router       /api/items/:id [delete]
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
