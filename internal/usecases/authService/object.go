package authservice

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"test/internal/domain"
	"test/internal/repository"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"
)

func hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("failed hashing")
	}
	return string(bytes), nil
}

func checkPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func generateToken(lenght int) (string, error) {
	bytes := make([]byte, lenght)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("error while decoding token: %v", err)
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

type redisStorage interface {
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
}

type Object struct {
	userService    repository.Object
	sessionService redisStorage
}

func NewObject(service repository.Object, redis redisStorage) *Object {
	return &Object{
		userService:    service,
		sessionService: redis}
}

func (s *Object) Register(login string, password string) error {
	if len(login) < 4 || len(password) < 4 {
		return errors.New("too short login or password")
	}
	hashedPassword, err := hash(password)
	if err != nil {
		return err
	}

	user := domain.User{
		Id:       uuid.New().String(),
		Login:    login,
		Password: hashedPassword,
	}

	if err := s.userService.CreateUser(user); err != nil {
		return err
	}
	return nil
}

func (s *Object) Login(login string, password string) (string, error) {
	user, err := s.userService.GetUser(login)
	if err != nil || !checkPassword(password, user.Password) {
		return "", errors.New("invalid password or login")
	}

	sessionToken, err := generateToken(32)
	if err != nil {
		return "", err
	}
	// Сохраняем токен в хранилище
	//s.userService.CreateSession(sessionToken, user.Id)
	s.sessionService.Set(context.Background(), sessionToken, user.Id, 20*time.Minute)

	return sessionToken, nil
}

func (s *Object) Auth(token string) (string, error) {
	//userId, err := s.userService.GetUserBySession(token)
	userId, err := s.sessionService.Get(context.Background(), token)
	if err != nil {
		return "", err
	}
	return userId, err
}
