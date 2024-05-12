package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/tejiriaustin/narx_api/env"
	"github.com/tejiriaustin/narx_api/models"
	"github.com/tejiriaustin/narx_api/repository"
	"github.com/tejiriaustin/narx_api/requests"
	"github.com/tejiriaustin/narx_api/response"
	"github.com/tejiriaustin/narx_api/services"
)

type DeviceController struct {
	conf *env.Environment
}

func NewDeviceController(conf *env.Environment) *DeviceController {
	return &DeviceController{
		conf: conf,
	}
}

func (c *DeviceController) SaveDeviceToken(
	deviceService services.DeviceServiceInterface,
	deviceRepo *repository.Repository[models.Devices],
) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var req requests.SaveDeviceToken

		err := ctx.BindJSON(&req)
		if err != nil {
			response.FormatResponse(ctx, http.StatusBadRequest, "Bad Request", nil)
			return
		}

		accountInfo, err := GetAccountInfo(ctx, c.conf.GetAsBytes(env.JwtSecret))
		if err != nil {
			response.FormatResponse(ctx, http.StatusUnauthorized, "Unauthorized access", nil)
			return
		}

		input := services.SaveDeviceTokenInput{
			AccountInfo: *accountInfo,
			DeviceToken: req.DeviceToken,
		}

		err = deviceService.SaveDeviceToken(ctx, input, deviceRepo)
		if err != nil {
			response.FormatResponse(ctx, http.StatusBadRequest, err.Error(), nil)
			return
		}

		response.FormatResponse(ctx, http.StatusOK, "successful", nil)
	}
}
