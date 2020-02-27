package gui

import (
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	"kpackui/static"
)

func ErrorContainer(err error) *fyne.Container {
	return fyne.NewContainerWithLayout(
		layout.NewGridLayoutWithColumns(1),
		widget.NewLabelWithStyle("Error: Retrieving contexts", fyne.TextAlignLeading, fyne.TextStyle{
			Bold: true,
		}),
		fyne.NewContainerWithLayout(
			layout.NewGridLayoutWithColumns(2),
			widget.NewIcon(static.SadEmojiIcon()),
			widget.NewLabel(err.Error()),
		),
	)
}
