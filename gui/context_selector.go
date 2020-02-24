package gui

import (
	"fmt"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

type ContextGetter interface {
	GetAll() ([]string, error)
}

func ContextSelector(win fyne.Window, contextGetter ContextGetter, onContextSelected func(string)) {
	grid := fyne.NewContainerWithLayout(layout.NewGridLayout(1))
	group1 := widget.NewGroup("Contexts")
	grid.AddObject(group1)

	contexts, err := contextGetter.GetAll()
	if err != nil {
		fmt.Printf("unable to load contexts: %s", err)
	}

	for _, name := range contexts {
		name := name
		group1.Append(widget.NewButton(name, func() {
			onContextSelected(name)
		}))
	}

	win.SetContent(grid)
	win.Show()
}
