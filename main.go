package main

import (
	"log"
	"os"
	"os/signal"
	"path"
	"ricardasfaturovas/oto-tui/config"
	"syscall"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	setupLog()
	config.Load()

	tview.Styles.PrimitiveBackgroundColor = tcell.ColorNames["none"]
	app := tview.NewApplication()
	oto = createOto(app)

	root := tview.NewFlex().
		SetDirection(0).
		AddItem(oto.pages, 0, 10, false).
		AddItem(oto.progressBar.container, 0, 1, false)
	defer func() {
		if oto.mpv.launchCmd.Process != nil {
			oto.mpv.launchCmd.Process.Kill()
		}
	}()

	// Catch CTRL+C or kill signals
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigc
		if oto.mpv.launchCmd.Process != nil {
			oto.mpv.launchCmd.Process.Kill()
		}
		os.Exit(0)
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
