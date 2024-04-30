package server

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/tejiriaustin/narx_api/controllers"
	"github.com/tejiriaustin/narx_api/env"
	"github.com/tejiriaustin/narx_api/middleware"
	"github.com/tejiriaustin/narx_api/repository"
	"github.com/tejiriaustin/narx_api/services"
)

func Start(
	ctx context.Context,
	service *services.Container,
	repo *repository.Container,
	conf *env.Environment,
) {
	router := gin.New()

	router.Use(
		middleware.DefaultStructuredLogs(),
		middleware.ReadPaginationOptions(),
		middleware.CORSMiddleware(),
	)
	log.Println("starting server...")

	controllers.BindRoutes(ctx, router, service, repo, conf)

	go func() {
		if err := router.Run(); err != nil {
			log.Fatal("shutting down server...:", err.Error())
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := router; err != nil {
		log.Fatal(err)
	}
}
