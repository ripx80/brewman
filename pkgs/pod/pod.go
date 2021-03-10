package pod

import (
	"fmt"
	"time"

	"github.com/ripx80/brewman/pkgs/brew"
	"github.com/ripx80/recipe"
)

/*Pod hold a kettle with recipe and task scheduler*/
type Pod struct {
	Kettle *brew.Kettle
	recipe *recipe.Recipe
	task   *Task
	stop   chan struct{}
	run    bool
}

/*Metric provide metrics for the pod*/
type Metric struct {
	StepName,
	Recipe string
	Running bool
	Step    StepMetric
	Kettle  brew.KettleMetric
}

/*Metric return a Metric struct with metrics*/
func (p *Pod) Metric() Metric {
	return Metric{
		StepName: p.task.step.Name,
		Running:  p.run,
		Recipe:   p.recipe.Global.Name,
		Step:     p.task.step.Metric,
		Kettle:   p.Kettle.Metric(),
	}
}

/*Stop the pod*/
func (p *Pod) Stop() {
	close(p.stop)
	p.stop = make(chan struct{}) // better to send something than stop?
}

/*New create a new pod*/
func New(kettle *brew.Kettle, recipe *recipe.Recipe) *Pod {
	return &Pod{
		Kettle: kettle,
		recipe: recipe,
		stop:   make(chan struct{}),
		task: &Task{
			Name: "Empty",
			step: idleStep,
		},
	}
}

/*Jobs return jobs in task list. can change during runtime*/
func (p *Pod) Jobs() map[string]StepMetric {
	m := make(map[string]StepMetric)
	for _, v := range p.task.Steps {
		m[v.Name] = v.Metric
	}
	return m
}

/*Run loop over the steps and check return values*/
func (p *Pod) Run() error {
	defer func() {
		p.run = false
		p.Stop()
	}()
	var s *Step
	p.run = true
	// num of steps can change, dont use range
	for i := 0; i < len(p.task.Steps); i++ {
		p.task.num = i
		s = p.task.Steps[i]
		s.Metric.Start = time.Now()
		s.Metric.End = time.Time{}
		s.Metric.TempStart, _ = p.Kettle.Temp.Get() // no buffer for first set
		p.task.step = s

		if err := s.F(); err != nil {
			s.Metric.End = time.Now()
			if err.Error() == brew.CancelErr {
				return nil
			}
			return err
		}
		s.Metric.End = time.Now()
	}
	return nil
}

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
func (p *Pod) Mash(extendRest int, confirm chan Quest) {
	task := &Task{
		Name: "Mash",
		Steps: []*Step{
			p.StepAgitatorOn(),
			p.StepTempUp("TempUp", p.recipe.Mash.InTemperatur),
			p.StepConfirm("Malt added? continue...", confirm),
		},
	}

	for num, rast := range p.recipe.Mash.Rests {
		task.Steps = append(
			task.Steps,
			p.StepTempUp(fmt.Sprintf("Rest %d TempUp", num+1), rast.Temperatur), // mabye use extra field in Step for additional information like rest num?
			p.StepTempHold(fmt.Sprintf("Rest %d TempHold", num+1), rast.Temperatur, time.Duration(rast.Time*60)*time.Second),
		)
	}
	lastRast := len(p.recipe.Mash.Rests) - 1
	pos := 0
	if p.recipe.Mash.Rests[lastRast].Temperatur >= p.recipe.Mash.OutTemperatur {
		pos = 2 // two steps back
		lastRast--
	}
	idx := len(task.Steps) - pos
	task.Steps = append(task.Steps, &Step{})
	copy(task.Steps[idx+1:], task.Steps[idx:])
	task.Steps[idx] = p.StepConfirmRest("jod test successful?", confirm, extendRest, p.recipe.Mash.Rests[lastRast].Temperatur)

	task.Steps = append(
		task.Steps,
		p.StepTempUp("Rest Out TempUp", p.recipe.Mash.OutTemperatur),
		p.StepMessage("Successful mash"),
	)
	p.task = task
	p.task.step = p.task.Steps[0]
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
