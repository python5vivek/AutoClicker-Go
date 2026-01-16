package main

import (
	"strconv"
	"sync/atomic"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
)

var clicking atomic.Bool
var latency atomic.Int64
var side atomic.Bool

func clickLoop() {
	for {
		if clicking.Load() {
			sid := side.Load()
			if sid {
				robotgo.Click("right")
			} else {
				robotgo.Click("left")
			}

			time.Sleep(time.Duration(latency.Load()) * time.Millisecond)
		} else {
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func main() {
	go clickLoop()
	latency.Store(10)
	side.Store(false)

	hook.Register(hook.KeyDown, []string{"g"}, func(e hook.Event) {
		clicking.Store(true)
		side.Store(false)
	})
	hook.Register(hook.KeyDown, []string{"h"}, func(e hook.Event) {
		clicking.Store(true)
		side.Store(true)
	})
	hook.Register(hook.KeyDown, []string{"j"}, func(e hook.Event) {
		clicking.Store(false)
	})
	a := app.New()
	w := a.NewWindow("Go AutoClicker")
	Title := widget.NewRichText(
		&widget.TextSegment{
			Text: "Fyne AutoClicker",
			Style: widget.RichTextStyle{
				SizeName: theme.SizeNameHeadingText,
			},
		},
	)
	GLabel := widget.NewLabel("'G' ----> Left Click")
	HLabel := widget.NewLabel("'H' ----> Right Click")
	JLabel := widget.NewLabel("'J' ----> Stop Clicking")
	slider := widget.NewSlider(0.0, 100.0)
	slider.SetValue(float64(latency.Load()))
	stoptime := widget.NewLabel("Clicking Latency (ms):" + strconv.FormatInt(latency.Load(), 10))
	slider.OnChanged = func(f float64) {
		latency.Store(int64(f))
		stoptime.SetText("Clicking Latency (ms):" + strconv.FormatInt(latency.Load(), 10))
	}
	w.SetContent(container.NewVBox(
		Title,
		GLabel,
		HLabel,
		JLabel,
		stoptime,
		slider,
	))

	w.Resize(fyne.NewSize(300, 150))
	s := hook.Start()
	defer hook.End()
	go hook.Process(s)
	w.ShowAndRun()
}
