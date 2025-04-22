package jdb

import (
	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/reg"
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

type Function struct {
	Flow        *Flow
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Code        string `json:"code"`
	Context     *Context
	Execute     func(chan *Context) *Context
	Rollback    func()
}

/**
* NewFunction
* @param name string, description string, execute func(chan *Context) *Context, rollback func()
* @return *Function
**/
func NewFunction(name, description string, execute func(chan *Context) *Context, rollback func()) *Function {
	return &Function{
		Id:          reg.Id("function"),
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
	Functions   []*Function
}

/**
* NewFlow
* @param name string
* @param description string
* @return *Flow
**/
func NewFlow(name, description string) *Flow {
	return &Flow{
		Id:          reg.Id("flow"),
		Name:        name,
		Description: description,
		Functions:   []*Function{},
	}
}

/**
* Describe
* @return et.Json
**/
func (s *Flow) Describe() et.Json {
	return et.Json{
		"id":          s.Id,
		"name":        s.Name,
		"description": s.Description,
		"functions":   s.Functions,
	}
}

/**
* Run
* @param status chan *Context
**/
func (s *Flow) Run(status chan *Context) {
	if status == nil {
		status = make(chan *Context)
		defer close(status)
	}

	for i, function := range s.Functions {
		go function.Execute(status)

		result := <-status
		if !result.Ok {
			console.Logf("Flow", "\u274C Error in step: %s - %s, Begin rollback...", function.Name, result.Message.Message)
			s.Rollback(i)
			return
		} else {
			console.Logf("Flow", "\u2705 Complete step: %s", function.Name)
		}
	}

	console.Logf("Flow", "\U0001F389 Executed successfully flow: %s", s.Name)
}

/**
* Rollback
* @param index int
**/
func (s *Flow) Rollback(index int) {
	for i := index; i >= 0; i-- {
		if s.Functions[i].Rollback != nil {
			console.Logf("Flow", "\u21A9 Rollback step: (%d) %s", i, s.Functions[i].Name)
			s.Functions[i].Rollback()
		}
	}

	console.Logf("Flow", "\U0001F504 Rollback completed")
}
