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

type EventBase struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Master      uuid.UUID          `json:"master" bson:"master"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"desc" bson:"desc"`
	Type        string             `json:"type" bson:"type"`
	Prior       Priority           `json:"prior" bson:"prior"`
	Start       time.Time          `json:"start" bson:"start"`
	End         time.Time          `json:"end" bson:"end"`
}

type Participant struct {
	Accepted bool      `json:"accepted,omitempty" bson:"accepted,omitempty"`
	UID      uuid.UUID `json:"uid" bson:"uid"`
}

type ParticipantWithEmail struct {
	Accepted bool   `json:"accepted"`
	Email    string `json:"email"`
}

type Event struct {
	EventBase    `json:",inline" bson:",inline"`
	Participants []Participant `json:"parts" bson:"parts"`
}

func (e *Event) ParticipantsUUIDS() []uuid.UUID {
	result := make([]uuid.UUID, len(e.Participants))
	for _, p := range e.Participants {
		result = append(result, p.UID)
	}
	return result
}

type EventWithMails struct {
	EventBase    `json:",inline" bson:",inline"`
	Participants []string `json:"parts" bson:"parts"`
}

type Chat struct {
	EventID  primitive.ObjectID `json:"event_id" bson:"event_id"`
	Messages []MessageWithMail  `json:"messages" bson:"messages"`
}

type Message struct {
	Sender  uuid.UUID `json:"sender" bson:"sender"`
	Content string    `json:"cont" bson:"cont"`
}

type MessageWithMail struct {
	Sender  string `json:"sender" bson:"sender"`
	Content string `json:"cont" bson:"cont"`
}

type FileLoad struct {
	Name string `json:"name"`
	Data []byte `json:"data"`
}

type FileDownload struct {
	Name string `json:"name"`
	Link string `json:"link"`
}

type Notification struct {
	To      []string `json:"emails"`
	Subject string   `json:"subject"`
	Content string   `json:"content"`
}
