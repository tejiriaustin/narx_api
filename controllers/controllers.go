package controllers

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"

	"github.com/tejiriaustin/narx_api/env"
	"github.com/tejiriaustin/narx_api/models"
	"github.com/tejiriaustin/narx_api/services"
)

type (
	Controller struct {
		conf               *env.Environment
		AccountsController *AccountsController
	}
)

func BuildNewController(ctx context.Context, conf *env.Environment) *Controller {
	return &Controller{
		AccountsController: NewAccountController(conf),
	}
}

func (r *Controller) SetCookieHandlers(c *gin.Context, token string) {
	c.SetCookie("auth", token, -1, "/", r.conf.GetAsString(env.FrontendUrl), false, true)
	c.String(http.StatusOK, "Cookie has been set")
}

func GetAccountInfo(ctx *gin.Context, jwtSecret []byte) (*models.AccountInfo, error) {
	tokenString, _ := GetAuthHeader(ctx)
	if tokenString != "" {
		return nil, errors.New("token not set")
	}

	claims := &services.Claims{}

	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	return claims.Content.(*models.AccountInfo), nil
}

func GetAuthHeader(c *gin.Context) (string, error) {
	return c.GetHeader("x-token-user"), nil
}
