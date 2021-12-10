package service

import (
	todo "to-do-list"
	"to-do-list/pkg/repository"
)

type TodoListService struct {
	repo repository.TodoList
}

func NewTodoListSevice(repo repository.TodoList) *TodoListService {
	return &TodoListService{repo: repo}
}

func (s *TodoListService) Create(userId int, list todo.TodoList) (int, error) {
	return s.repo.Create(userId, list)
}

func (s *TodoListService) GetAll(userId int) ([]todo.TodoList, error) {
	return s.repo.GetAll(userId)
}

func (s *TodoListService) GetById(userId, listId int) (todo.TodoList, error) {
	return s.repo.GetById(userId, listId)
}

func (s *TodoListService) DeleteList(userId, listId int) error {
	return s.repo.DeleteList(userId, listId)
}

func (s *TodoListService) UpdateList(userId, listId int, input todo.UpdateListInput) error {
	if err := input.Validate(); err != nil {
		return err
	}
	return s.repo.UpdateList(userId, listId, input)
}
