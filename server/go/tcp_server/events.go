package tcp_server

import "time"

type Event struct {
	time      time.Time
	message   string
	level     int
	isAnError bool
}

func NewEvent(message string, isAnError bool, level int) *Event {
	return &Event{
		time:      time.Now(),
		message:   message,
		isAnError: isAnError,
		level:     level,
	}
}

func (e *Event) Time() time.Time {
	return e.time
}

func (e *Event) Message() string {
	return e.message
}

func (e *Event) IsAnError() bool {
	return e.isAnError
}

func (e *Event) Level() int {
	return e.level
}
