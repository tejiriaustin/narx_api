package repository

import "go.mongodb.org/mongo-driver/bson"

type QueryFilter struct {
	filters bson.D
}

func newQF() *QueryFilter {
	return &QueryFilter{filters: bson.D{}}
}

func NewQueryFilter() *QueryFilter {
	return newQF()
}

func (qf *QueryFilter) AddFilter(name string, value interface{}) *QueryFilter {
	qf.filters = append(qf.filters, bson.E{Key: name, Value: value})
	return qf
}

func (qf *QueryFilter) GetFilters() bson.D {
	return qf.filters
}
