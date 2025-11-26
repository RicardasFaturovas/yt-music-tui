package main

import (
	"github.com/rivo/tview"
	"ricardasfaturovas/oto-tui/config"
)

type Oto struct {
	app         *tview.Application
	pages       *tview.Pages
	progressBar *ProgressBar
	root        *tview.Flex
}

var oto *Oto

func buildLayout(app *tview.Application, mpv *MPV) *Oto {
	config := config.NewConfig()
	ytClient := NewYoutubeClient(config.InvidiousUrl)

	progressBar := NewProgressBar(mpv, app.QueueUpdateDraw)
	searchLayout := NewSearchLayout(mpv, progressBar.TrackProgressBar, app.SetFocus, config, ytClient)

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
