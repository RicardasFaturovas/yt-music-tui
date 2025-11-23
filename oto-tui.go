package main

import (
	"github.com/blang/mpv"
	"github.com/rivo/tview"
)

type Oto struct {
	app   *tview.Application
	pages *tview.Pages
	mpv   *mpv.Client
}

var oto *Oto

func createOto(app *tview.Application) *Oto {
	newOto := &Oto{
		app:   app,
		pages: tview.NewPages(),
		mpv:   createMpvClient(),
	}

	searchLayout := buildSearchLayout()
	newOto.pages.AddAndSwitchToPage("search", searchLayout, true)

	return newOto
}
