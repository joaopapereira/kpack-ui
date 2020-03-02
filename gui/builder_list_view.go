package gui

import (
	"log"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	"kpackui/builder"
)

func NewBuildersScreen(getter *builder.CustomClusterGetter) fyne.CanvasObject {
	builders, err := getter.GetAll()
	if err != nil {
		log.Fatalf("cannot retrieve custom builders: %s", err.Error())
	}
	container := fyne.NewContainerWithLayout(layout.NewGridLayout(3))

	for _, clusterBuilder := range builders {
		container.AddObject(
			fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(160, 70)),
				widget.NewScrollContainer(fyne.NewContainerWithLayout(
					layout.NewGridLayout(1),
					widget.NewLabel(clusterBuilder.Name),
					widget.NewLabel(clusterBuilder.Tag),
				)),
			),
		)
	}

	return container
}
