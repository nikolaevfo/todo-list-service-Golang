package repository

import (
	todo "to-do-list"

	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user todo.User) (int, error)
	GetUser(username, password string) (todo.User, error)
}
type TodoList interface {
	Create(userId int, list todo.TodoList) (int, error)
	GetAll(userId int) ([]todo.TodoList, error)
	GetById(userId, listId int) (todo.TodoList, error)
	DeleteList(userId, listId int) error
	UpdateList(userId, listId int, input todo.UpdateListInput) error
}
type TodoItem interface {
	CreateItem(listId int, item todo.TodoItem) (int, error)
	GetAllItems(userId, listId int) ([]todo.TodoItem, error)
	GetItemById(userId, itemId int) (todo.TodoItem, error)
	UpdateItem(userId, itemId int, input todo.UpdateItemInput) error
	DeleteItem(userId, itemId int) error
}

type Repository struct {
	Authorization
	TodoList
	TodoItem
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		TodoList:      NewTodoListPostgres(db),
		TodoItem:      NewTodoItemPostgres(db),
	}
}
