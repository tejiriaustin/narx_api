package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type SharedInterface interface {
	Initialize(id primitive.ObjectID, now time.Time)
	GetId() string
	SetID(id primitive.ObjectID)
	SetUsedProjection(flag bool)
	DidUseProjection() bool
	SetUpdatedAt()
}

type Shared struct {
	ID             primitive.ObjectID `json:"id" bson:"_id"`
	CreatedAt      *time.Time         `json:"created_at" bson:"created_at"`
	UpdatedAt      *time.Time         `json:"updated_at" bson:"updated_at"`
	DeletedAt      *time.Time         `json:"deleted_at" bson:"deleted_at"`
	usedProjection bool               `bson:"-"`
}

func (s Shared) GetId() string {
	return s.ID.Hex()
}

func (s Shared) SetUpdatedAt() {
	t := time.Now().UTC()
	s.UpdatedAt = &t
}

func (s Shared) SetID(id primitive.ObjectID) {
	s.ID = id
}

func (s Shared) Initialize(id primitive.ObjectID, now time.Time) {
	s.ID = id
	t := now.UTC()
	s.CreatedAt = &t
}

func (s Shared) SetUsedProjection(flag bool) {
	s.usedProjection = flag
}
func (s Shared) DidUseProjection() bool {
	return s.usedProjection
}

var _ SharedInterface = (*Shared)(nil)
