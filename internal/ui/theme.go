package ui

import (
	"ricardasfaturovas/oto-tui/internal/config"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Theme struct {
	backgroundColor    tcell.Color
	primaryTextColor   tcell.Color
	secondaryTextColor tcell.Color
	borderColor        tcell.Color
	titleColor         tcell.Color
}

var theme = Theme{
	backgroundColor:    tcell.NewRGBColor(42, 39, 63),
	primaryTextColor:   tcell.NewRGBColor(224, 222, 244),
	borderColor:        tcell.NewRGBColor(62, 143, 176),
	titleColor:         tcell.NewRGBColor(234, 154, 151),
	secondaryTextColor: tcell.NewRGBColor(156, 207, 216),
}

func loadTheme(c *config.Config) {
	tview.Styles.PrimitiveBackgroundColor = theme.backgroundColor
	tview.Styles.PrimaryTextColor = theme.primaryTextColor
	tview.Styles.BorderColor = theme.borderColor
	tview.Styles.TitleColor = theme.titleColor
	tview.Styles.SecondaryTextColor = theme.secondaryTextColor
}
