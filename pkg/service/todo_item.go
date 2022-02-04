package service

import (
	"WebApi"
	"WebApi/pkg/repository"
)

type TodoItemService struct {
	repo     repository.TodoItem
	listRepo repository.TodoList
}

func NewTodoItemService(repo repository.TodoItem, listRepo repository.TodoList) *TodoItemService {
	return &TodoItemService{
		repo:     repo,
		listRepo: listRepo,
	}
}

func (s *TodoItemService) Create(userId, listId int, item WebApi.TodoItem) (int, error) {
	_, err := s.listRepo.GetById(listId, userId)
	if err != nil {
		return 0, err
	}

	return s.repo.Create(listId, item)
}

func (s *TodoItemService) GetAll(userId, listId int) ([]WebApi.TodoItem, error) {
	return s.repo.GetAll(userId, listId)
}

func (s *TodoItemService) GetById(itemId, userId int) (WebApi.TodoItem, error) {
	return s.repo.GetById(itemId, userId)
}

func (s *TodoItemService) Delete(itemId, userId int) error {
	return s.repo.Delete(itemId, userId)
}

func (s *TodoItemService) Update(itemId, userId int, input WebApi.UpdateTodoItemInput) error {
	if err := input.Validate(); err != nil {
		return err
	}
	return s.repo.Update(itemId, userId, input)
}
