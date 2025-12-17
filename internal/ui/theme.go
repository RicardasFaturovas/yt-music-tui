package ui

import (
	"ricardasfaturovas/oto-tui/internal/config"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func loadTheme(c *config.Config) {
	tview.Styles.PrimitiveBackgroundColor = tcell.Color(c.Theme.BackgroundColor)
	tview.Styles.PrimaryTextColor = tcell.Color(c.Theme.PrimaryTextColor)
	tview.Styles.BorderColor = tcell.Color(c.Theme.BorderColor)
	tview.Styles.TitleColor = tcell.Color(c.Theme.TitleColor)
	tview.Styles.SecondaryTextColor = tcell.Color(c.Theme.SecondaryTextColor)
}
