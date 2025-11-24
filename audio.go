package main

import (
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/blang/mpv"
)

type MPV struct {
	launchCmd *exec.Cmd
	client    *mpv.Client
}

func playAudio(videoId string) {
	audioURL := "https://www.youtube.com/watch?v=" + videoId

	oto.mpv.client.Loadfile(audioURL, mpv.LoadFileModeReplace)
}

func togglePause(shouldPause bool) {
	oto.mpv.client.SetPause(shouldPause)
}

func createMpvClient() *MPV {
	launchMpvCmd := exec.Command("mpv", "--input-ipc-server=/tmp/mpvsocket", "--idle")
	launchMpvCmd.Stdout = nil
	launchMpvCmd.Stderr = nil

	if err := launchMpvCmd.Start(); err != nil {
		log.Panicln("Error launching mpv: ", err)
	}

	socketPath := "/tmp/mpv-socket"

	for i := 0; i < 50; i++ {
		if _, err := os.Stat(socketPath); err == nil {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	ipcc := mpv.NewIPCClient("/tmp/mpvsocket")
	c := mpv.NewClient(ipcc)

	if err := c.SetProperty("vo", "null"); err != nil {
		log.Fatalf("failed to set vo: %v", err)
	}

	return &MPV{
		client:    c,
		launchCmd: launchMpvCmd,
	}
}

func turnOffMpv() {

}
