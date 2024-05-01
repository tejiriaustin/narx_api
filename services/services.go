package services

import (
	"context"
	"github.com/tejiriaustin/narx_api/constants"
	"github.com/tejiriaustin/narx_api/env"
	"github.com/tejiriaustin/narx_api/publisher"
	"log"
)

type (
	Container struct {
		AccountsService AccountsServiceInterface
		Publisher       publisher.PublishInterface
	}

	Pager struct {
		Page    int64
		PerPage int64
	}
)

func NewService(conf *env.Environment) *Container {
	log.Println("Creating Container...")
	return &Container{
		AccountsService: NewAccountsService(conf),
	}
}

func GetPageNumberFromContext(ctx context.Context) int64 {
	n, ok := ctx.Value(constants.ContextKeyPageNumber).(int64)
	if !ok {
		return 0
	}
	return n
}

func GetPerPageLimitFromContext(ctx context.Context) int64 {
	l, ok := ctx.Value(constants.ContextKeyPerPageLimit).(int64)
	if !ok {
		return 0
	}
	return l
}
