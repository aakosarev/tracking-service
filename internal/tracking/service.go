package tracking

import (
	"fmt"
	"github.com/aakosarev/tracking-service/pkg/logging"
)

type Service struct {
	storage Storage
	logger  *logging.Logger
}

func NewService(storage Storage, logger *logging.Logger) Service {
	return Service{
		storage: storage,
		logger:  logger,
	}
}

func (s *Service) Insert(databaseData DatabaseData) error {
	err := s.storage.Insert(databaseData)
	if err != nil {
		return fmt.Errorf("failed to insert into the database, %v", err)
	}
	return nil
}
