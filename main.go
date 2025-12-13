package main

import (
	"log"
	"os"
	"os/signal"
	"path"
	"ricardasfaturovas/oto-tui/internal"
	"ricardasfaturovas/oto-tui/internal/config"
	"ricardasfaturovas/oto-tui/internal/ui"
	"syscall"

	"github.com/rivo/tview"
)

func main() {
	setupLog()

	app := tview.NewApplication()
	mpv := internal.NewMPV()
	config := config.NewConfig()
	if err := config.Load(); err != nil {
		log.Println(err)
	}

	ytClient := internal.NewYoutubeClient(config.InvidiousUrl)

	oto := ui.NewOto(app, mpv, config, ytClient)

	defer func() {
		if mpv.LaunchCmd.Process != nil {
			mpv.LaunchCmd.Process.Kill()
		}
	}()

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigc
		if mpv.LaunchCmd.Process != nil {
			mpv.LaunchCmd.Process.Kill()
		}
		os.Exit(0)
	}()

	if err := app.SetRoot(oto.Root, true).EnableMouse(true).EnablePaste(true).Run(); err != nil {
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
