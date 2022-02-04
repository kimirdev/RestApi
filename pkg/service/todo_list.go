package service

import (
	"WebApi"
	"WebApi/pkg/repository"
)

type TodoListService struct {
	repo repository.TodoList
}

func NewTodoListService(repo repository.TodoList) *TodoListService {
	return &TodoListService{repo: repo}
}

func (s *TodoListService) Create(userId int, list WebApi.TodoList) (int, error) {
	return s.repo.Create(userId, list)
}

func (s *TodoListService) GetAll(userId int) ([]WebApi.TodoList, error) {
	return s.repo.GetAll(userId)
}

func (s *TodoListService) GetById(listId, userId int) (WebApi.TodoList, error) {
	return s.repo.GetById(listId, userId)
}

func (s *TodoListService) Delete(listId, userId int) error {
	return s.repo.Delete(listId, userId)
}

func (s *TodoListService) Update(listId, userId int, input WebApi.UpdateTodoListInput) error {
	if err := input.Validate(); err != nil {
		return err
	}
	return s.repo.Update(listId, userId, input)
}
