package gui

import (
	"sort"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/pkg/errors"
)

type ContextGetter interface {
	GetAll() ([]string, error)
}

func SelectContext(a fyne.App, getter ContextGetter, onContextSelected func(string, fyne.Window), onError func(error) *fyne.Container) fyne.Window {
	w := a.NewWindow("Kpack gui - Context Selector")
	contextSelector := NewContextSelector()
	contextSelector.Show(w, getter, func(name string) {
		onContextSelected(name, w)
		w.Hide()
	}, func(err error) {
		w.SetContent(
			onError(err),
		)
	})
	w.ShowAndRun()
	return w
}

func NewContextSelector() *ContextSelector {
	return &ContextSelector{}
}

type ContextSelector struct {
	contextButtons []*widget.Button
}

func (c *ContextSelector) Show(win fyne.Window, contextGetter ContextGetter, onContextSelected func(string), onError func(error)) {
	contexts, err := contextGetter.GetAll()
	if err != nil {
		onError(errors.Wrap(err, "on context select"))
		return
	}

	grid := fyne.NewContainerWithLayout(layout.NewGridLayout(1))
	group1 := widget.NewGroup("Contexts")
	grid.AddObject(group1)

	sort.Strings(contexts)
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
