package repository

import (
	"WebApi"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type AuthSql struct {
	db *sqlx.DB
}

func NewAuthSql(db *sqlx.DB) *AuthSql {
	return &AuthSql{db: db}
}

func (s *AuthSql) CreateUser(user WebApi.User) (int, error) {
	query := fmt.Sprintf("INSERT INTO %s (name, username, password_hash) values ($1, $2, $3) RETURNING id", usersTable)
	var id int

	row := s.db.QueryRow(query, user.Name, user.Username, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (s *AuthSql) GetUser(username, password string) (WebApi.User, error) {
	var user WebApi.User

	query := fmt.Sprintf("SELECT id FROM %s WHERE username=$1 AND password_hash=$2", usersTable)

	err := s.db.Get(&user, query, username, password)
	if err != nil {
		fmt.Printf("Error: %s\n Username: %s\n Password: %s\n\n", err.Error(), username, password)
	}
	return user, err
}
