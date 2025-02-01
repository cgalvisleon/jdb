package jdb

import "github.com/cgalvisleon/et/console"

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

func (s *Flow) Run(status chan *Result) {
	if status == nil {
		status = make(chan *Result)
		defer close(status)
	}

	for i, step := range s.Steps {
		go step.Execute(status)

		result := <-status
		if !result.Ok {
			console.Logf("Flow", `❌ Error in step: %s - %s, Begin rollback...`, step.Name, result.Message)
			s.rollback(i)
			return
		} else {
			console.Logf("Flow", `✅ Complete step: %s`, step.Name)
		}
	}

	console.Logf("Flow", `🎉 Executed successfully flow: %s`, s.Name)
}

func (s *Flow) rollback(index int) {
	for i := index; i >= 0; i-- {
		if s.Steps[i].Rollback != nil {
			console.Logf("Flow", `Rollback step: (%d) %s`, i, s.Steps[i].Name)
			s.Steps[i].Rollback()
		}
	}
	console.Logf("Flow", `Rollback completed`)
}
