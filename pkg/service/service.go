package service

import (
	todo "to-do-list"
	"to-do-list/pkg/repository"
)

type Authorization interface {
	CreateUser(user todo.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}
type TodoList interface {
	Create(userId int, list todo.TodoList) (int, error)
	GetAll(userId int) ([]todo.TodoList, error)
	GetById(userId, listId int) (todo.TodoList, error)
	DeleteList(userId, listId int) error
	UpdateList(userId, listId int, input todo.UpdateListInput) error
}
type TodoItem interface {
	CreateItem(userId, listId int, input todo.TodoItem) (int, error)
	GetAllItems(userId, listId int) ([]todo.TodoItem, error)
	GetItemById(userId, itemId int) (todo.TodoItem, error)
	UpdateItem(userId, itemId int, input todo.UpdateItemInput) error
	DeleteItem(userId, itemId int) error
}

type Service struct {
	Authorization
	TodoList
	TodoItem
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		// инициализация новых репозиториев
		Authorization: NewAuthService(repos.Authorization),
		TodoList:      NewTodoListSevice(repos.TodoList),
		TodoItem:      newTodoItemService(repos.TodoItem, repos.TodoList),
	}
}
