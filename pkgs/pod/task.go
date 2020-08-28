package pod

import (
	"time"
)

/*
Todo: Remove idleStep and think about *step in task
change *Pod to *Task for steps
*/
var idleStep = &Step{
	Name: "idle",
	F:    func() error { return nil },
}

/*StepMetric represents the Metrics of each step*/
type StepMetric struct {
	Start,
	End time.Time
	Hold time.Duration
	TempStart,
	TempEnd float64
}

/*Step has a name and a call function which will be call by task.Run() in sequence*/
type Step struct {
	Name   string
	F      func() error
	Metric StepMetric
}

/*Task holds multiple steps and executes in sequence*/
type Task struct {
	Name  string
	Steps []*Step
	step  *Step //actual step working on
	num   int   // if i know the num i dont need the step pointer
}

/*Quest stimple struct with question and answer*/
type Quest struct {
	Msg string
	Asw bool
}

/*StepAgitatorOn tuns the Agiator on if defined in config*/
func (p *Pod) StepAgitatorOn() *Step {
	return &Step{
		Name: "AgiatorOn",
		F: func() error {
			if p.Kettle.Agitator != nil && !p.Kettle.Agitator.State() {
				p.Kettle.Agitator.On()
			}
			return nil
		}}
}

/*StepAgitatorOff tuns the Agiator on if defined in config*/
func (p *Pod) StepAgitatorOff() *Step {
	return &Step{
		Name: "AgiatorOff",
		F: func() error {
			if p.Kettle.Agitator != nil && p.Kettle.Agitator.State() {
				p.Kettle.Agitator.Off()
			}
			return nil
		}}
}

/*StepTempUp increase the temp as task step*/
func (p *Pod) StepTempUp(name string, temp float64) *Step {
	return &Step{
		Name: name,
		F: func() error {
			return p.Kettle.TempUp(p.stop, temp)
		},
		Metric: StepMetric{
			TempEnd: temp,
		},
	}
}

/*StepTempHold hold the temp as task step*/
func (p *Pod) StepTempHold(name string, temp float64, time time.Duration) *Step {
	return &Step{
		Name: name,
		F: func() error {
			return p.Kettle.TempHold(p.stop, temp, time)
		},
		Metric: StepMetric{
			TempEnd: temp,
			Hold:    time,
		}}
}

/*StepConfirm hold the temp as task step*/
func (p *Pod) StepConfirm(msg string, confirm chan Quest) *Step {
	return &Step{
		Name: "Confirm",
		F: func() error {
			confirm <- Quest{Msg: msg, Asw: false}
			quest := <-confirm // will blocks on stop
			if !quest.Asw {
				p.task.Steps = append(p.task.Steps, nil /* use the zero value of the element type */)
				copy(p.task.Steps[p.task.num+1:], p.task.Steps[p.task.num:])
				p.task.Steps[p.task.num] = p.StepConfirm(msg, confirm)
			}
			return nil
		},
	}
}

/*StepConfirmRest hold the temp as task step*/
func (p *Pod) StepConfirmRest(msg string, confirm chan Quest, extendRest int, temp float64) *Step {
	return &Step{
		Name: "ConfirmRest",
		F: func() error {
			confirm <- Quest{Msg: msg, Asw: false}
			quest := <-confirm
			if !quest.Asw {
				p.task.Steps = append(p.task.Steps[:p.task.num+1],
					append(
						[]*Step{
							p.StepTempHold("TempHold", temp, time.Duration(extendRest*60)*time.Second),
							p.StepConfirmRest(msg, confirm, extendRest, temp),
						},
						p.task.Steps[p.task.num+1:]...)...)
			}
			return nil
		},
	}
}

/*StepMessage set a message as name like finish or start
use sleep for reading buffers (ui)
*/
func (p *Pod) StepMessage(msg string) *Step {
	return &Step{
		Name: msg,
		F:    func() error { time.Sleep(1 * time.Second); return nil },
	}
}
