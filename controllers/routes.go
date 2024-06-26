package controllers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/tejiriaustin/narx_api/env"
	"github.com/tejiriaustin/narx_api/repository"
	"github.com/tejiriaustin/narx_api/response"
	"github.com/tejiriaustin/narx_api/services"
	"github.com/tejiriaustin/narx_api/utils"
)

func BindRoutes(
	ctx context.Context,
	routerEngine *gin.Engine,
	sc *services.Container,
	repos *repository.Container,
	conf *env.Environment,
) {

	controllers := BuildNewController(ctx, conf)

	passwordGenerator := utils.RandomStringGenerator()

	r := routerEngine.Group("/v1")

	r.GET("/health", func(c *gin.Context) {
		response.FormatResponse(c, http.StatusOK, "OK", nil)
	})

	accounts := r.Group("/user")
	{
		accounts.POST("/sign-up", controllers.AccountsController.SignUp(passwordGenerator, sc.AccountsService, repos.AccountsRepo))
		accounts.POST("/login", controllers.AccountsController.Login(sc.AccountsService, repos.AccountsRepo))
		accounts.POST("/forgot-password", controllers.AccountsController.ForgotPassword(sc.AccountsService, repos.AccountsRepo, sc.Publisher))
		accounts.POST("/reset-password", controllers.AccountsController.ResetPassword(sc.AccountsService, repos.AccountsRepo))
		accounts.GET("/")
		accounts.PUT("/edit-account", controllers.AccountsController.EditAccount(sc.AccountsService, repos.AccountsRepo))
	}

	sensors := r.Group("/sensors")
	{
		sensors.POST("/add", controllers.SensorController.AddSensor(passwordGenerator, sc.SensorService, repos.SensorRepo))
		sensors.PUT("/update", controllers.SensorController.UpdateSensor(sc.SensorService, repos.SensorRepo))
		sensors.GET("/:sensor_id", controllers.SensorController.GetSensor(sc.SensorService, repos.SensorRepo))
		sensors.GET("/list", controllers.SensorController.ListSensor(sc.SensorService, repos.SensorRepo))
		sensors.DELETE("/:sensor_id", controllers.SensorController.DeleteSensor(sc.SensorService, repos.SensorRepo))
	}

	devices := r.Group("/devices")
	{
		devices.POST("", controllers.DeviceController.SaveDeviceToken(sc.DeviceService, repos.DevicesRepo))
	}
}
