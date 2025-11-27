package main

import (
	"ricardasfaturovas/oto-tui/config"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Oto struct {
	root *tview.Flex
}

func NewOto(app *tview.Application, mpv *MPV, config *config.Config, ytClient *YoutubeClient) *Oto {
	tview.Styles.PrimitiveBackgroundColor = tcell.ColorNames["none"]

	progressBar := NewProgressBar(mpv, app.QueueUpdateDraw)
	searchLayout := NewSearchLayout(mpv, progressBar.TrackProgressBar, app, config, ytClient)
	pages := tview.NewPages()
	root := tview.NewFlex().SetDirection(0)

	newOto := &Oto{
		root: root,
	}

	root.
		AddItem(pages, 0, 10, true).
		AddItem(progressBar.container, 0, 1, false)
	pages.AddAndSwitchToPage("search", searchLayout.container, true)
	return newOto
}
