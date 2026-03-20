package usecases

import (
	"image"
	"test/internal/domain"
)

type TaskService interface {
	NewTask(name string, parameters any) (string, error)
	GetTaskResult(id string) (image.Image, error)
	GetTaskStatus(id string) (string, error)
	PutTask(domain.Task) error
}
