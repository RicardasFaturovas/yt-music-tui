package main

import (
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func getNextFocus(focusableElements []tview.Primitive) *tview.Primitive {
	var nextFocusEl tview.Primitive
	for i, el := range focusableElements {
		if el.HasFocus() {
			if i == len(focusableElements)-1 {
				nextFocusEl = focusableElements[0]
			} else {
				nextFocusEl = focusableElements[i+1]
			}
		}

	}
	if nextFocusEl == nil {
		log.Panicln("Could not find next focus element")
	}
	return &nextFocusEl
}

func getPreviousFocus(focusableElements []tview.Primitive) *tview.Primitive {
	var previousFocusEl tview.Primitive
	for i, el := range focusableElements {
		if el.HasFocus() {
			if i == 0 {
				previousFocusEl = focusableElements[len(focusableElements)-1]
			} else {
				previousFocusEl = focusableElements[i-1]
			}
		}

	}
	if previousFocusEl == nil {
		log.Panicln("Could not find previous focus element")
	}
	return &previousFocusEl
}

func focusInputCaptureCallback(event *tcell.EventKey, focusableElements []tview.Primitive) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyCtrlN:
		nextFocusEl := getNextFocus(focusableElements)
		oto.app.SetFocus(*nextFocusEl)
	case tcell.KeyCtrlP:
		previousFocusEl := getPreviousFocus(focusableElements)
		oto.app.SetFocus(*previousFocusEl)
	}
	return event
}
