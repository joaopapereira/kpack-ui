package gui

import (
	"log"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"

	"kpackui/kpack"
)

func NewBuilderView(app fyne.App, builder KpackBuilder) {
	w := app.NewWindow("Builder")
	log.Printf("builder: %s", builder.Name())
	clusterBuilder, ok := builder.(*kpack.ClusterBuilder)
	if !ok {
		log.Fatal("Not a builder...... :(")
	}
	var fields []*widget.FormItem

	fields = append(fields, widget.NewFormItem("Name", widget.NewLabel(clusterBuilder.Name())))
	fields = append(fields, widget.NewFormItem("Tag", widget.NewLabel(clusterBuilder.Tag())))
	fields = append(fields, widget.NewFormItem("Generated Image", widget.NewLabel(clusterBuilder.Image)))
	fields = append(fields, widget.NewFormItem("Stack Used", widget.NewLabel(clusterBuilder.Stack)))
	fields = append(fields, widget.NewFormItem("Store Used", widget.NewLabel(clusterBuilder.Store)))
	fields = append(fields, widget.NewFormItem("Last Build Successful", widget.NewLabel(builtText(clusterBuilder))))

	var buildpacks []*widget.FormItem
	for _, buildpack := range clusterBuilder.Buildpacks {
		buildpacks = append(buildpacks, widget.NewFormItem(buildpack.ID, widget.NewLabel(buildpack.Version)))
	}
	fields = append(fields, widget.NewFormItem("Buildpacks", widget.NewForm(buildpacks...)))

	content := fyne.NewContainer(widget.NewForm(fields...))

	w.SetContent(content)
	w.Show()
}

func builtText(builder *kpack.ClusterBuilder) string {
	if builder.BuiltSuccess {
		return "Success"
	}
	return "Failed"
}
