package tracking_more

import (
	"fmt"
	"github.com/aakosarev/tracking-service/pkg/logging"
)

type Storage interface {
	CreateOrUpdate(databaseData DatabaseData) error
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
		return fmt.Errorf("failed to insert into the database, %v", err)
	}
	return nil
}
