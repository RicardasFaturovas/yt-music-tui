package main

import "github.com/rivo/tview"

type Oto struct {
	app   *tview.Application
	pages *tview.Pages
}

var oto *Oto

func createOto(app *tview.Application) *Oto {
	newOto := &Oto{
		app:   app,
		pages: tview.NewPages(),
	}

	searchLayout := buildSearchLayout()
	newOto.pages.AddAndSwitchToPage("search", searchLayout, true)

	return newOto
}
