package parcels_app

import (
	"context"
	"fmt"
	"github.com/aakosarev/tracking-service/pkg/logging"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type storage struct {
	client *dynamodb.Client
	logger *logging.Logger
}

func NewStorage(client *dynamodb.Client, logger *logging.Logger) *storage {
	return &storage{
		client: client,
		logger: logger,
	}
}

func (s *storage) CreateOrUpdate(databaseData DatabaseData) error {
	states := []types.AttributeValue{}

	for _, v := range databaseData.States {
		m := &types.AttributeValueMemberM{
			Value: map[string]types.AttributeValue{
				"location": &types.AttributeValueMemberS{Value: v.Location},
				"date":     &types.AttributeValueMemberS{Value: v.Date},
				"carrier":  &types.AttributeValueMemberS{Value: string(v.Carrier)},
				"status":   &types.AttributeValueMemberS{Value: v.Status},
			},
		}
		states = append(states, m)
	}

	carriers := []types.AttributeValue{}
	for _, v := range databaseData.Carriers {
		c := &types.AttributeValueMemberS{
			Value: v,
		}
		carriers = append(carriers, c)
	}

	_, err := s.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("parcels_app"),
		Item: map[string]types.AttributeValue{
			"tracking_number": &types.AttributeValueMemberS{Value: databaseData.TrackingNumber},
			"status":          &types.AttributeValueMemberS{Value: databaseData.Status},
			"sub_status":      &types.AttributeValueMemberS{Value: databaseData.SubStatus},
			"origin":          &types.AttributeValueMemberS{Value: databaseData.Origin},
			"destination":     &types.AttributeValueMemberS{Value: databaseData.Destination},
			"from":            &types.AttributeValueMemberS{Value: databaseData.From},
			"to":              &types.AttributeValueMemberS{Value: databaseData.To},
			"weight":          &types.AttributeValueMemberS{Value: databaseData.Weight},
			"carriers":        &types.AttributeValueMemberL{Value: carriers},
			"states":          &types.AttributeValueMemberL{Value: states},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to insert into the database, %w", err)
	}
	return nil
}
