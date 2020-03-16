package main

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

const refreshInterval = 500 * time.Millisecond

var (
	pod *tview.Table
	app *tview.Application
)

func refresh() {
	tick := time.NewTicker(refreshInterval)
	cfg := &tview.TableCell{Expansion: 1, Align: tview.AlignCenter, Color: tcell.ColorYellow}
	for {
		select {
		case <-tick.C:
			//app.QueueUpdateDraw(func() { drawLastRow(pod) })
			now := time.Now()

			app.QueueUpdateDraw(func() {
				drawCell(pod, (pod.GetRowCount() - 1), 2, cfg, fmt.Sprintf(now.Format("15:04:05")))
			})
		}
	}
}

func drawCell(t *tview.Table, rowCount int, num int, cfg *tview.TableCell, value string) {
	cell := *cfg // copy
	cell.Text = value
	t.SetCell(rowCount, num, &cell)
}

func drawRow(t *tview.Table, rowCount int, content []string, cfg *tview.TableCell) {
	for idx, v := range content {
		drawCell(t, rowCount, idx, cfg, v)
	}
}

func getStringTime() string {
	t := time.Now()
	return fmt.Sprintf(t.Format("15:04:05"))
}

const logo string = `  _________       ___.               ___.
  \_   ___ \___.__\_ |__   __________\_ |_________  ______  _  __
  /    \  \<   |  || __ \_/ __ \_  __ | __ \_  __ _/ __ \ \/ \/ /
  \     \___\___  || \_\ \  ___/|  | \| \_\ |  | \\  ___/\     /
   \______  / ____||___  /\___  |__|  |___  |__|   \___  >\/\_/
          \/\/         \/     \/          \/           \/ ripx80
`

func main() {

	leftCfg := &tview.TableCell{Expansion: 0, Align: tview.AlignCenter, Color: tcell.ColorYellow}

	podName := "Hotwater"
	app = tview.NewApplication()
	timeNow := getStringTime()

	left := tview.NewTable().SetBorders(true)
	drawRow(left, left.GetRowCount(), []string{"Version: ", "1.0"}, leftCfg)
	drawRow(left, left.GetRowCount(), []string{"Recipe: ", "TagIPA"}, leftCfg)

	// add as list pods for Hotwater, Masher, Cooker so you can switch.
	// need to set the reciept
	middle := tview.NewList().
		// AddItem("<?>", "Help", '?', nil).
		// AddItem("<l>", "Logs", 'l', nil).
		AddItem("Quit", "Press to exit", 'q', func() {
			app.Stop()
		})

	middle.SetWrapAround(true)
	middle.ShowSecondaryText(false)

	logoBox := tview.NewTextView()
	logoBox.SetText(logo)
	logoBox.SetTextColor(tcell.ColorDarkRed)
	logoBox.SetTextAlign(tview.AlignLeft)

	pod = tview.NewTable()
	pod.SetBorder(true)
	pod.SetBorders(true)
	pod.SetTitle(fmt.Sprintf("  Pod: [::b]%s ", podName)).SetTitleAlign(1).SetTitleColor(tcell.ColorDarkRed)

	drawRow(pod, pod.GetRowCount(), []string{"[::b]Step", "[::b]StartTime", "[::b]Time", "[::b]HoldTime", "[::b]TempStart", "[::b]Temp", "[::b]TempEnd", "[::b]State", "[::b]Fail"}, &tview.TableCell{Expansion: 1, Align: tview.AlignCenter, Color: tcell.ColorAqua})
	drawRow(pod, pod.GetRowCount(), []string{"increase", timeNow, timeNow, "60:00", "43.23", "53.34", "62.00", "on", "0"}, &tview.TableCell{Expansion: 1, Align: tview.AlignCenter, Color: tcell.ColorYellow})

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(left, 0, 1, false).
			AddItem(middle, 0, 1, true).
			AddItem(logoBox, 0, 2, false), 0, 1, true).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(pod, 0, 1, false), 0, 3, false)

	//go refresh()

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}
