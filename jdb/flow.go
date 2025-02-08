package jdb

import (
	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/utility"
)

type Message struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Context struct {
	Ok      bool        `json:"ok"`
	Message Message     `json:"message"`
	Data    interface{} `json:"data"`
}

/**
* NewContext
* @param data interface{}
* @return *Context
**/
func NewContext(data interface{}) *Context {
	return &Context{
		Message: Message{},
	}
}

type Step struct {
	Flow        *Flow
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	context     *Context
	Execute     func(chan *Context) *Context
	Rollback    func()
}

/**
* NewStep
* @param name string, description string, execute func(chan *Context) *Context, rollback func()
* @return *Step
**/
func NewStep(name, description string, execute func(chan *Context) *Context, rollback func()) *Step {
	return &Step{
		Id:          utility.RecordId("step", ""),
		Name:        name,
		Description: description,
		Execute:     execute,
		Rollback:    rollback,
	}
}

type Flow struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Steps       []*Step
}

/**
* NewFlow
* @param name string
* @param description string
* @return *Flow
**/
func NewFlow(name, description string) *Flow {
	return &Flow{
		Id:          utility.RecordId("flow", ""),
		Name:        name,
		Description: description,
		Steps:       []*Step{},
	}
}

func (s *Flow) Run(status chan *Context) {
	if status == nil {
		status = make(chan *Context)
		defer close(status)
	}

	for i, step := range s.Steps {
		go step.Execute(status)

		result := <-status
		if !result.Ok {
			console.Logf("Flow", "\u274C Error in step: %s - %s, Begin rollback...", step.Name, result.Message.Message)
			s.rollback(i)
			return
		} else {
			console.Logf("Flow", "\u2705 Complete step: %s", step.Name)
		}
	}

	console.Logf("Flow", "\U0001F389 Executed successfully flow: %s", s.Name)
}

func (s *Flow) rollback(index int) {
	for i := index; i >= 0; i-- {
		if s.Steps[i].Rollback != nil {
			console.Logf("Flow", "\u21A9 Rollback step: (%d) %s", i, s.Steps[i].Name)
			s.Steps[i].Rollback()
		}
	}
	console.Logf("Flow", "\U0001F504 Rollback completed")
}
