package main

import (
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"github.com/rivo/tview"
)

type ProgressBar struct {
	currentSong *tview.TextView
	bar         *tview.TextView
	container   *tview.Flex
}

func buildProgressBar() *ProgressBar {
	currentSong := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText("Test song")

	bar := tview.NewTextView().
		SetTextAlign(tview.AlignCenter)
	container := tview.NewFlex().
		SetDirection(0).
		AddItem(currentSong, 1, 1, false).
		AddItem(bar, 1, 1, false)
	_, _, width, _ := container.GetRect()
	fill := strings.Repeat("━", width*4)
	container.SetBorder(true)
	container.SetTitle("Currently playing")

	bar.
		SetText("0:00 |" + fill + "| 4:20")

	progressBar := ProgressBar{
		currentSong: currentSong,
		bar:         bar,
		container:   container,
	}
	return &progressBar
}

func updateBar(songName string) {
	oto.progressBar.currentSong.SetText(songName)

	for {
		playlistCount, _ := oto.mpv.client.GetFloatProperty("playlist-playing-pos")
		isPaused, pauseErr := oto.mpv.client.Pause()
		if pauseErr != nil {
			log.Println("Error getting pause/idle")
		}
		if isPaused {
			time.Sleep(1 * time.Second)
			continue
		}

		currentProgress, positionErr := oto.mpv.client.Position()
		if positionErr != nil {
			log.Println("Error getting position")
		}

		duration, durationErr := oto.mpv.client.Duration()
		if durationErr != nil {
			log.Println("Error getting duration")
		}

		percent, err := oto.mpv.client.PercentPosition()
		if err != nil {
			log.Println("Error getting position")
		}

		_, _, width, _ := oto.progressBar.container.GetRect()
		elapsedTime := time.Duration(currentProgress * float64(time.Second))
		timeLeft := time.Duration((duration - currentProgress) * float64(time.Second))

		currentPosition := int(math.Round(percent / 100 * float64(width/2)))

		fill := fmtProgress(currentPosition, width/2)

		oto.app.QueueUpdateDraw(func() {
			oto.progressBar.bar.SetText(fmt.Sprintf("%s |%s| %s", fmtDuration(elapsedTime), fill, fmtDuration(timeLeft)))
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
