package repository

import "test/internal/domain"

type TaskManager interface {
	Send(domain.Task) error
	Close()
}
