package main

import (
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/blang/mpv"
)

type MPV struct {
	launchCmd *exec.Cmd
	client    *mpv.Client
	stdout    io.Reader
}

func NewMPV() *MPV {
	lavfiData, err := os.ReadFile("waves.lavfi")
	if err != nil {
		log.Fatal(err)
	}
	filterGraph := strings.ReplaceAll(string(lavfiData), "\n", "")

	IPCPath := "/tmp/mpvsocket"
	launchMpvCmd := exec.Command(
		"mpv",
		"--no-terminal",
		"--input-ipc-server="+IPCPath,
		"--lavfi-complex="+filterGraph,
		"--vo=tct",
		"--really-quiet",
		"--idle",
	)

	stdout, _ := launchMpvCmd.StdoutPipe()
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
		stdout:    stdout,
	}
}

func (m *MPV) PlaySong(videoId string) {
	audioURL := "https://www.youtube.com/watch?v=" + videoId

	m.client.Loadfile(audioURL, mpv.LoadFileModeReplace)
	m.client.SetPause(false)
}

func (m *MPV) TogglePause() {
	p, _ := m.client.Pause()
	m.client.SetPause(!p)
}

func (m *MPV) GetCurrentSong() (string, error) {
	return m.client.Filename()
}
