package main

import (
	"log"
	"os"
	"path"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	setupLog()

	tview.Styles.PrimitiveBackgroundColor = tcell.ColorNames["none"]
	app := tview.NewApplication()
	mpv := NewMPV()
	oto = buildLayout(app, mpv)

	defer func() {
		if mpv.launchCmd.Process != nil {
			mpv.launchCmd.Process.Kill()
		}
	}()

	if err := app.SetRoot(oto.root, true).EnableMouse(true).EnablePaste(true).Run(); err != nil {
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
