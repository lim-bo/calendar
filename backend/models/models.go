package models

import (
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserCredentials struct {
	Email string `json:"mail,omitempty"`
	Pass  string `json:"pass,omitempty"`
}

type UserCredentialsRegister struct {
	UserCredentials `json:",inline"`
	FirstName       string `json:"f_name"`
	SecondName      string `json:"s_name"`
	ThirdName       string `json:"t_name,omitempty"`
	Department      string `json:"dep"`
	Position        string `json:"pos"`
}

type Priority int8

var (
	PriorityHigh = Priority(3)
	PriorityMid  = Priority(2)
	PriorityLow  = Priority(1)
)

type Event struct {
	ID           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Master       uuid.UUID          `json:"master" bson:"master"`
	Name         string             `json:"name" bson:"name"`
	Description  string             `json:"desc" bson:"desc"`
	Type         string             `json:"type" bson:"type"`
	Prior        Priority           `json:"prior" bson:"prior"`
	Start        time.Time          `json:"start" bson:"start"`
	End          time.Time          `json:"end" bson:"end"`
	Participants []uuid.UUID        `json:"parts" bson:"parts"`
}
