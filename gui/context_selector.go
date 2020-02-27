package gui

import (
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/pkg/errors"
)

type ContextGetter interface {
	GetAll() ([]string, error)
}

func NewContextSelector() *contextSelector {
	return &contextSelector{}
}

type contextSelector struct {
	contextButtons []*widget.Button
}

func (c *contextSelector) Show(win fyne.Window, contextGetter ContextGetter, onContextSelected func(string), onError func(error)) {
	contexts, err := contextGetter.GetAll()
	if err != nil {
		onError(errors.Wrap(err, "on context select"))
		return
	}

	grid := fyne.NewContainerWithLayout(layout.NewGridLayout(1))
	group1 := widget.NewGroup("Contexts")
	grid.AddObject(group1)

	for _, name := range contexts {
		name := name
		button := widget.NewButton(name, func() {
			onContextSelected(name)
		})
		c.contextButtons = append(c.contextButtons, button)
		group1.Append(button)
	}

	win.SetContent(grid)
	win.Show()
}
