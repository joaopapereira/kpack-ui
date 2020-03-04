package gui

import (
	"log"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"image/color"

	"kpackui/builder"
	"kpackui/kpack"
)

func NewBuildersScreen(getter *builder.CustomClusterGetter) fyne.CanvasObject {
	builders, err := getter.GetAll()
	if err != nil {
		log.Fatalf("cannot retrieve custom builders: %s", err.Error())
	}
	container := fyne.NewContainerWithLayout(layout.NewGridLayout(3))

	for _, clusterBuilder := range builders {
		success := newBuilderSuccessBox(clusterBuilder)
		success.build()
		container.AddObject(
			success,
		)
	}

	return container
}

type builderBox struct {
	widget.Box
	background color.Color
	builder    kpack.CustomClusterBuilder
}

func (b *builderBox) build() {
	b.Children = append(b.Children,
		fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(160, 70)),
			widget.NewScrollContainer(fyne.NewContainerWithLayout(
				layout.NewGridLayout(1),
				widget.NewLabel(b.builder.Name),
				widget.NewLabel(b.builder.Tag),
			)),
		),
	)
}

func (b *builderBox) CreateRenderer() fyne.WidgetRenderer {
	b.ExtendBaseWidget(b)

	return &boxRenderer{objects: b.Children, layout: layout.NewVBoxLayout(), box: b}
}

func newBuilderSuccessBox(builder kpack.CustomClusterBuilder) *builderBox {
	return &builderBox{
		Box: widget.Box{
			BaseWidget: widget.BaseWidget{},
			Horizontal: false,
		},
		background: color.RGBA{
			R: 0xff,
			G: 0x0,
			B: 0x0,
			A: 0xff,
		},
		builder: builder,
	}
}

type boxRenderer struct {
	layout fyne.Layout

	objects []fyne.CanvasObject
	box     *builderBox
}

func (b *boxRenderer) MinSize() fyne.Size {
	return b.layout.MinSize(b.objects)
}

func (b *boxRenderer) Layout(size fyne.Size) {
	b.layout.Layout(b.objects, size)
}

func (b *boxRenderer) BackgroundColor() color.Color {
	if b.box.background == nil {
		return theme.BackgroundColor()
	}

	return b.box.background
}

func (b *boxRenderer) Objects() []fyne.CanvasObject {
	return b.objects
}

func (b *boxRenderer) Refresh() {
	b.objects = b.box.Children
	for _, child := range b.objects {
		child.Refresh()
	}
	b.Layout(b.box.Size())

	canvas.Refresh(b.box)
}

func (b *boxRenderer) Destroy() {
}
