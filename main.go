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
	oto = createOto(app)

	tview.Styles.PrimitiveBackgroundColor = tcell.ColorNames["none"]

	if err := app.SetRoot(oto.pages, true).EnableMouse(true).EnablePaste(true).Run(); err != nil {
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

func center(p tview.Primitive, width, height int) tview.Primitive {
	return tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(p, height, 1, false).
			AddItem(nil, 0, 1, false), width, 1, false).
		AddItem(nil, 0, 1, false)
}
