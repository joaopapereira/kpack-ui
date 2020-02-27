package gui

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

func ShowSpinner(w fyne.Window, title string) {
	progressBar := widget.NewProgressBarInfinite()

	w.SetContent(widget.NewVBox(
		widget.NewLabel(title), progressBar))
}
