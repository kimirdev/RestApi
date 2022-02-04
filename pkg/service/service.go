package service

import (
	"WebApi"
	"WebApi/pkg/repository"
)

type Authorization interface {
	CreateUser(user WebApi.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(accessToken string) (int, error)
}

type TodoList interface {
	Create(userId int, list WebApi.TodoList) (int, error)
	GetAll(userId int) ([]WebApi.TodoList, error)
	GetById(listId, userId int) (WebApi.TodoList, error)
	Delete(listId, userId int) error
	Update(listId, userId int, input WebApi.UpdateTodoListInput) error
}

type TodoItem interface {
	Create(userId, listId int, item WebApi.TodoItem) (int, error)
	GetAll(userId, listId int) ([]WebApi.TodoItem, error)
	GetById(itemId, userId int) (WebApi.TodoItem, error)
	Delete(itemId, userId int) error
	Update(itemId, userId int, input WebApi.UpdateTodoItemInput) error
}

type Service struct {
	Authorization
	TodoList
	TodoItem
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		TodoList:      NewTodoListService(repos.TodoList),
		TodoItem:      NewTodoItemService(repos.TodoItem, repos.TodoList),
	}
}
