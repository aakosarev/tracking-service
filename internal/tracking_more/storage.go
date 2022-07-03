package tracking_more

import (
	"context"
	"fmt"
	"github.com/aakosarev/tracking-service/pkg/logging"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
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
	trackinfo := []types.AttributeValue{}
	for _, v := range databaseData.Trackinfo {
		m := &types.AttributeValueMemberM{
			Value: map[string]types.AttributeValue{
				"checkpoint_date":               &types.AttributeValueMemberS{Value: v.CheckpointDate},
				"checkpoint_delivery_status":    &types.AttributeValueMemberS{Value: v.CheckpointDeliveryStatus},
				"checkpoint_delivery_substatus": &types.AttributeValueMemberS{Value: v.CheckpointDeliverySubstatus},
				"location":                      &types.AttributeValueMemberS{Value: v.Location},
				"tracking_detail":               &types.AttributeValueMemberS{Value: v.TrackingDetail},
			},
		}
		trackinfo = append(trackinfo, m)
	}
	_, err := s.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("tracking_more"),
		Item: map[string]types.AttributeValue{
			"tracking_number":         &types.AttributeValueMemberS{Value: databaseData.TrackingNumber},
			"courier_code":            &types.AttributeValueMemberS{Value: databaseData.CourierCode},
			"lastest_checkpoint_time": &types.AttributeValueMemberS{Value: databaseData.LastestCheckpointTime},
			"latest_event":            &types.AttributeValueMemberS{Value: databaseData.LatestEvent},
			"trackinfo":               &types.AttributeValueMemberL{Value: trackinfo},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to insert into the database, %w", err)
	}
	return nil
}

func (s *storage) GetAll() ([]InputData, error) {
	result, err := s.client.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: aws.String("tracking_more"),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get all items from database, %w", err)
	}
	inputData := []InputData{}
	for _, v := range result.Items {
		var trackNumber, courierCode string
		if err = attributevalue.Unmarshal(v["tracking_number"], &trackNumber); err != nil {
			return nil, fmt.Errorf("failed to unmarshal, %w", err)
		}
		if err = attributevalue.Unmarshal(v["courier_code"], &courierCode); err != nil {
			return nil, fmt.Errorf("failed to unmarshal, %w", err)
		}
		input := InputData{
			TrackingNumber: trackNumber,
			CourierCode:    courierCode,
		}
		inputData = append(inputData, input)
	}
	return inputData, nil
}
