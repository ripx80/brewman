package cmd

import (
	"fmt"
	"io/ioutil"
	"math"
	"sync"
	"time"

	"github.com/gdamore/tcell"
	log "github.com/ripx80/brave/log/logger"
	"github.com/ripx80/brewman/pkgs/pod"
	"github.com/rivo/tview"
)

const brand string = `  _________       ___.               ___.
  \_   ___ \___.__\_ |__   __________\_ |_________  ______  _  __
  /    \  \<   |  || __ \_/ __ \_  __ | __ \_  __ _/ __ \ \/ \/ /
  \     \___\___  || \_\ \  ___/|  | \| \_\ |  | \\  ___/\     /
   \______  / ____||___  /\___  |__|  |___  |__|   \___  >\/\_/
          \/\/         \/     \/          \/           \/ ripx80
`

/*
- add change recipe
- add locks on buffer access
- add wg group and stop on routines
- add finish Step to display
- bug: if you display logs in modal and a modal question appears no focus!
		use u.content to display logs

future
- Get list of jobs from pod (to see what will happen)
	- calculate time of ending
- at the moment we dont get short jobs like AgiatorOn
*/

type ui struct {
	content   *tview.Table
	container *tview.Flex
	app       *tview.Application
	left      *tview.Table
	right     *tview.TextView
	options   *tview.List
	commands  *tview.List
	modal     *tview.Modal
	buffers   [3]buffer
	active    uint
	instant   chan struct{}
}

type buffer struct {
	n string
	b []pod.PodMetric
	m sync.Mutex
}

func (b *buffer) Metric(m pod.PodMetric) {
	b.m.Lock()
	defer b.m.Unlock()

	l := len(b.b)

	if l == 0 && !m.Running {
		return
	}

	if l == 0 {
		b.b = append(b.b, m)
		return
	}

	if b.b[l-1].Step.Start != m.Step.Start {
		b.b[l-1].Step.End = time.Now() // its a hack
		b.b = append(b.b, m)
		return
	}
	b.b[l-1] = m
}

func (b *buffer) Clear() {
	b.b = []pod.PodMetric{}
}

func (u *ui) Metrics() {
	for {
		select {
		case <-time.After(500 * time.Millisecond):
			u.buffers[0].Metric(cfg.pods.hotwater.Metric())
			u.buffers[1].Metric(cfg.pods.masher.Metric())
			u.buffers[2].Metric(cfg.pods.cooker.Metric())
		}
	}
}

func (u *ui) refresh() {
	var quest pod.Quest
	for {
		select {
		case quest = <-cfg.confirm:
		case <-u.instant:
		case <-time.After(1 * time.Second):
		}
		if quest != (pod.Quest{}) {
			tp := u.commands
			if u.options.GetFocusable().HasFocus() {
				tp = u.options
			}

			u.modal = tview.NewModal().
				SetText(quest.Msg).
				AddButtons([]string{"Yes", "No"}).
				SetDoneFunc(func(buttonIndex int, buttonLabel string) {
					if buttonLabel == "Yes" {
						cfg.confirm <- pod.Quest{Msg: "y", Asw: true}
					}
					if buttonLabel == "No" {
						cfg.confirm <- pod.Quest{Msg: "n", Asw: false}
					}
					u.container.RemoveItem(u.modal)
					u.app.SetFocus(tp)
				})
			u.app.SetFocus(u.modal)
			u.container.AddItem(u.modal, 0, 1, true)
			quest = pod.Quest{}
			continue
		}
		u.app.QueueUpdateDraw(func() {
			u.Content()
		})
	}
}

func cell(t *tview.Table, rowCount int, num int, cfg *tview.TableCell, value string) {
	cell := *cfg // copy
	cell.Text = value
	t.SetCell(rowCount, num, &cell)
}

func row(t *tview.Table, rowCount int, content []string, cfg *tview.TableCell) {
	for idx, v := range content {
		cell(t, rowCount, idx, cfg, v)
	}
}

func timeString(t time.Time) string {
	return fmt.Sprintf(t.Format("15:04:05"))
}

func (u *ui) Content() {
	title := " Pod: [::b]%s(%s) "
	u.content.Clear().SetBorder(true)
	u.content.SetBorders(true)

	if len(u.buffers[u.active].b) == 0 {
		u.content.SetTitle(fmt.Sprintf(title, u.buffers[u.active].n, "S")).SetTitleAlign(1).SetTitleColor(tcell.ColorDeepPink)
		return
	}

	m := (u.buffers[u.active].b)[len(u.buffers[u.active].b)-1]
	run := "S"
	if m.Running {
		run = "R"
	}
	u.content.SetTitle(fmt.Sprintf(title, u.buffers[u.active].n, run)).SetTitleAlign(1).SetTitleColor(tcell.ColorDeepPink)
	row(u.content,
		u.content.GetRowCount(),
		[]string{
			"[::b]Step",
			"[::b]StartTime",
			"[::b]Time",
			"[::b]HoldTime",
			"[::b]TempStart",
			"[::b]Temp",
			"[::b]TempEnd",
			"[::b]State",
			"[::b]Fail",
		},
		&tview.TableCell{
			Expansion: 1,
			Align:     tview.AlignCenter,
			Color:     tcell.ColorAqua,
		})

	now := time.Now()
	var runtime time.Duration
	for _, m := range u.buffers[u.active].b {
		if m.Step.End.IsZero() {
			runtime = now.Sub(m.Step.Start).Round(time.Second)
		} else {
			runtime = m.Step.End.Sub(m.Step.Start).Round(time.Second)
		}
		var hold string
		if m.Step.Hold == 0 {
			hold = "-"
		} else {
			hold = fmt.Sprintf("%02d:%02.f", m.Step.Hold/time.Minute, math.Mod(m.Step.Hold.Seconds(), 60))
		}

		var state string
		if m.Kettle.Heater == false {
			state = fmt.Sprintf("[#8080ff::b]%t[#ffffff]", m.Kettle.Heater)
		} else {
			state = fmt.Sprintf("[red::b]%t[white]", m.Kettle.Heater)
		}

		row(u.content, u.content.GetRowCount(), []string{
			m.StepName,
			timeString(m.Step.Start),
			fmt.Sprintf("%02d:%02.f", runtime/time.Minute, math.Mod(runtime.Seconds(), 60)),
			hold,
			fmt.Sprintf("%0.2f", m.Step.TempStart),
			fmt.Sprintf("%0.2f", m.Kettle.Temp),
			fmt.Sprintf("%0.2f", m.Step.TempEnd),
			state,
			fmt.Sprintf("%d", m.Kettle.Fail),
		},
			&tview.TableCell{
				Expansion: 1,
				Align:     tview.AlignCenter,
				Color:     tcell.ColorYellow,
			},
		)
	}
}

func logo(brand string) *tview.TextView {
	return tview.NewTextView().
		SetText(brand).
		SetTextColor(tcell.ColorRed).
		SetTextAlign(tview.AlignLeft)
}

// inefficent buf as arg use a choice
func (u *ui) Commands(name string, pod *pod.Pod) {
	u.commands.Clear().
		AddItem("Run", "start pod", 'r', func() {
			if !pod.Metric().Running {
				u.content.Clear()
				// clear buffer
				u.buffers[u.active].Clear()
				go func() {
					defer cfg.wg.Done()
					cfg.wg.Add(1)
					if err := pod.Run(); err != nil {
						log.WithFields(log.Fields{
							"pod":   name,
							"error": err,
						}).Error("pod job run failed")
					}
					//draw finish row
				}()
				u.Instant()
			}

		}).
		AddItem("Stop", "stop pod", 's', func() {
			pod.Stop()
		}).
		// AddItem("Recipe", "change recipe", 'c', func() {}).
		AddItem("Back", "go back menu", 'b', func() {
			u.commands.Clear()
			u.app.SetFocus(u.options)
		})
	u.app.SetFocus(u.commands)
}

func (u *ui) Logs(fp string) {
	content, err := ioutil.ReadFile(fp)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("can not read from file")
	}
	u.modal = tview.NewModal().
		SetText(fmt.Sprintf("%s", content)).
		AddButtons([]string{"close"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			u.container.RemoveItem(u.modal)
			u.app.SetFocus(u.options)
		})

	u.app.SetFocus(u.modal)
	u.container.AddItem(u.modal, 0, 1, true)
}

func (u *ui) Instant() {
	u.instant <- struct{}{}
}

func (u *ui) Options() *tview.List {
	return tview.NewList().
		AddItem("Hotwater", "Select Hotwater Pod", 'h', func() {
			u.Commands("hotwater", cfg.pods.hotwater)
			u.active = 0
			u.Instant()

		}).
		AddItem("Masher", "Select Masher pod", 'm', func() {
			u.Commands("masher", cfg.pods.masher)
			u.active = 1
			u.Instant()
		}).
		AddItem("Cooker", "Select Cooker pod", 'c', func() {
			u.Commands("cooker", cfg.pods.cooker)
			u.active = 2
			u.Instant()
		}).
		AddItem("Logs", "Logs", 'l', func() {
			u.Logs(logfile)
		}).
		AddItem("Quit", "Press to exit", 'q', func() {
			u.app.Stop()
		}).
		SetWrapAround(true).
		ShowSecondaryText(false)
}

func view() error {
	var view ui
	view = ui{
		content:   tview.NewTable(),
		container: tview.NewFlex(),
		app:       tview.NewApplication(),
		left:      tview.NewTable(),
		right:     logo(brand),
		options:   tview.NewList(),
		commands:  tview.NewList().ShowSecondaryText(false),
		modal:     tview.NewModal(),
		buffers: [3]buffer{
			buffer{n: "Hotwater", b: []pod.PodMetric{}, m: sync.Mutex{}},
			buffer{n: "Masher", b: []pod.PodMetric{}, m: sync.Mutex{}},
			buffer{n: "Cooker", b: []pod.PodMetric{}, m: sync.Mutex{}},
		},
		instant: make(chan struct{}),
		active:  0,
	}
	view.options = view.Options()

	leftCfg := &tview.TableCell{Expansion: 0, Align: tview.AlignLeft, Color: tcell.ColorAqua}

	row(view.left, view.left.GetRowCount(), []string{"[::b]Version: ", version}, leftCfg)
	row(view.left, view.left.GetRowCount(), []string{"[::b]Recipe: ", cfg.recipe.Global.Name}, leftCfg)

	view.container = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(view.left, 0, 1, false).
			AddItem(view.options, 0, 1, true).
			AddItem(view.commands, 0, 1, true).
			AddItem(view.right, 0, 3, false), 0, 1, true).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(view.content, 0, 1, false), 0, 3, false)

	go view.Metrics()
	go view.refresh()
	view.instant <- struct{}{}
	return view.app.SetRoot(view.container, true).Run()
}
