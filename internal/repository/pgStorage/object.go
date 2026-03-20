package pgstorage

import (
	"database/sql"
	"errors"
	"test/internal/domain"

	_ "github.com/lib/pq"
)

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorsge(connStr string) (*PostgresStorage, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStorage{db: db}, nil
}

func (rs *PostgresStorage) GetTask(key string) (*domain.Task, error) {
	var value domain.Task
	err := rs.db.QueryRow("SELECT id, image_base, status, filter_name, filter_parametes, result FROM tasks WHERE id = $1", key).
		Scan(&value.Id, &value.ImageBase, &value.Status, &value.FilterName, &value.FilterParametes, &value.Result)
	if err == sql.ErrNoRows {
		return nil, errors.New("task not found")
	} else if err != nil {
		return nil, err
	}

	return &value, nil
}
func (rs *PostgresStorage) PutTask(key string, value domain.Task) error {
	query := `
        INSERT INTO tasks (id, image_base, filter_name, filter_parametes, status, result)
        VALUES ($1, $2, $3, $4, $5, $6)
        ON CONFLICT (id) DO UPDATE
        SET image_base = $2, filter_name = $3, filter_parametes = $4, status = $5, result = $6
    `

	_, err := rs.db.Exec(query, key, value.ImageBase, value.FilterName, value.FilterParametes, value.Status, value.Result)
	if err != nil {
		return err
	}

	return nil
}

func (rs *PostgresStorage) PostTask(key string, value domain.Task) error {
	query := `
        INSERT INTO tasks (id, image_base, filter_name, filter_parametes, status, result)
        VALUES ($1, $2, $3, $4, $5, $6)
    `
	_, err := rs.db.Exec(query, key, value.ImageBase, value.FilterName, value.FilterParametes, value.Status, value.Result)
	if err != nil {
		return err
	}

	return nil
}

func (rs *PostgresStorage) GetUser(login string) (*domain.User, error) {
	var value domain.User
	err := rs.db.QueryRow("Select id, login, password From users Where login = $1", login).Scan(&value.Id, &value.Login, &value.Password)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	} else if err != nil {
		return nil, err
	}

	return &value, nil
}
func (rs *PostgresStorage) CreateUser(user domain.User) error {
	query := `
        INSERT INTO users (id, login, password)
        VALUES ($1, $2, $3)
        ON CONFLICT (login) DO NOTHING
    `

	result, err := rs.db.Exec(query, user.Id, user.Login, user.Password)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("user with this login already exists")
	}

	return nil
}

func (rs *PostgresStorage) CreateSession(sessionID, userID string) error {
	query := `
        INSERT INTO sessions (user_id, session_id)
        VALUES ($1, $2)
    `
	_, err := rs.db.Exec(query, userID, sessionID)
	if err != nil {
		return err
	}
	return nil
}

func (rs *PostgresStorage) GetUserBySession(sessionID string) (string, error) {
	var userID string
	err := rs.db.QueryRow("SELECT user_id FROM sessions WHERE session_id = $1", sessionID).Scan(&userID)
	if err == sql.ErrNoRows {
		return "", errors.New("session not found")
	} else if err != nil {
		return "", err
	}

	return userID, nil
}

func (s *PostgresStorage) IsFetched(url string) (bool, error) {
	query := `
        SELECT 1
        FROM crawled_pages
        WHERE url = $1
        LIMIT 1
    `
	var exists int
	err := s.db.QueryRow(query, url).Scan(&exists)

	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
