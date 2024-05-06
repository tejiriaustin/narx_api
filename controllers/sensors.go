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
	"github.com/tejiriaustin/narx_api/utils"
)

type SensorController struct {
	conf *env.Environment
}

func NewSensorController(conf *env.Environment) *SensorController {
	return &SensorController{
		conf: conf,
	}
}

func (s *SensorController) AddSensor(
	passwordGen utils.StrGenFunc,
	sensorService services.SensorServiceInterface,
	sensorRepo *repository.Repository[models.Sensor],
) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var req requests.CreateSensorRequest

		err := ctx.BindJSON(&req)
		if err != nil {
			response.FormatResponse(ctx, http.StatusBadRequest, "Bad Request", nil)
			return
		}

		accountInfo, err := GetAccountInfo(ctx, s.conf.GetAsBytes(env.JwtSecret))
		if err != nil {
			response.FormatResponse(ctx, http.StatusUnauthorized, err.Error(), nil)
			return
		}

		input := services.CreateSensorInput{
			Name:        req.Name,
			IpAddress:   req.IpAddress,
			AccountInfo: accountInfo,
		}

		sensor, err := sensorService.CreateSensor(ctx, input, passwordGen, sensorRepo)
		if err != nil {
			response.FormatResponse(ctx, http.StatusBadRequest, err.Error(), nil)
			return
		}

		response.FormatResponse(ctx, http.StatusOK, "successful", response.SingleSensorResponse(sensor))
	}
}

func (s *SensorController) GetSensor(
	sensorService services.SensorServiceInterface,
	sensorRepo *repository.Repository[models.Sensor],
) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		sensorId := ctx.Param("sensor_id")

		_, err := GetAccountInfo(ctx, s.conf.GetAsBytes(env.JwtSecret))
		if err != nil {
			response.FormatResponse(ctx, http.StatusUnauthorized, "unauthorised access", nil)
			return
		}

		sensor, err := sensorService.GetSensor(ctx, sensorId, sensorRepo)
		if err != nil {
			response.FormatResponse(ctx, http.StatusBadRequest, err.Error(), nil)
			return
		}

		response.FormatResponse(ctx, http.StatusOK, "successful", response.SingleSensorResponse(sensor))
	}
}

func (s *SensorController) ListSensor(
	sensorService services.SensorServiceInterface,
	sensorRepo *repository.Repository[models.Sensor],
) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		accountInfo, err := GetAccountInfo(ctx, s.conf.GetAsBytes(env.JwtSecret))
		if err != nil {
			response.FormatResponse(ctx, http.StatusUnauthorized, "unauthorised access", nil)
			return
		}

		query := ctx.Param("query")

		input := services.ListSensorsInput{
			Pager: services.Pager{
				Page:    services.GetPageNumberFromContext(ctx),
				PerPage: services.GetPerPageLimitFromContext(ctx),
			},
			Filters: services.SensorListFilters{Query: query, AccountId: accountInfo.Id},
		}

		sensors, paginator, err := sensorService.ListSensors(ctx, input, sensorRepo)
		if err != nil {
			response.FormatResponse(ctx, http.StatusBadRequest, err.Error(), nil)
			return
		}

		payload := map[string]interface{}{
			"records": response.MultipleSensorResponse(sensors),
			"meta":    paginator,
		}

		response.FormatResponse(ctx, http.StatusOK, "successful", payload)
	}
}

func (s *SensorController) UpdateSensor(
	sensorService services.SensorServiceInterface,
	sensorRepo *repository.Repository[models.Sensor],
) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var req requests.UpdateSensorRequest

		_, err := GetAccountInfo(ctx, s.conf.GetAsBytes(env.JwtSecret))
		if err != nil {
			response.FormatResponse(ctx, http.StatusUnauthorized, "unauthorised access", nil)
			return
		}

		err = ctx.BindJSON(&req)
		if err != nil {
			response.FormatResponse(ctx, http.StatusBadRequest, "Bad Request", nil)
			return
		}

		input := services.UpdateSensorInput{
			ID:        req.ID,
			Name:      req.Name,
			IpAddress: req.IpAddress,
		}

		_, err = sensorService.UpdateSensor(ctx, input, sensorRepo)
		if err != nil {
			response.FormatResponse(ctx, http.StatusBadRequest, err.Error(), nil)
			return
		}

		response.FormatResponse(ctx, http.StatusOK, "successful", nil)
	}
}

func (s *SensorController) DeleteSensor(
	sensorService services.SensorServiceInterface,
	sensorRepo *repository.Repository[models.Sensor],
) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		sensorId := ctx.Param("sensor_id")

		_, err := GetAccountInfo(ctx, s.conf.GetAsBytes(env.JwtSecret))
		if err != nil {
			response.FormatResponse(ctx, http.StatusUnauthorized, "unauthorised access", nil)
			return
		}

		err = sensorService.DeleteSensor(ctx, sensorId, sensorRepo)
		if err != nil {
			response.FormatResponse(ctx, http.StatusBadRequest, err.Error(), nil)
			return
		}

		response.FormatResponse(ctx, http.StatusOK, "successful", nil)
	}
}
