package repository

import (
	"context"
	"github.com/tejiriaustin/narx_api/models"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	Creator[T models.SharedInterface] interface {
		Create(ctx context.Context, data T) (T, error)
	}
	Finder[T models.SharedInterface] interface {
		FindOne(ctx context.Context, queryFilter *QueryFilter, projection *QueryProjection, findOneOptions ...*options.FindOneOptions) (T, error)
	}
	Deleter[T models.SharedInterface] interface {
		DeleteMany(ctx context.Context, queryFilter *QueryFilter) error
	}
	Updater[T models.SharedInterface] interface {
		Update(ctx context.Context, dataObject T) (T, error)
	}
)
