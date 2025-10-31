package main

import (
	"log"
	"os/exec"
)

func playAudio(videoId string) {
	audioURL := "https://www.youtube.com/watch?v=" + videoId
	cmd := exec.Command("mpv", "--no-video", "--quiet", audioURL)
	cmd.Stdout = nil
	cmd.Stderr = nil

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	log.Println("Playing audio...")
	cmd.Wait()
}
