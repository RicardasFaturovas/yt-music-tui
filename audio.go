package main

import (
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/blang/mpv"
)

type MPV struct {
	launchCmd *exec.Cmd
	client    *mpv.Client
}

func NewMPV() *MPV {
	IPCPath := "/tmp/mpvsocket"
	launchMpvCmd := exec.Command("mpv", "--input-ipc-server="+IPCPath, "--idle")
	launchMpvCmd.Stdout = nil
	launchMpvCmd.Stderr = nil

	if err := launchMpvCmd.Start(); err != nil {
		log.Panicln("Error launching mpv: ", err)
	}

	for range 50 {
		if _, err := net.Dial("unix", IPCPath); err == nil {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	ipcc := mpv.NewIPCClient(IPCPath)
	c := mpv.NewClient(ipcc)

	if err := c.SetProperty("vo", "null"); err != nil {
		log.Fatalf("failed to set vo: %v", err)
	}

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigc
		if launchMpvCmd.Process != nil {
			launchMpvCmd.Process.Kill()
		}
		os.Exit(0)
	}()

	return &MPV{
		client:    c,
		launchCmd: launchMpvCmd,
	}
}

func (m *MPV) PlaySong(videoId string) {
	audioURL := "https://www.youtube.com/watch?v=" + videoId

	m.client.Loadfile(audioURL, mpv.LoadFileModeReplace)
}

func (m *MPV) TogglePause(shouldPause bool) {
	m.client.SetPause(shouldPause)
}

func (m *MPV) GetCurrentSong() (string, error) {
	return m.client.Filename()
}
