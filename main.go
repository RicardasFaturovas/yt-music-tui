package main

import (
	"log"
	"os"
	"path"
	"ricardasfaturovas/y-tui/config"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	setupLog()
	config.Load()
	app := tview.NewApplication()

	root := tview.NewTreeNode("results").
		SetColor(tcell.ColorRed)
	tree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)

	tree.SetBorder(true)
	tree.SetTitle("Search results")

	searhTerms := tview.NewInputField()
	searhTerms.SetLabel("Search terms: ")
	searhTerms.SetFieldBackgroundColor(tcell.ColorNames["none"])
	searhTerms.SetBorder(true)

	searchButton := tview.NewButton("Search")
	searchButton.SetStyle(tcell.StyleDefault.Background(tcell.ColorNames["none"]))
	searchButton.SetBackgroundColorActivated(tcell.ColorNames["none"])
	searchButton.SetBorder(true)
	searchButton.SetSelectedFunc(func() { searchYoutube(searhTerms, root) })

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
		AddItem(tree, 0, 4, false)

	if err := app.SetRoot(flex, true).EnableMouse(true).EnablePaste(true).Run(); err != nil {
		panic(err)
	}

}

func searchYoutube(searchTerms *tview.InputField, resultList *tview.TreeNode) {
	searchValue := searchTerms.GetText()
	results := getSearchResults(searchValue)
	resultList.ClearChildren()

	for _, v := range results {
		node := tview.NewTreeNode(v.Title)
		node.SetSelectable(true)
		node.SetSelectedFunc(func() { playAudio(v.VideoId) })
		resultList.AddChild(node)
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
