package service

import (
	"speed-up/service/application/exception"
	"speed-up/service/application/ports"
	repositories "speed-up/service/application/repository"
)

type MemoryService struct {
	Repository repositories.DataRepository
}

func NewMemoryService(repository repositories.DataRepository) ports.MemoryData {
	return MemoryService{
		Repository: repository,
	}
}

func (m MemoryService) Set(key, value string) error {

	if key == "" {
		return exception.ThrowValidationError("Key is null")
	}

	if value == "" {
		return exception.ThrowValidationError("Value is null")
	}

	m.Repository.Set(key, value)

	return nil
}

func (m MemoryService) Get(key string) (string, error) {

	if key == "" {
		return "", exception.ThrowValidationError("Key is null")
	}

	return m.Repository.Get(key), nil
}
