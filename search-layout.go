package main

import (
	"encoding/json"
	"log"
	"os"
	"path"
	"path/filepath"
	"ricardasfaturovas/oto-tui/config"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type SearchLayout struct {
	mpv                *MPV
	container          *tview.Flex
	progressBarHandler func(songName string)
}

func NewSearchLayout(mpv *MPV, progressBarHandler func(songName string)) *SearchLayout {
	searchTerms := tview.NewInputField()
	searchTerms.SetLabel("Search terms: ")
	searchTerms.SetFieldBackgroundColor(tcell.ColorNames["none"])
	searchTerms.SetBorder(true)

	searchButton := tview.NewButton("Search")
	searchButton.SetStyle(tcell.StyleDefault.Background(tcell.ColorNames["none"]))
	searchButton.SetBackgroundColorActivated(tcell.ColorNames["none"])
	searchButton.SetBorder(true)

	searchRow := tview.NewFlex().
		SetDirection(1).
		AddItem(searchTerms, 0, 5, true).
		AddItem(tview.NewBox(), 0, 1, false).
		AddItem(searchButton, 0, 2, false)
	searchRow.SetBorder(true)
	searchRow.SetBorderPadding(2, 2, 1, 1)
	searchRow.SetTitle("Search youtube")

	searchResults := tview.NewTreeView()
	searchResults.SetBorder(true)
	searchResults.SetTitle("Search results")

	flex := tview.NewFlex().
		SetDirection(0).
		AddItem(searchRow, 0, 2, true).
		AddItem(searchResults, 0, 7, false)

	layout := &SearchLayout{
		mpv:                mpv,
		container:          flex,
		progressBarHandler: progressBarHandler,
	}

	focusableElements := []tview.Primitive{searchTerms, searchButton, searchResults}

	searchButton.SetSelectedFunc(func() {
		searchText := searchTerms.GetText()
		resultsTree := layout.buildSearchResultTree(searchText, searchResults)
		searchResults.SetRoot(resultsTree)
		searchResults.SetCurrentNode(resultsTree)
	})

	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		return focusInputCaptureCallback(event, focusableElements)
	})

	return layout
}

func (s *SearchLayout) buildSearchResultTree(searchText string, treeView *tview.TreeView) *tview.TreeNode {
	results := getSearchResults(searchText)

	root := tview.NewTreeNode(".").
		SetColor(tcell.ColorNames["none"])

	for _, v := range results {
		if v.Title != "" {
			playNode := tview.NewTreeNode("Play")
			playNode.SetSelectable(true)
			playNode.SetReference(v)

			playlistNode := tview.NewTreeNode("Add to playlist")
			playlistNode.SetReference(v)

			songNode := tview.NewTreeNode(v.Title)
			songNode.SetExpanded(false)
			songNode.AddChild(playNode)
			songNode.AddChild(playlistNode)

			root.AddChild(songNode)
		}
	}

	s.attachResultTreeHandlers(treeView)

	return root
}

func (s *SearchLayout) attachResultTreeHandlers(treeView *tview.TreeView) {
	treeView.SetSelectedFunc(func(node *tview.TreeNode) {
		if node.GetLevel() == 1 {
			treeView.GetRoot().CollapseAll()
			treeView.GetRoot().Expand()
			node.SetExpanded(!node.IsExpanded())
		} else if node.GetLevel() == 2 {
			reference := node.GetReference().(YoutubeVideo)
			if node.GetText() == "Play" {
				s.playSongHandler(reference)
			} else if node.GetText() == "Add to playlist" {
				s.loadPlaylistsHandler(node, reference)
			}
		} else {
			if node.GetText() == "Create new playlist" {
				reference := node.GetReference().(YoutubeVideo)
				s.createPlaylistHandler(reference)
			} else {
				reference := node.GetReference().(YoutubeVideo)
				s.addToPlaylist(reference, node.GetText())
			}
		}
	})
}

func (s *SearchLayout) playSongHandler(song YoutubeVideo) {
	currentSong, _ := s.mpv.GetCurrentSong()
	isPaused := false
	if strings.Contains(currentSong, song.VideoId) {
		isPaused = !isPaused
		s.mpv.TogglePause(isPaused)
	} else {
		s.mpv.PlaySong(song.VideoId)
		go s.progressBarHandler(song.Title)
	}
}

func (s *SearchLayout) loadPlaylistsHandler(target *tview.TreeNode, song YoutubeVideo) {
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
			node.SetReference(song)
			target.AddChild(node)
		}
	}

	createPlaylistNode := tview.NewTreeNode("Create new playlist")
	createPlaylistNode.SetReference(song)
	createPlaylistNode.SetSelectable(true)

	target.AddChild(createPlaylistNode)
	target.SetExpanded(!target.IsExpanded())
}

func (s *SearchLayout) createPlaylistHandler(song YoutubeVideo) {
	playlistNameInput := s.newPlaylistPopup()

	playlistNameInput.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:
			playListName := playlistNameInput.GetText()
			s.addToPlaylist(song, playListName+".playlist")
			oto.pages.RemovePage("new-playlist")
		case tcell.KeyEsc:
			oto.pages.RemovePage("new-playlist")
		}

	})
}

func (s *SearchLayout) addToPlaylist(song YoutubeVideo, playlistName string) {
	home, _ := os.UserHomeDir()
	fullPath := path.Join(home, config.Get().PlaylistPath, playlistName)

	f, err := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Panicln(err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	if err := enc.Encode(song); err != nil {
		log.Panicln(err)
	}
}

func (s *SearchLayout) newPlaylistPopup() *tview.InputField {
	inputField := tview.NewInputField().
		SetLabel("New playlist").
		SetFieldWidth(0).
		SetAcceptanceFunc(tview.InputFieldMaxLength(50))

	inputField.SetTitle("New playlist").
		SetBorder(true).
		SetBorderPadding(1, 0, 2, 2)
	popup := center(inputField, 60, 5)
	oto.pages.AddPage("new-playlist", popup, true, true)
	oto.app.SetFocus(inputField)

	return inputField
}
