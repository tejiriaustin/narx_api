package services

import (
	"context"
	"errors"
	"github.com/tejiriaustin/narx_api/env"
	"github.com/tejiriaustin/narx_api/models"
	"github.com/tejiriaustin/narx_api/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"
)

type DeviceService struct {
	conf *env.Environment
}

func NewDeviceService(conf *env.Environment) *DeviceService {
	return &DeviceService{
		conf: conf,
	}
}

type (
	SaveDeviceTokenInput struct {
		DeviceToken string             `json:"deviceToken" bson:"device_token"`
		AccountInfo models.AccountInfo `json:"accountInfo" bson:"account_info"`
	}
)

func (s *DeviceService) SaveDeviceToken(ctx context.Context,
	input SaveDeviceTokenInput,
	devicesRepo *repository.Repository[models.Devices],
) error {
	if input.DeviceToken == "" {
		log.Println("device token is required")
		return errors.New("device token is required")
	}

	now := time.Now()

	device := models.Devices{
		Shared: models.Shared{
			ID:        primitive.NewObjectID(),
			CreatedAt: &now,
		},
		AccountInfo: input.AccountInfo,
		DeviceToken: input.DeviceToken,
	}
	_, err := devicesRepo.Create(ctx, device)
	if err != nil {
		return err
	}
	return nil
}
