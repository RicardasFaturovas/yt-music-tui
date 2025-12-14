package internal

import (
	_ "embed"
	"io"
	"log"
	"net"
	"os/exec"
	"strings"
	"time"

	"github.com/blang/mpv"
)

type MPV struct {
	LaunchCmd *exec.Cmd
	Client    *mpv.Client
	Stdout    io.Reader
}

//go:embed lavfi/bars.lavfi
var fileByte []byte

func NewMPV() *MPV {

	filterGraph := strings.ReplaceAll(string(fileByte), "\n", "")

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

	return &MPV{
		Client:    c,
		LaunchCmd: launchMpvCmd,
		Stdout:    stdout,
	}
}

func (m *MPV) PlaySong(videoId string) {
	audioURL := "https://www.youtube.com/watch?v=" + videoId

	m.Client.Loadfile(audioURL, mpv.LoadFileModeReplace)
	m.Client.SetPause(false)
}

func (m *MPV) TogglePause() {
	p, _ := m.Client.Pause()
	m.Client.SetPause(!p)
}

func (m *MPV) GetCurrentSong() (string, error) {
	return m.Client.Filename()
}
