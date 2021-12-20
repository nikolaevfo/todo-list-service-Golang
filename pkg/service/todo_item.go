package service

import (
	todo "to-do-list"
	"to-do-list/pkg/repository"
)

// Методы сервиса вызывают соответствующие методы из модуля repository,
// передаем данные на уровень ниже

// структура сервиса по работе с задачами
// содержит два репозитория, для связи задач с их списками
type TodoItemService struct {
	repo     repository.TodoItem
	listRepo repository.TodoList
}

// конструктор для создания сервиса по работе с задачами
func newTodoItemService(repo repository.TodoItem, listRepo repository.TodoList) *TodoItemService {
	return &TodoItemService{repo: repo, listRepo: listRepo}
}

func (s *TodoItemService) CreateItem(userId, listId int, item todo.TodoItem) (int, error) {
	// осуществляем проверку на наличие соотвтетствующего списка
	_, err := s.listRepo.GetById(userId, listId)
	if err != nil {
		// если лист не существует
		return 0, err
	}

	return s.repo.CreateItem(listId, item)
}

func (s *TodoItemService) GetAllItems(userId, listId int) ([]todo.TodoItem, error) {
	return s.repo.GetAllItems(userId, listId)
}

func (s *TodoItemService) GetItemById(userId, itemId int) (todo.TodoItem, error) {
	return s.repo.GetItemById(userId, itemId)
}

func (s *TodoItemService) DeleteItem(userId, itemId int) error {
	return s.repo.DeleteItem(userId, itemId)
}

func (s *TodoItemService) UpdateItem(userId, itemId int, input todo.UpdateItemInput) error {
	return s.repo.UpdateItem(userId, itemId, input)
}
