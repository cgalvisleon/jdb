package jdb

import (
	"fmt"

	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/event"
	"github.com/cgalvisleon/et/timezone"
	"github.com/cgalvisleon/et/utility"
)

/**
* eventEmiter
* @param message event.Message
**/
func (s *Model) eventEmiter(message event.Message) {
	if s.eventsEmiter == nil {
		s.eventsEmiter = make(map[string]event.Handler)
	}

	eventEmiter, ok := s.eventsEmiter[message.Channel]
	if !ok {
		console.Alert(fmt.Errorf(MSG_EVENT_NOT_FOUND, message.Channel, s.Name))
		return
	}

	eventEmiter(message)
}

/**
* On
* @param channel string, handler EventHandler
**/
func (s *Model) On(channel string, handler event.Handler) *Model {
	if s.eventsEmiter == nil {
		s.eventsEmiter = make(map[string]event.Handler)
	}

	s.eventsEmiter[channel] = handler

	return s
}

/**
* Emit
* @param channel string, data et.Json
**/
func (s *Model) Emit(channel string, data et.Json) *Model {
	if s.eventEmiterChannel == nil {
		console.Alert(fmt.Errorf("event channel not found (%s)", channel))
	}

	message := event.Message{
		CreatedAt: timezone.NowTime(),
		FromId:    s.Db.Id,
		Id:        utility.UUID(),
		Channel:   channel,
		Data:      data,
	}

	s.eventEmiterChannel <- message

	return s
}
