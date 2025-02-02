package jdb

import (
	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/utility"
)

type Result struct {
	Ok      bool
	Message string
	Data    interface{}
}

type Step struct {
	Flow        *Flow
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Execute     func(chan *Result) *Result
	Rollback    func()
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

func (s *Flow) Run(status chan *Result) {
	if status == nil {
		status = make(chan *Result)
		defer close(status)
	}

	for i, step := range s.Steps {
		go step.Execute(status)

		result := <-status
		if !result.Ok {
			console.Logf("Flow", "\u274C Error in step: %s - %s, Begin rollback...", step.Name, result.Message)
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
