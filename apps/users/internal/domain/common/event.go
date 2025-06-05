package common

type Event interface {
	EventName() string
}

type EventBus interface {
	Publish(event Event) error
}
