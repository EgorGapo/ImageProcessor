package ramStorage

import (
	"errors"
	"test/internal/domain"
)

type RamStorage struct {
	tasks    map[string]domain.Task
	sessions map[string]string
	users    map[string]domain.User
}

func NewRamStorage() *RamStorage {
	return &RamStorage{
		tasks:    make(map[string]domain.Task),
		users:    make(map[string]domain.User),
		sessions: make(map[string]string),
	}
}

func (rs *RamStorage) GetTask(key string) (*domain.Task, error) {
	value, exists := rs.tasks[key]
	if !exists {
		return nil, errors.New("task not found")
	}
	return &value, nil
}

func (rs *RamStorage) PutTask(key string, value domain.Task) error {
	rs.tasks[key] = value
	return nil
}

func (rs *RamStorage) PostTask(key string, value domain.Task) error {
	if _, exists := rs.tasks[key]; exists {
		return errors.New("task already exists")
	}
	rs.tasks[key] = value
	return nil
}

func (rs *RamStorage) GetUser(login string) (*domain.User, error) {
	user, exists := rs.users[login]
	if !exists {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func (rs *RamStorage) CreateUser(user domain.User) error {
	if _, exists := rs.users[user.Login]; exists {
		return errors.New("user already exists")
	}
	rs.users[user.Login] = user
	return nil
}

func (rs *RamStorage) CreateSession(sessionID, userID string) error {
	if _, exists := rs.sessions[sessionID]; exists {
		return errors.New("session already exists")
	}
	rs.sessions[sessionID] = userID
	return nil
}

func (rs *RamStorage) GetUserBySession(sessionID string) (string, error) {
	userID, exists := rs.sessions[sessionID]
	if !exists {
		return "", errors.New("session not found")
	}
	return userID, nil
}
