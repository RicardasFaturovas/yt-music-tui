package main

import (
	"log"
	"os"
	"path"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	setupLog()
	loadConfig()
	app := tview.NewApplication()

	list := tview.NewList()
	list.SetBorder(true).SetTitle("Music results")

	searhTerms := tview.NewInputField()
	searhTerms.SetLabel("Search terms: ")
	searhTerms.SetFieldBackgroundColor(tcell.ColorNames["none"])
	searhTerms.SetBorder(true)

	searchButton := tview.NewButton("Search")
	searchButton.SetStyle(tcell.StyleDefault.Background(tcell.ColorNames["none"]))
	searchButton.SetBackgroundColorActivated(tcell.ColorNames["none"])
	searchButton.SetBorder(true)
	searchButton.SetSelectedFunc(func() { searchYoutube(searhTerms, list) })

	searchRow := tview.NewFlex().
		SetDirection(1).
		AddItem(searhTerms, 0, 5, true).
		AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorNames["none"]), 0, 1, false).
		AddItem(searchButton, 0, 2, false)
	searchRow.SetBorder(true)
	searchRow.SetBorderPadding(2, 2, 1, 1)

	flex := tview.NewFlex().
		SetDirection(0).
		AddItem(searchRow, 0, 1, true).
		AddItem(list, 0, 4, false)

	if err := app.SetRoot(flex, true).EnableMouse(true).EnablePaste(true).Run(); err != nil {
		panic(err)
	}

}

func searchYoutube(searchTerms *tview.InputField, resultList *tview.List) {
	log.Printf("SEARCHING YT")
	searchValue := searchTerms.GetText()
	results := getSearchResults(searchValue)

	for _, v := range results {
		log.Println("PLS")
		log.Println(v.Title)
		resultList.AddItem(v.Title, "", 'a', func() {
			playAudio(v.VideoId)
		})
	}
}

func setupLog() {
	tmpDir, _ := os.Getwd()
	logFile := path.Join(tmpDir, "go-music.log")
	file, e := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if e != nil {
		log.Fatalf("Error opening file %s", logFile)
	}

	log.SetOutput(file)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
