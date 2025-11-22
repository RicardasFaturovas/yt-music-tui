package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"ricardasfaturovas/y-tui/config"

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
			playlistNode.SetSelectedFunc(func() { loadPlaylists(playlistNode, v) })

			songNode.AddChild(playNode)
			songNode.AddChild(playlistNode)
			resultList.AddChild(songNode)
		}
	}
}

func loadPlaylists(target *tview.TreeNode, song YoutubeVideo) {
	// TODO: Need to figure out better directory management
	home, _ := os.UserHomeDir()
	playlistPath := home + "/" + config.Get().PlaylistPath

	target.ClearChildren()
	files, err := os.ReadDir(playlistPath)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".playlist" {
			node := tview.NewTreeNode(file.Name()).
				SetSelectable(true)

			node.SetSelectedFunc(func() { addToPlaylist(song, filepath.Join(playlistPath, file.Name())) })
			target.AddChild(node)
		}
	}

	createPlaylistNode := tview.NewTreeNode("Create new playlist")
	createPlaylistNode.SetSelectable(true)
	createPlaylistNode.SetSelectedFunc(func() { createPlaylistPopup(song, playlistPath) })

	target.AddChild(createPlaylistNode)
	target.SetExpanded(!target.IsExpanded())
}

func createPlaylistPopup(song YoutubeVideo, playlistPath string) {
	playlistNameInput := newPlaylistPopup()
	playlistNameInput.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:
			playListName := playlistNameInput.GetText()
			addToNewPlaylist(song, playListName, playlistPath)
			oto.pages.RemovePage("new-playlist")
		case tcell.KeyEsc:
			oto.pages.RemovePage("new-playlist")
		}

	})
}

func addToNewPlaylist(song YoutubeVideo, playlistName string, playlistPath string) {
	addToPlaylist(song, filepath.Join(playlistPath, playlistName+".playlist"))
}

func addToPlaylist(song YoutubeVideo, playlistPath string) {
	f, err := os.OpenFile(playlistPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Panicln(err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	if err := enc.Encode(song); err != nil {
		log.Panicln(err)
	}
}

func newPlaylistPopup() *tview.InputField {
	inputField := tview.NewInputField().
		SetLabel("New playlist").
		SetFieldWidth(0).
		SetAcceptanceFunc(tview.InputFieldMaxLength(50))

	inputField.SetTitle("New playlist").
		SetBorder(true).
		SetBorderPadding(1, 0, 2, 2)

	oto.pages.
		AddPage("new-playlist", center(inputField, 60, 5), true, true)

	return inputField
}
