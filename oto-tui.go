package main

import (
	"github.com/rivo/tview"
)

type Oto struct {
	app         *tview.Application
	pages       *tview.Pages
	progressBar *ProgressBar
	root        *tview.Flex
}

var oto *Oto

func buildLayout(app *tview.Application, mpv *MPV) *Oto {
	progressBar := NewProgressBar(mpv, app.QueueUpdateDraw)
	searchLayout := NewSearchLayout(mpv, progressBar.TrackProgressBar, app.SetFocus)

	root := tview.NewFlex().SetDirection(0)

	newOto := &Oto{
		app:         app,
		pages:       tview.NewPages(),
		progressBar: progressBar,
		root:        root,
	}

	root.
		AddItem(newOto.pages, 0, 10, true).
		AddItem(newOto.progressBar.container, 0, 1, false)
	newOto.pages.AddAndSwitchToPage("search", searchLayout.container, true)
	return newOto
}
