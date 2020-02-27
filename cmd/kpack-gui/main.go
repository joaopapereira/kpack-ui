package main

import (
	"fyne.io/fyne/app"
	"fyne.io/fyne/theme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"

	"kpackui/gui"
	"kpackui/k8s"
)

func main() {
	a := app.NewWithID("com.github.kpack-gui")
	a.SetIcon(theme.FyneLogo())

	gui.SelectContext(a, k8s.NewContextGetter())
}
