package main

import (
	"github.com/rivo/tview"
)

type Oto struct {
	app         *tview.Application
	pages       *tview.Pages
	mpv         *MPV
	progressBar *ProgressBar
}

var oto *Oto

func createOto(app *tview.Application) *Oto {
	newOto := &Oto{
		app:         app,
		pages:       tview.NewPages(),
		mpv:         createMpvClient(),
		progressBar: buildProgressBar(),
	}

	searchLayout := buildSearchLayout()
	newOto.pages.AddAndSwitchToPage("search", searchLayout, true)

	return newOto
}
