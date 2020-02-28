package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/theme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"

	"kpackui/gui"
	"kpackui/k8s"
)

func main() {
	a := app.NewWithID("com.github.kpack-gui")
	a.SetIcon(theme.FyneLogo())

	contextGetter := k8s.NewContextGetter()
	mainView := gui.NewKpackMainView(contextGetter, k8s.ConnectToCluster)
	var contextView fyne.Window
	contextView = gui.SelectContext(a, contextGetter, func(context string) {
		mainView.LoadUI(a, context, func() {
			contextView.Show() // SegFault........
		})
	}, gui.ErrorContainer)
}
