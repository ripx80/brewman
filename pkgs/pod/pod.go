package pod

import (
	"fmt"
	"time"

	"github.com/ripx80/brewman/pkgs/brew"
	"github.com/ripx80/recipe"
)

type Pod struct {
	Kettle *brew.Kettle
	recipe *recipe.Recipe
	task   *Task
	stop   chan struct{}
	run    bool
}

type PodMetric struct {
	StepName,
	Recipe string
	Running bool
	Step    StepMetric
	Kettle  brew.KettleMetric
}

func (p *Pod) Metric() PodMetric {
	return PodMetric{
		StepName: p.task.step.Name,
		Running:  p.run,
		Recipe:   p.recipe.Global.Name,
		Step:     p.task.step.Metric,
		Kettle:   p.Kettle.Metric(),
	}
}

func (p *Pod) Stop() {
	close(p.stop)
	p.stop = make(chan struct{}) // better to send something than stop?
}

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
	defer func() { p.run = false }()
	var s *Step
	p.run = true
	// num of steps can change, dont use range
	for i := 0; i < len(p.task.Steps); i++ {
		p.task.num = i
		s = p.task.Steps[i]
		s.Metric.Start = time.Now()
		s.Metric.TempStart, _ = p.Kettle.Temp.Get() // no buffer for first set
		p.task.step = s

		if err := s.F(); err != nil {
			s.Metric.End = time.Now()
			if err.Error() == "cancled" {
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
			p.StepTempUp(fmt.Sprintf("Rest %d TempUp", num), rast.Temperatur), // mabye use extra field in Step for additional information like rest num?
			p.StepTempHold(fmt.Sprintf("Rest %d TempHold", num), rast.Temperatur, time.Duration(rast.Time*60)*time.Second),
		)
		if len(p.recipe.Mash.Rests)-2 == num {
			task.Steps = append(
				task.Steps,
				p.StepConfirmRest("jod test successful?", confirm, extendRest, rast.Temperatur),
			)
		}
	}
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
