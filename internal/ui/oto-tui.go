package ui

import (
	"ricardasfaturovas/oto-tui/internal"
	"ricardasfaturovas/oto-tui/internal/config"

	"github.com/rivo/tview"
)

type Oto struct {
	Root *tview.Flex
}

func NewOto(app *tview.Application, mpv *internal.MPV, config *config.Config, ytClient *internal.YoutubeClient) *Oto {
	loadTheme(config)

	progressBar := NewProgressBar(mpv, app.QueueUpdateDraw)
	searchLayout := NewSearchLayout(mpv, progressBar.TrackProgressBar, app, config, ytClient)
	pages := tview.NewPages()
	root := tview.NewFlex().SetDirection(0)

	newOto := &Oto{
		Root: root,
	}

	root.
		AddItem(pages, 0, 10, true).
		AddItem(progressBar.container, 0, 1, false)
	pages.AddAndSwitchToPage("search", searchLayout.container, true)
	return newOto
}
