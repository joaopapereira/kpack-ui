package gui

import (
	"log"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"

	"kpackui/builder"
)

func NewBuildersScreen(getter *builder.CustomClusterGetter) fyne.CanvasObject {
	_, err := getter.GetAll()
	if err != nil {
		log.Fatalf("cannot retrieve custom builders: %s", err.Error())
	}
	return widget.NewVBox()
}
