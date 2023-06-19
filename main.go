package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	var labels *[]Mirror
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
				if ph == 0 {
					labels, _ = calcQcPeak(lc, wl)
				} else {
					labels, _ = calcAcPeak(lc, wl)
				}

				app.QueueUpdateDraw(func() {
					list.Clear(true)
					for _, m := range *labels {
						if ph == 0 {
							list.AddInputField(
								fmt.Sprintf("(%s, %s, %s, %s, %s, %s) ~%.2f: ", m.h, m.k, m.l, m.m, m.n, m.o, calcTowTheta(wl, m.N, lc)), "", 10, tview.InputFieldFloat, nil)
						} else {
							list.AddInputField(fmt.Sprintf("(%s, %s, %s) ~%.2f: ", m.h, m.k, m.l, calcTowTheta(wl, m.N, lc)), "", 10, tview.InputFieldFloat, nil)
						}

					}
					app.SetFocus(list)
				})
			}()

		}).
		AddButton("Save", func() {
			name := cfg.GetFormItem(0).(*tview.InputField).GetText()
			ph, _ := cfg.GetFormItem(1).(*tview.DropDown).GetCurrentOption()
			wl, _ := strconv.ParseFloat(cfg.GetFormItem(3).(*tview.InputField).GetText(), 64)

			go func() {
				f, err := os.Create(name)
				if err != nil {
					panic(err)
				}
				defer f.Close()

				w := csv.NewWriter(f)

				// header
				if ph == 0 {
					if err := w.Write([]string{"#h", "k", "l", "m", "n", "o", "2theta", "NR", "lattice constant"}); err != nil {
						panic(err)
					}
				} else {
					if err := w.Write([]string{"#h", "k", "l", "2theta", "NR", "lattice constant"}); err != nil {
						panic(err)
					}
				}

				for i, m := range *labels {
					txt := list.GetFormItem(i).(*tview.InputField).GetText()
					if txt == "" {
						continue
					}

					th, _ := strconv.ParseFloat(txt, 64)
					if ph == 0 {
						if err := w.Write([]string{m.h, m.k, m.l, m.m, m.n, m.o, txt, calcNR(th), calcLatticeConstant(m.N, wl, th)}); err != nil {
							panic(err)
						}
					} else {
						if err := w.Write([]string{m.h, m.k, m.l, txt, calcNR(th), calcLatticeConstant(m.N, wl, th)}); err != nil {
							panic(err)
						}
					}
				}

				w.Flush()
			}()
		}).
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
