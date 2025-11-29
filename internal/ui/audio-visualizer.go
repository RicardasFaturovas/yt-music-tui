package ui

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/rivo/tview"
)

const WIDTH = 80

// ASCII characters from dark to light
var asciiChars = " .:-=+*#%@"

// get brightness from background color only
var colorRegex = regexp.MustCompile(`\x1b\[48;2;(\d+);(\d+);(\d+)m`)

func rgbToLuminance(r, g, b int) float64 {
	return 0.2126*float64(r) + 0.7152*float64(g) + 0.0722*float64(b)
}

func luminanceToASCII(lum float64) string {
	index := int((lum / 255.0) * float64(len(asciiChars)-1))
	return string(asciiChars[index])
}

type AudioVisualizer struct {
	tctOutput io.Reader
	container *tview.TextView
	app       *tview.Application
}

func NewAudioVisualizer(tctOutput io.Reader, app *tview.Application) *AudioVisualizer {
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(false)

	return &AudioVisualizer{
		tctOutput: tctOutput,
		container: textView,
		app:       app,
	}
}

func (a *AudioVisualizer) visualize() {
	scanner := bufio.NewScanner(a.tctOutput)
	var count int
	var frame strings.Builder

	for scanner.Scan() {
		line := scanner.Text()
		log.Println(line)

		matches := colorRegex.FindAllStringSubmatch(line, -1)
		for _, m := range matches {
			r, _ := strconv.Atoi(m[1])
			g, _ := strconv.Atoi(m[2])
			b, _ := strconv.Atoi(m[3])

			lum := rgbToLuminance(r, g, b)
			ascii := luminanceToASCII(lum)
			frame.WriteString(ascii)

			count++
			if count == WIDTH {
				frame.WriteString("\n")
				count = 0
			}
		}

		// Update TextView when a row is complete
		if count == 0 && frame.Len() > 0 {
			output := frame.String()
			frame.Reset()
			a.app.QueueUpdateDraw(func() {
				fmt.Fprint(a.container, output)
			})
		}
	}
}
