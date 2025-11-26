package main

import (
	"github.com/rivo/tview"
)

type Oto struct {
	app         *tview.Application
	pages       *tview.Pages
	progressBar *ProgressBar
}

var oto *Oto

func buildLayout(app *tview.Application, mpv *MPV) *Oto {
	progressBar := NewProgressBar(mpv)
	searchLayout := NewSearchLayout(mpv, progressBar.TrackProgressBar)

	newOto := &Oto{
		app:         app,
		pages:       tview.NewPages(),
		progressBar: progressBar,
	}

	newOto.pages.AddAndSwitchToPage("search", searchLayout.container, true)
	return newOto
}
