package main

import (
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication().EnableMouse(true)

	list := tview.NewForm()
	cfg := tview.NewForm()
	root := tview.NewFlex().
		AddItem(cfg, 0, 1, true).
		AddItem(list, 0, 1, false)

	list.
		SetItemPadding(0).
		SetBorder(true).
		SetTitle("Peak list")

	cfg.
		AddInputField("* File name", "peeklist.csv", 30, nil, nil).
		AddDropDown("* Phase", []string{"QC", "1/1AC"}, 0, nil).
		AddInputField("* Lattice constant", "", 10, tview.InputFieldFloat, nil).
		AddInputField("* X-ray wave length", "1.540593", 10, tview.InputFieldFloat, nil).
		AddButton("Calc", func() {
			ph, _ := cfg.GetFormItem(1).(*tview.DropDown).GetCurrentOption()
			lc, _ := strconv.ParseFloat(cfg.GetFormItem(2).(*tview.InputField).GetText(), 64)
			wl, _ := strconv.ParseFloat(cfg.GetFormItem(3).(*tview.InputField).GetText(), 64)
			go func() {
				var labels *[]string
				if ph == 0 {
					labels, _ = calcQcPeak(lc, wl)
				} else {
					labels, _ = calcAcPeak(lc, wl)
				}

				app.QueueUpdateDraw(func() {
					list.Clear(true)
					for _, l := range *labels {
						list.AddInputField(l, "", 10, tview.InputFieldFloat, nil)
					}
					app.SetFocus(list)
				})
			}()

		}).
		AddButton("Save", nil).
		AddButton("Quit", func() {
			app.Stop()
		}).
		SetItemPadding(0)
	cfg.
		SetBorder(true).
		SetTitle("Config")

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Modifiers() == tcell.ModAlt {
			switch event.Rune() {
			case 'l':
				app.SetFocus(list)
				return nil
			case 'h':
				app.SetFocus(cfg)
				return nil
			}
		}

		return event
	})

	if err := app.SetRoot(root, true).Run(); err != nil {
		panic(err)
	}
}
