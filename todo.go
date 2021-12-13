package todo

import "errors"

// Описываем структуры листов, задач и их списков,
// а так же структуры для их обновления.
// Структуры соответствуют базе данных.
// Применяются при чтении запросов клиента.
// Добавлены json-теги для корректного чтения и
// вывода данных при запросах, а также теги для базы данных.
type TodoList struct {
	Id          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title" binding:"required"`
	Description string `json:"description" db:"description"`
}

type UsersList struct {
	Id     int
	UserId int
	ListId int
}

type TodoItem struct {
	Id          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title" binding:"required"`
	Description string `json:"description" db:"description"`
	Done        string `json:"done" db:"done"`
}

type ListsItem struct {
	Id     int
	ListId int
	ItemID int
}

type UpdateListInput struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

func (i UpdateListInput) Validate() error {
	if i.Title == nil && i.Description == nil {
		return errors.New("update structure has no values")
	}

	return nil
}

type UpdateItemInput struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Done        *bool   `json:"done"`
}

func (i UpdateItemInput) Validate() error {
	if i.Title == nil && i.Description == nil && i.Done == nil {
		return errors.New("update structure has no values")
	}

	return nil
}
