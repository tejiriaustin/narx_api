package repository

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson"
)

type QuerySort struct {
	sort bson.D
}

func NewDefaultQuerySort() *QuerySort {
	return &QuerySort{
		sort: bson.D{
			bson.E{
				Key: "_id", Value: -1},
		},
	}
}

func NewQuerySort() *QuerySort {
	return &QuerySort{
		sort: bson.D{},
	}
}

func (qs *QuerySort) AddSort(field string, value interface{}) (*QuerySort, error) {
	if field == "" {
		return qs, errors.New("field Name is required")
	}

	for _, elem := range qs.sort {
		if elem.Key == field {
			return qs, errors.New("field already added")
		}
	}

	qs.sort = append(qs.sort, bson.E{Key: field, Value: value})
	return qs, nil
}

func (qs *QuerySort) GetSort() bson.D {
	return qs.sort
}
