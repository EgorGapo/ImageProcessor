package repository

import (
	"test/internal/domain"
)

type Object interface {
	GetTask(key string) (*domain.Task, error)
	PutTask(key string, value domain.Task) error
	PostTask(key string, value domain.Task) error
	GetUser(login string) (*domain.User, error)
	CreateUser(user domain.User) error
	CreateSession(sessionID, userID string) error
	GetUserBySession(sessionID string) (string, error)
}
