package main

import (
	"log"
	"os/exec"

	"github.com/blang/mpv"
)

func playAudio(videoId string) {
	audioURL := "https://www.youtube.com/watch?v=" + videoId

	oto.mpv.Loadfile(audioURL, mpv.LoadFileModeReplace)
}

func togglePause(shouldPause bool) {
	oto.mpv.SetPause(shouldPause)
}

func createMpvClient() *mpv.Client {
	launchMpvCmd := exec.Command("mpv", "--input-ipc-server=/tmp/mpvsocket", "--idle")
	launchMpvCmd.Stdout = nil
	launchMpvCmd.Stderr = nil

	if err := launchMpvCmd.Start(); err != nil {
		log.Panicln("Error launching mpv: ", err)
	}

	ipcc := mpv.NewIPCClient("/tmp/mpvsocket")
	c := mpv.NewClient(ipcc)

	if err := c.SetProperty("vo", "null"); err != nil {
		log.Fatalf("failed to set vo: %v", err)
	}
	return c
}
