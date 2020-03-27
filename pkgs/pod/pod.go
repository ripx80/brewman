package pod

import (
	"fmt"
	"time"

	"github.com/ripx80/brewman/pkgs/brew"
	"github.com/ripx80/recipe"
)

/*
Todo: Remove idleStep and think about *step in task
*/

type Pod struct {
	Kettle *brew.Kettle
	recipe *recipe.Recipe
	task   *Task
	stop   chan struct{}
}

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
}

type PodMetric struct {
	StepName,
	Recipe string
	Step   StepMetric
	Kettle brew.KettleMetric
}

// func (pm *PodMetric) String() string{
// 	return fmt.Sprintf(

// 	)
// }

var idleStep = &Step{
	Name: "idle",
	F:    func() error { return nil },
}

func (p *Pod) Metric() *PodMetric {
	return &PodMetric{
		StepName: p.task.step.Name,
		Recipe:   p.recipe.Global.Name,
		Step:     p.task.step.Metric,
		Kettle:   p.Kettle.Metric(),
	}
}

func New(kettle *brew.Kettle, recipe *recipe.Recipe, stop chan struct{}) *Pod {
	return &Pod{
		Kettle: kettle,
		recipe: recipe,
		stop:   stop,
		task: &Task{
			Name: "Empty",
			step: idleStep,
		},
	}
}

/*Run loop over the steps and check return values*/
func (p *Pod) Run() error {
	for _, s := range p.task.Steps {
		s.Metric.Start = time.Now()
		s.Metric.TempStart, _ = p.Kettle.Temp.Get() // no buffer for first set
		p.task.step = s
		if err := s.F(); err != nil {
			return err
		}
		s.Metric.End = time.Now()
	}
	return nil
}

// func (p *Pod) Metric() *Metric {
// 	return &Metric{
// 		Pod:       p.Name,
// 		Step:      p.task.step.Name,
// 		StartTime: fmt.Sprintf(p.task.step.Start.Format("15:04:05")),
// 		Time,
// 		HoldTime,
// 		TempStart,
// 		Temp,
// 		TempEnd,
// 		State,
// 		Fail,*/
// 		Recipe: p.recipe.Global.Name,
// 	}
// }

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

/*StepAgitatorOn tuns the Agiator on if defined in config*/
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

//need a good way for a confirmation with channels

/*StepConfirm hold the temp as task step*/
// func (p *Pod) Confirm() *Step {
// 	return &Step{
// 		Name: "TempHold",
// 		F: func() error {
// 			return confirm("start mashing? <y/n>")
// 		}}
// }

/*Hotwater Task template*/
func (p *Pod) Hotwater(temp float64) {
	p.task = &Task{
		Name: "Hotwater",
		Steps: []*Step{
			p.StepTempUp("TempUp", temp),
			p.StepTempHold("TempHold", temp, 0),
		},
	}
	p.task.step = p.task.Steps[0]
}

/*Mash Task template*/
func (p *Pod) Mash(extendRest int) {
	task := &Task{
		Name: "Mash",
		Steps: []*Step{
			p.StepAgitatorOn(),
			p.StepTempUp("TempUp", p.recipe.Mash.InTemperatur),
			//p.Confirm(), Malt added
		},
	}
	for num, rast := range p.recipe.Mash.Rests {
		task.Steps = append(
			task.Steps,
			p.StepTempUp(fmt.Sprintf("Rest %d TempUp", num), rast.Temperatur), // mabye use extra field in Step for additional information like rest num?
			p.StepTempHold(fmt.Sprintf("Rest %d TempHold", num), rast.Temperatur, time.Duration(rast.Time*60)*time.Second),
		)
	}

	p.task = task
	p.task.step = p.task.Steps[0]
	//p.Confirm(), jod and if not correct append a new ExtendRest

}

/*MashRast can jump to a specific rast. Not Index Safe*/
func (p *Pod) MashRast(num int) {
	rast := p.recipe.Mash.Rests[num]
	p.task = &Task{
		Name: "MashRast",
		Steps: []*Step{
			p.StepAgitatorOn(),
			p.StepTempUp(fmt.Sprintf("Rest %d TempUp", num), rast.Temperatur),
			p.StepTempHold(fmt.Sprintf("Rest %d TemHold", num), rast.Temperatur, time.Duration(rast.Time*60)*time.Second),
		},
	}
	p.task.step = p.task.Steps[0]
}

/*Cook implements cooking programm*/
func (p *Pod) Cook(temp float64) {
	p.task = &Task{
		Name: "MashRast",
		Steps: []*Step{
			p.StepTempUp("TempUp", temp),
			p.StepTempHold("TempHold", temp, time.Duration(p.recipe.Cook.Time*60)*time.Second),
		},
	}
	p.task.step = p.task.Steps[0]
}

/*Validate Task template*/
func (p *Pod) Validate(temp float64) {
	p.task = &Task{
		Name: "Validate",
		Steps: []*Step{
			p.StepTempUp("TempUp", temp+1),
		},
	}
	p.task.step = p.task.Steps[0]
}
