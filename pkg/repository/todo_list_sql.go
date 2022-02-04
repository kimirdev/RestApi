package repository

import (
	"WebApi"
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
)

type TodoListSql struct {
	db *sqlx.DB
}

func NewTodoListSql(db *sqlx.DB) *TodoListSql {
	return &TodoListSql{db: db}
}

func (s *TodoListSql) Create(userId int, list WebApi.TodoList) (int, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return 0, err
	}

	createListQueqry := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", todoListsTable)
	createUsersListsQuery := fmt.Sprintf("INSERT INTO %s (user_id, list_id) VALUES($1, $2)", usersListsTable)

	row := tx.QueryRow(createListQueqry, list.Title, list.Description)
	var listId int
	if err := row.Scan(&listId); err != nil {
		tx.Rollback()
		return 0, err
	}

	_, err = tx.Exec(createUsersListsQuery, userId, listId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return listId, tx.Commit()
}

func (s *TodoListSql) GetAll(userId int) ([]WebApi.TodoList, error) {
	var lists []WebApi.TodoList
	query := fmt.Sprintf(
		`SELECT tl.id, tl.title, tl.description 
		FROM %s tl 
		INNER JOIN %s ul ON tl.id = ul.list_id 
		WHERE ul.user_id = $1`,
		todoListsTable, usersListsTable)

	if err := s.db.Select(&lists, query, userId); err != nil {
		return nil, err
	}

	return lists, nil
}

func (s *TodoListSql) GetById(listId, userId int) (WebApi.TodoList, error) {
	var list WebApi.TodoList
	query := fmt.Sprintf(
		`SELECT tl.id, tl.title, tl.description 
		FROM %s tl 
		INNER JOIN %s ul ON tl.id = ul.list_id 
		WHERE ul.user_id = $1 AND tl.id = $2`,
		todoListsTable, usersListsTable)

	err := s.db.Get(&list, query, userId, listId)

	return list, err
}

func (s *TodoListSql) Delete(listId, userId int) error {
	query := fmt.Sprintf(
		`DELETE FROM %s tl USING %s ul 
		WHERE tl.id = ul.list_id 
		AND ul.user_id = $1 
		AND ul.list_id = $2`,
		todoListsTable, usersListsTable)

	_, err := s.db.Exec(query, userId, listId)

	return err
}

func (s *TodoListSql) Update(listId, userId int, input WebApi.UpdateTodoListInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprint("title=$", argId))
		args = append(args, *input.Title)
		argId++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprint("description=$", argId))
		args = append(args, *input.Description)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(
		`UPDATE %s tl SET %s 
		FROM %s ul 
		WHERE tl.id = ul.list_id 
		AND ul.user_id = $%d 
		AND ul.list_id = $%d`,
		todoListsTable, setQuery, usersListsTable, argId, argId+1)

	args = append(args, userId, listId)
	_, err := s.db.Exec(query, args...)

	log.Printf("Update query: %s", query)
	log.Printf("args: %s", args)

	return err
}
