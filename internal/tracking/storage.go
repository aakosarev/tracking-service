package tracking

import (
	"context"
	"fmt"
	"github.com/aakosarev/tracking-service/pkg/logging"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Storage struct {
	client *dynamodb.Client
	logger *logging.Logger
}

func NewStorage(client *dynamodb.Client, logger *logging.Logger) Storage {
	return Storage{
		client: client,
		logger: logger,
	}
}

func (s *Storage) Insert(databaseData DatabaseData) error {
	item, err := attributevalue.MarshalMap(databaseData)
	if err != nil {
		return fmt.Errorf("failed to marshal, %v", err)
	}
	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String("trackings"),
	}
	if _, err := s.client.PutItem(context.TODO(), input); err != nil {
		return fmt.Errorf("failed to insert into the database, %v", err)
	}
	return nil
}
