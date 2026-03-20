package taskService

import (
	"image"
	"test/internal/domain"
	"test/internal/repository"
	"test/pkg"

	"github.com/google/uuid"
)

type Object struct {
	service repository.Object
	sender  repository.TaskManager
}

func NewObject(service repository.Object, sender repository.TaskManager) *Object {
	return &Object{service: service, sender: sender}
}

func (s *Object) NewTask(name string, parameters any) (string, error) {

	task := domain.Task{
		Id:              uuid.New().String(),
		ImageBase:       "images/test2.png",
		Status:          "running",
		FilterName:      name,
		FilterParametes: parameters,
		Result:          "none",
	}
	if err := s.sender.Send(task); err != nil {
		return "", err
	}

	if err := s.service.PostTask(task.Id, task); err != nil {
		return "", err
	}
	return task.Id, nil
}

func (s *Object) GetTaskResult(id string) (image.Image, error) {
	task, err := s.service.GetTask(id)
	if err != nil {
		return nil, err
	}
	if task.Status == "running" {
		return nil, nil
	}

	img, err := pkg.FromFileToImage(task.Result)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func (s *Object) GetTaskStatus(id string) (string, error) {
	task, err := s.service.GetTask(id)
	if err != nil {
		return "", err
	}
	return task.Status, nil
}

func (s *Object) PutTask(task domain.Task) error {
	return s.service.PutTask(task.Id, task)
}
