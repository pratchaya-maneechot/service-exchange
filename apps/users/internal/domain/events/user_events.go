package events

import "github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/entities"

type UserCreated entities.User

var UserCreatedEvent = "UserCreated"

func (e UserCreated) EventName() string {
	return UserCreatedEvent
}

type UserLineProfileLinked entities.LineProfile

var UserLineProfileLinkedEvent = "UserLineProfileLinked"

func (e UserLineProfileLinked) EventName() string {
	return UserLineProfileLinkedEvent
}
