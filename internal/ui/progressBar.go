package ui

import (
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"ricardasfaturovas/oto-tui/internal"

	"github.com/rivo/tview"
)

type ProgressBar struct {
	mpv               *internal.MPV
	currentSong       *tview.TextView
	bar               *tview.TextView
	container         *tview.Flex
	updateDrawHandler func(func()) *tview.Application
}

func NewProgressBar(mpv *internal.MPV, updateDrawHandler func(func()) *tview.Application) *ProgressBar {
	currentSong := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText("Nothing playing")

	bar := tview.NewTextView().
		SetTextAlign(tview.AlignCenter)
	container := tview.NewFlex().
		SetDirection(0).
		AddItem(currentSong, 1, 1, false).
		AddItem(bar, 1, 1, false).
		AddItem(tview.NewBox(), 1, 1, false)
	_, _, width, _ := container.GetRect()
	fill := strings.Repeat("━", width*4)
	container.SetBorder(true)
	container.SetTitle("Currently playing")

	bar.
		SetText("0:00 |" + fill + "| 0:00")

	progressBar := ProgressBar{
		currentSong:       currentSong,
		bar:               bar,
		container:         container,
		mpv:               mpv,
		updateDrawHandler: updateDrawHandler,
	}
	return &progressBar
}

func (p *ProgressBar) TrackProgressBar(songName string) {
	p.currentSong.SetText(songName)

	for {
		playlistCount, _ := p.mpv.Client.GetFloatProperty("playlist-playing-pos")
		isPaused, pauseErr := p.mpv.Client.Pause()
		if pauseErr != nil {
			log.Println("Error getting pause")
		}
		if isPaused {
			time.Sleep(1 * time.Second)
			continue
		}

		currentProgress, positionErr := p.mpv.Client.Position()
		if positionErr != nil {
			log.Println("Error getting position")
		}

		duration, durationErr := p.mpv.Client.Duration()
		if durationErr != nil {
			log.Println("Error getting duration")
		}

		percent, err := p.mpv.Client.PercentPosition()
		if err != nil {
			log.Println("Error getting position")
		}

		_, _, width, _ := p.container.GetRect()
		elapsedTime := time.Duration(currentProgress * float64(time.Second))
		timeLeft := time.Duration((duration - currentProgress) * float64(time.Second))

		currentPosition := int(math.Round(percent / 100 * float64(width/2)))

		fill := fmtProgress(currentPosition, width/2)

		p.updateDrawHandler(func() {
			p.bar.SetText(fmt.Sprintf("%s |%s| %s", fmtDuration(elapsedTime), fill, fmtDuration(timeLeft)))
		})

		if playlistCount < 0 {
			log.Println("IDLE Turning off")
			break
		}

		<-time.After(time.Second)
	}
}

func fmtDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

func fmtProgress(currentPostion int, length int) string {
	return fmt.Sprintf("%s%s%s",
		strings.Repeat("━", currentPostion),
		"I",
		strings.Repeat("━", length-currentPostion),
	)
}
