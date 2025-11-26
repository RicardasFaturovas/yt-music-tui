package main

import (
	"log"
	"os"
	"path"
	"ricardasfaturovas/oto-tui/config"

	"github.com/rivo/tview"
)

func main() {
	setupLog()

	app := tview.NewApplication()
	mpv := NewMPV()
	config := config.NewConfig()
	ytClient := NewYoutubeClient(config.InvidiousUrl)

	oto := NewOto(app, mpv, config, ytClient)

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
