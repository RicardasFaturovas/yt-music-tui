package main

import (
	"log"
	"os"
	"path"
	"ricardasfaturovas/oto-tui/config"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	setupLog()
	config.Load()

	tview.Styles.PrimitiveBackgroundColor = tcell.ColorNames["none"]
	app := tview.NewApplication()
	mpv := NewMPV()
	oto = buildLayout(app, mpv)

	root := tview.NewFlex().
		SetDirection(0).
		AddItem(oto.pages, 0, 10, true).
		AddItem(oto.progressBar.container, 0, 1, false)

	defer func() {
		if mpv.launchCmd.Process != nil {
			mpv.launchCmd.Process.Kill()
		}
	}()

	if err := app.SetRoot(root, true).EnableMouse(true).EnablePaste(true).Run(); err != nil {
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
