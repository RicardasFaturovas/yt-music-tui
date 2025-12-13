package ui

import (
	"encoding/json"
	"log"
	"os"
	"path"
	"path/filepath"
	"ricardasfaturovas/oto-tui/internal"
	"ricardasfaturovas/oto-tui/internal/config"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type SearchLayout struct {
	mpv                *internal.MPV
	progressBarHandler func(songName string)
	app                *tview.Application
	config             *config.Config
	youtubeClient      *internal.YoutubeClient
	container          *tview.Pages
}

func NewSearchLayout(
	mpv *internal.MPV,
	progressBarHandler func(songName string),
	app *tview.Application,
	config *config.Config,
	youtubeClient *internal.YoutubeClient,
) *SearchLayout {
	base := tview.NewPages()

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

	albumCarousel := NewAlbumCarousel(app)
	visualizer := NewAudioVisualizer(mpv.Stdout, app)
	go func() {
		visualizer.visualize()
	}()

	searchResults.SetChangedFunc(func(node *tview.TreeNode) {
		if node.GetLevel() == 1 {
			go albumCarousel.CycleRight()
		}
	})

	visual := tview.NewFlex().
		SetDirection(0).
		AddItem(albumCarousel.container, 0, 1, false).
		AddItem(visualizer.container, 0, 1, false)

	visual.SetBorder(true)
	visual.SetBorderPadding(2, 2, 4, 0)

	middleRow := tview.NewFlex().
		SetDirection(1).
		AddItem(searchResults, 0, 1, false).
		AddItem(visual, 0, 1, false)

	flex := tview.NewFlex().
		SetDirection(0).
		AddItem(searchRow, 0, 2, true).
		AddItem(middleRow, 0, 7, false)

	base.AddPage("main", flex, true, true)

	layout := &SearchLayout{
		mpv:                mpv,
		container:          base,
		progressBarHandler: progressBarHandler,
		app:                app,
		config:             config,
		youtubeClient:      youtubeClient,
	}

	focusableElements := []tview.Primitive{searchTerms, searchButton, searchResults}

	searchButton.SetSelectedFunc(func() {
		searchText := searchTerms.GetText()
		resultsTree := layout.buildSearchResultTree(searchText, searchResults)
		searchResults.SetRoot(resultsTree)
		searchResults.SetCurrentNode(resultsTree)
	})

	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		return focusInputCaptureCallback(event, focusableElements, layout.app.SetFocus)
	})

	return layout
}

func (s *SearchLayout) buildSearchResultTree(searchText string, treeView *tview.TreeView) *tview.TreeNode {
	results, err := s.youtubeClient.GetSearchResults(searchText)
	if err != nil {
		s.newErrorPopup("Error", "Error fetching music from youtube", 5*time.Second)
	}

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
			reference := node.GetReference().(internal.YoutubeSong)
			if node.GetText() == "Play" {
				s.playSongHandler(reference)
			} else if node.GetText() == "Add to playlist" {
				s.loadPlaylistsHandler(node, reference)
			}
		} else {
			if node.GetText() == "Create new playlist" {
				reference := node.GetReference().(internal.YoutubeSong)
				s.createPlaylistHandler(reference)
			} else {
				reference := node.GetReference().(internal.YoutubeSong)
				s.addToPlaylist(reference, node.GetText())
			}
		}
	})
}

func (s *SearchLayout) playSongHandler(song internal.YoutubeSong) {
	currentSong, _ := s.mpv.GetCurrentSong()
	log.Println(song.Title, song.VideoId, currentSong)
	if strings.Contains(currentSong, song.VideoId) {
		s.mpv.TogglePause()
	} else {
		s.mpv.PlaySong(song.VideoId)
		go s.progressBarHandler(song.Title)
	}
}

func (s *SearchLayout) loadPlaylistsHandler(target *tview.TreeNode, song internal.YoutubeSong) {
	// TODO: Need to figure out better directory management
	home, _ := os.UserHomeDir()
	playlistPath := home + "/" + s.config.PlaylistPath

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

func (s *SearchLayout) createPlaylistHandler(song internal.YoutubeSong) {
	playlistNameInput := s.newPlaylistPopup()

	playlistNameInput.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:
			playListName := playlistNameInput.GetText()
			s.addToPlaylist(song, playListName+".playlist")
			s.container.RemovePage("new-playlist")
		case tcell.KeyEsc:
			s.container.RemovePage("new-playlist")
		}
	})
}

func (s *SearchLayout) addToPlaylist(song internal.YoutubeSong, playlistName string) {
	home, _ := os.UserHomeDir()
	fullPath := path.Join(home, s.config.PlaylistPath, playlistName)

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
	s.container.AddPage("new-playlist", popup, true, true)
	s.app.SetFocus(inputField)

	return inputField
}

func (s *SearchLayout) newErrorPopup(title string, description string, timeout time.Duration) {
	popupId := "error"
	textView := tview.NewTextView().
		SetText(description)
	textView.SetTextAlign(tview.AlignCenter).SetTitle(title).SetBorder(true)

	popup := topRight(textView, 60, 5)
	s.container.AddPage(popupId, popup, true, true)

	go func() {
		time.Sleep(timeout)
		s.app.QueueUpdateDraw(func() {
			s.container.RemovePage(popupId)
		})
	}()

}
