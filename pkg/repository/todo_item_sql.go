package repository

import (
	"WebApi"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type TodoItemSql struct {
	db *sqlx.DB
}

func NewTodoItemSql(db *sqlx.DB) *TodoItemSql {
	return &TodoItemSql{db: db}
}

func (s *TodoItemSql) Create(listId int, item WebApi.TodoItem) (int, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return 0, err
	}

	createItemQuery := fmt.Sprintf("INSERT INTO %s (title, description) values ($1, $2) RETURNING id",
		todoItemsTable)

	row := s.db.QueryRow(createItemQuery, item.Title, item.Description)

	var itemId int
	err = row.Scan(&itemId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	createListItemsQuery := fmt.Sprintf("INSERT INTO %s (list_id, item_id) values ($1, $2)", listsItemsTable)
	_, err = s.db.Exec(createListItemsQuery, listId, itemId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return itemId, tx.Commit()
}

func (s *TodoItemSql) GetAll(userId, listId int) ([]WebApi.TodoItem, error) {
	var items []WebApi.TodoItem
	query := fmt.Sprintf(
		`SELECT ti.id, ti.title, ti.description, ti.done 
		FROM %s ti 
		INNER JOIN %s li ON ti.id = li.item_id 
		INNER JOIN %s ul ON ul.list_id = li.list_id
		WHERE ul.user_id = $1 AND li.list_id = $2`,
		todoItemsTable, listsItemsTable, usersListsTable)

	if err := s.db.Select(&items, query, userId, listId); err != nil {
		return nil, err
	}

	return items, nil
}

func (s *TodoItemSql) GetById(itemId, userId int) (WebApi.TodoItem, error) {
	var item WebApi.TodoItem

	query := fmt.Sprintf(
		`SELECT ti.id, ti.title, ti.description, ti.done
		FROM %s ti
		INNER JOIN %s li ON li.item_id = ti.id
		INNER JOIN %s ul ON ul.list_id = li.list_id
		WHERE ti.id = $1 AND ul.user_id = $2`,
		todoItemsTable, listsItemsTable, usersListsTable,
	)

	err := s.db.Get(&item, query, itemId, userId)

	return item, err
}

func (s *TodoItemSql) Delete(itemId, userId int) error {
	query := fmt.Sprintf(
		`DELETE FROM %s ti USING %s li, %s ul
		WHERE ti.id = $1 AND ti.id = li.item_id AND li.list_id = ul.list_id AND ul.user_id = $2`,
		todoItemsTable, listsItemsTable, usersListsTable,
	)

	_, err := s.db.Exec(query, itemId, userId)

	return err
}

func (s *TodoItemSql) Update(itemId, userId int, input WebApi.UpdateTodoItemInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *input.Title)
		argId++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	if input.Done != nil {
		setValues = append(setValues, fmt.Sprintf("done=$%d", argId))
		args = append(args, *input.Done)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(
		`UPDATE %s ti SET %s 
		FROM %s li, %s ul
		WHERE ti.id = li.item_id AND li.list_id = ul.list_id AND ul.user_id = $%d AND ti.id = $%d`,
		todoItemsTable, setQuery, listsItemsTable, usersListsTable, argId, argId+1,
	)

	args = append(args, userId, itemId)

	_, err := s.db.Exec(query, args...)

	return err
}
