package service

import (
	todo "to-do-list"
	"to-do-list/pkg/repository"
)

// Интерфейсы называем исходя из их доменной зоны бизнес-логики.
// Создаем объединяющую структуру Service, и объявим ее конструктор.
// Сервисы обращаются к базе данных, поэтому конструктор принимает указатель на структуру repository.
// Это и есть внедрение зависимостей

type Authorization interface {
	// метод принимает данные о пользователе, возвращает id и ошибку
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

// описываем струтуру сервиса, состоящую из интерфейсов
type Service struct {
	Authorization
	TodoList
	TodoItem
}

// конструктор сервиса, в котором инициализируются сервисы авторизации,
// сервисы работы со списками и задачами
// данные уходят на слой ниже, в repository
func NewService(repos *repository.Repository) *Service {
	// инициализация сервиса
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		TodoList:      NewTodoListSevice(repos.TodoList),
		TodoItem:      newTodoItemService(repos.TodoItem, repos.TodoList),
	}
}
