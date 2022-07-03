package tracking_more

import (
	"github.com/aakosarev/tracking-service/pkg/logging"
)

type Storage interface {
	CreateOrUpdate(databaseData DatabaseData) error
	GetAll() ([]InputData, error)
}

type service struct {
	storage Storage
	logger  *logging.Logger
}

func NewService(storage Storage, logger *logging.Logger) *service {
	return &service{
		storage: storage,
		logger:  logger,
	}
}

func (s *service) CreateOrUpdate(databaseData DatabaseData) error {
	err := s.storage.CreateOrUpdate(databaseData)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) UpdateAll() ([]InputData, error) {
	inputData, err := s.storage.GetAll()
	if err != nil {
		return nil, err
	}
	return inputData, nil
}
