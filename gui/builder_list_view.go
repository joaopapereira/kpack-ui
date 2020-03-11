package gui

import (
	"image/color"
	"log"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"kpackui/builder"
)

type KpackBuilder interface {
	Name() string
	Tag() string
}

func NewBuildersScreen(getter *builder.CustomClusterGetter) fyne.CanvasObject {
	builders, err := getter.GetAll()
	if err != nil {
		log.Fatalf("cannot retrieve custom builders: %s", err.Error())
	}
	container := fyne.NewContainerWithLayout(layout.NewGridLayout(3))

	for _, clusterBuilder := range builders {
		container.AddObject(
			&builderSuccessWidget{
				builder: &clusterBuilder,
			},
		)
	}

	return container
}

var (
	green = &color.RGBA{R: 0, G: 128, B: 0, A: 255}
)

type builderSuccessWidget struct {
	widget.BaseWidget
	background color.Color
	builder    KpackBuilder
}

func (b *builderSuccessWidget) MinSize() fyne.Size {
	b.ExtendBaseWidget(b)
	return b.BaseWidget.MinSize()
}

func (b *builderSuccessWidget) Refresh() {
	b.background = green

	b.BaseWidget.Refresh()
}

func (b *builderSuccessWidget) CreateRenderer() fyne.WidgetRenderer {
	b.ExtendBaseWidget(b)

	name := canvas.NewText(b.builder.Name(), color.Black)
	tag := canvas.NewText(b.builder.Tag(), color.Black)
	background := canvas.NewRectangle(green)
	return &builderWidgetRenderer{
		builderName: name,
		builderTag:  tag,
		objects: []fyne.CanvasObject{
			background,
			name,
			tag,
		},
		background: background,
		builderWidget: b,
	}
}

type builderWidgetRenderer struct {
	builderTag  *canvas.Text
	builderName *canvas.Text
	background  *canvas.Rectangle

	builderWidget *builderSuccessWidget
	objects []fyne.CanvasObject
}

func (b *builderWidgetRenderer) MinSize() fyne.Size {
	return fyne.NewSize(160+theme.Padding()*2, 70+theme.Padding()*2)
}

func (b *builderWidgetRenderer) Layout(size fyne.Size) {
	inner := size.Subtract(fyne.NewSize(theme.Padding()*2, theme.Padding()*2))
	b.background.Resize(inner)
	b.background.Move(fyne.NewPos(0, 0))

	textSize := int(float32(size.Height) * .1)
	textMin := fyne.CurrentApp().Driver().RenderedTextSize(b.builderName.Text, textSize, fyne.TextStyle{Bold: false})
	b.builderName.TextSize = textSize
	b.builderName.Resize(fyne.NewSize(size.Width, textMin.Height))
	b.builderName.Move(fyne.NewPos(0, textMin.Height))

	textMin = fyne.CurrentApp().Driver().RenderedTextSize(b.builderTag.Text, textSize, fyne.TextStyle{Bold: false})
	b.builderTag.TextSize = textSize
	b.builderTag.Resize(fyne.NewSize(size.Width, textMin.Height))
	b.builderTag.Move(fyne.NewPos(0, size.Height-textMin.Height))
}

func (b *builderWidgetRenderer) BackgroundColor() color.Color {
	if b.builderWidget.background == nil {
		return green
	}

	return b.builderWidget.background
}

func (b *builderWidgetRenderer) Objects() []fyne.CanvasObject {
	return b.objects
}

func (b *builderWidgetRenderer) Refresh() {
	b.Layout(b.builderWidget.Size())

	canvas.Refresh(b.builderWidget)
}

func (b *builderWidgetRenderer) Destroy() {
}
