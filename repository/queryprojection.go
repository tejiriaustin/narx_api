package repository

import "go.mongodb.org/mongo-driver/bson"

type QueryProjection struct {
	projection bson.M
}

// NewQueryProjection returns a new QueryProjection struct.
// Each call to NewQueryProjection() will empty the projections and restart.
func NewQueryProjection() QueryProjection {
	return QueryProjection{
		projection: bson.M{},
	}
}

// AddProjection will add a projection choice to a field. Opt is an integer between 0 (exclude) or 1(include).
// If an opt of > 1 is provided, AddProjection ignores and proceeds.
// AddProjection is not go-routine safe and should not be called in multiple go routines.
// Multiple calls with the same key would lead to updating the opt each time.
func (qp QueryProjection) AddProjection(field string, opt int) QueryProjection {
	if opt < 0 || opt > 1 {
		return qp
	}

	qp.projection[field] = opt
	return qp
}

// GetProjection returns the collected projections
func (qp QueryProjection) GetProjection() bson.M {
	return qp.projection
}
