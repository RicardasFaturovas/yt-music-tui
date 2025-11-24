package main

import (
	"fmt"
	"log"
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

func updateBar() {
	for {
		_, idleErr := oto.mpv.client.Idle()
		isPaused, pauseErr := oto.mpv.client.Pause()

		if idleErr != nil || pauseErr != nil {
			log.Panicln("Error getting mpv idle/paused status")
		}

		// if isIdle {
		// 	break
		// }
		//
		if isPaused {
			time.Sleep(1 * time.Second)
			continue
		}
		currentProgress, positionErr := oto.mpv.client.Position()
		if positionErr != nil {
			log.Println("Error getting position")
		}

		log.Println(currentProgress)

		oto.app.QueueUpdateDraw(func() {
			duration := time.Duration(currentProgress * float64(time.Second))
			oto.progressBar.bar.SetText(fmt.Sprintf(fmtDuration(duration)))
		})

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
