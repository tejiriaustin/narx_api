package services

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/tejiriaustin/narx_api/env"
	"github.com/tejiriaustin/narx_api/models"
	"github.com/tejiriaustin/narx_api/repository"
	"github.com/tejiriaustin/narx_api/utils"
)

type (
	SensorService struct {
		conf *env.Environment
	}

	CreateSensorInput struct {
		Name        string              `json:"name" bson:"name"`
		IpAddress   string              `json:"ipAddress" bson:"ip_address"`
		AccountInfo *models.AccountInfo `json:"accountInfo" bson:"account_info"`
	}

	UpdateSensorInput struct {
		ID        string `json:"id" bson:"id"`
		Name      string `json:"name" bson:"name"`
		IpAddress string `json:"ipAddress" bson:"ip_address"`
	}

	SensorListFilters struct {
		Query     string // for partial free hand lookups
		AccountId string
	}

	ListSensorsInput struct {
		Pager
		Projection *repository.QueryProjection
		Sort       *repository.QuerySort
		Filters    SensorListFilters
	}
)

func NewSensorService(conf *env.Environment) *SensorService {
	return &SensorService{
		conf: conf,
	}
}

var _ SensorServiceInterface = (*SensorService)(nil)

func (s *SensorService) CreateSensor(ctx context.Context,
	input CreateSensorInput,
	passwordGen utils.StrGenFunc,
	sensorRepo *repository.Repository[models.Sensor],
) (*models.Sensor, error) {
	if input.Name == "" {
		return nil, errors.New("sensor name cannot be empty")
	}
	if input.IpAddress == "" {
		return nil, errors.New("please set your sensors ip address")
	}

	now := time.Now().UTC()
	sensor := models.Sensor{
		Shared: models.Shared{
			ID:        primitive.NewObjectID(),
			CreatedAt: &now,
		},
		AccountInfo: models.AccountInfo{},
		Name:        input.Name,
		IpAddress:   input.IpAddress,
		Status:      "good",
		Token:       passwordGen(),
	}
	sensor, err := sensorRepo.Create(ctx, sensor)
	if err != nil {
		return nil, err
	}
	return &sensor, nil
}

func (s *SensorService) UpdateSensor(ctx context.Context,
	input UpdateSensorInput,
	sensorRepo *repository.Repository[models.Sensor],
) (*models.Sensor, error) {
	fields := map[string]interface{}{}

	if input.Name != "" {
		fields["name"] = input.Name
	}
	if input.IpAddress != "" {
		fields["ip_address"] = input.IpAddress
	}

	updates := map[string]interface{}{
		"$set": fields,
	}

	id, err := primitive.ObjectIDFromHex(input.ID)
	if err != nil {
		return nil, errors.New("invalid id")
	}

	filter := repository.NewQueryFilter().AddFilter(models.FieldId, id)
	err = sensorRepo.UpdateMany(ctx, filter, updates)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *SensorService) GetSensor(ctx context.Context,
	sensorId string,
	sensorRepo *repository.Repository[models.Sensor],
) (*models.Sensor, error) {
	id, err := primitive.ObjectIDFromHex(sensorId)
	if err != nil {
		return nil, errors.New("invalid id")
	}
	filter := repository.NewQueryFilter().AddFilter(models.FieldId, id)

	sensor, err := sensorRepo.FindOne(ctx, filter, nil, nil)
	if err != nil {
		return nil, err
	}

	return &sensor, nil
}

func (s *SensorService) ListSensors(ctx context.Context,
	input ListSensorsInput,
	sensorRepo *repository.Repository[models.Sensor],
) ([]models.Sensor, *repository.Paginator, error) {
	filter := repository.NewQueryFilter()

	if input.Filters.AccountId != "" {
		filter.AddFilter("account_info._id", input.Filters.AccountId)
	}

	if input.Filters.Query != "" {
		freeHandFilters := []map[string]interface{}{
			{"name": map[string]interface{}{"$regex": input.Filters.Query, "$options": "i"}},
			{"ip_address": map[string]interface{}{"$regex": input.Filters.Query, "$options": "i"}},
			{"status": map[string]interface{}{"$regex": input.Filters.Query, "$options": "i"}},
			{"account_info.first_name": map[string]interface{}{"$regex": input.Filters.Query, "$options": "i"}},
			{"account_info.last_name": map[string]interface{}{"$regex": input.Filters.Query, "$options": "i"}},
			{"account_info.email": map[string]interface{}{"$regex": input.Filters.Query, "$options": "i"}},
		}
		filter.AddFilter("$or", freeHandFilters)
	}

	account, paginator, err := sensorRepo.Paginate(ctx, filter, input.PerPage, input.Page, input.Projection, input.Sort)
	if err != nil {
		return nil, nil, err
	}

	return account, paginator, nil
}

func (s *SensorService) DeleteSensor(ctx context.Context,
	sensorId string,
	sensorRepo *repository.Repository[models.Sensor],
) error {
	id, err := primitive.ObjectIDFromHex(sensorId)
	if err != nil {
		return errors.New("invalid id")
	}
	filter := repository.NewQueryFilter().AddFilter(models.FieldId, id)

	return sensorRepo.DeleteMany(ctx, filter)
}
