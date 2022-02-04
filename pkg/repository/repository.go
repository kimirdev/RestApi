package repository

import (
	"WebApi"

	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user WebApi.User) (int, error)
	GetUser(username, password string) (WebApi.User, error)
}

type TodoList interface {
	Create(userId int, list WebApi.TodoList) (int, error)
	GetAll(userId int) ([]WebApi.TodoList, error)
	GetById(listId, userId int) (WebApi.TodoList, error)
	Delete(listId, userId int) error
	Update(listId, userId int, input WebApi.UpdateTodoListInput) error
}

type TodoItem interface {
	Create(listId int, item WebApi.TodoItem) (int, error)
	GetAll(userId, listId int) ([]WebApi.TodoItem, error)
	GetById(itemId, userId int) (WebApi.TodoItem, error)
	Delete(itemId, userId int) error
	Update(itemId, userId int, input WebApi.UpdateTodoItemInput) error
}

type Repository struct {
	Authorization
	TodoList
	TodoItem
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthSql(db),
		TodoList:      NewTodoListSql(db),
		TodoItem:      NewTodoItemSql(db),
	}
}
