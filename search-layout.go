package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var focusableElements []tview.Primitive

func buildSearchLayout() *tview.Flex {
	root := tview.NewTreeNode(".").
		SetColor(tcell.ColorNames["none"])
	searchResults := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)

	searchResults.SetBorder(true)
	searchResults.SetTitle("Search results")

	searchTerms := tview.NewInputField()
	searchTerms.SetLabel("Search terms: ")
	searchTerms.SetFieldBackgroundColor(tcell.ColorNames["none"])
	searchTerms.SetBorder(true)

	searchButton := tview.NewButton("Search")
	searchButton.SetStyle(tcell.StyleDefault.Background(tcell.ColorNames["none"]))
	searchButton.SetBackgroundColorActivated(tcell.ColorNames["none"])
	searchButton.SetBorder(true)
	searchButton.SetSelectedFunc(func() { searchYoutube(searchTerms, root) })

	searchRow := tview.NewFlex().
		SetDirection(1).
		AddItem(searchTerms, 0, 5, true).
		AddItem(tview.NewBox(), 0, 1, false).
		AddItem(searchButton, 0, 2, false)
	searchRow.SetBorder(true)
	searchRow.SetBorderPadding(2, 2, 1, 1)
	searchRow.SetTitle("Search youtube")

	flex := tview.NewFlex().
		SetDirection(0).
		AddItem(searchRow, 0, 1, true).
		AddItem(searchResults, 0, 4, false)

	addFocusableElements(searchTerms, searchButton, searchResults)
	return flex
}

func addFocusableElements(elements ...tview.Primitive) {
	for _, el := range elements {
		focusableElements = append(focusableElements, el)
	}
}

func getNextFocus() *tview.Primitive {
	var nextFocusEl tview.Primitive
	for i, el := range focusableElements {
		if el.HasFocus() {
			if i == len(focusableElements)-1 {
				nextFocusEl = focusableElements[0]
			} else {
				nextFocusEl = focusableElements[i+1]
			}
		}

	}
	if nextFocusEl == nil {
		log.Panicln("Could not find next focus element")
	}
	return &nextFocusEl
}

func getPreviousFocus() *tview.Primitive {
	var previousFocusEl tview.Primitive
	for i, el := range focusableElements {
		if el.HasFocus() {
			if i == 0 {
				previousFocusEl = focusableElements[len(focusableElements)-1]
			} else {
				previousFocusEl = focusableElements[i-1]
			}
		}

	}
	if previousFocusEl == nil {
		log.Panicln("Could not find previous focus element")
	}
	return &previousFocusEl
}

func searchYoutube(searchTerms *tview.InputField, resultList *tview.TreeNode) {
	searchValue := searchTerms.GetText()
	results := getSearchResults(searchValue)
	resultList.ClearChildren()

	for _, v := range results {
		if v.Title != "" {
			songNode := tview.NewTreeNode(v.Title)
			songNode.SetExpanded(false)

			songNode.SetSelectedFunc(func() {
				resultList.CollapseAll()
				resultList.Expand()
				songNode.SetExpanded(!songNode.IsExpanded())
			})

			playNode := tview.NewTreeNode("Play")
			playNode.SetSelectable(true)
			playNode.SetSelectedFunc(func() { playAudio(v.VideoId) })

			playlistNode := tview.NewTreeNode("Add to playlist")
			playlistNode.SetSelectedFunc(func() { addToPlaylist(v) })

			songNode.AddChild(playNode)
			songNode.AddChild(playlistNode)
			resultList.AddChild(songNode)
		}
	}
}

func addToPlaylist(song YoutubeVideo) {
	f, err := os.Create("playlist.json")
	if err != nil {
		log.Panicln(err)
	}
	defer f.Close()
	as_json, err := json.MarshalIndent(song, "", "\t")
	if err != nil {
		log.Panicln(err)
	}
	f.Write(as_json)
}
