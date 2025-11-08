package main

import (
	"log"
	"os"
	"path"
	"ricardasfaturovas/y-tui/config"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	setupLog()
	config.Load()
	app := tview.NewApplication()
	tview.Styles.PrimitiveBackgroundColor = tcell.ColorNames["none"]
	searchLayout := buildSearchLayout()

	searchLayout.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlN:
			nextFocusEl := getNextFocus()
			app.SetFocus(*nextFocusEl)
		case tcell.KeyCtrlP:
			previousFocusEl := getPreviousFocus()
			app.SetFocus(*previousFocusEl)
		}
		return event
	})
	if err := app.SetRoot(searchLayout, true).EnableMouse(true).EnablePaste(true).Run(); err != nil {
		panic(err)
	}
}

func setupLog() {
	tmpDir, _ := os.Getwd()
	logFile := path.Join(tmpDir, "go-music.log")
	file, e := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if e != nil {
		log.Fatalf("Error opening file %s", logFile)
	}

	log.SetOutput(file)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
