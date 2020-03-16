package gui

import (
	"fmt"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"kpackui"
	"kpackui/builder"
	"kpackui/kpack"
)

const preferenceCurrentTab = "currentTab"
const builderCurrentTab = "currentTab"

func NewKpackMainView(getter ContextGetter, connectionBuilder func(context string) (kpackui.KpackConnectionManager, error)) *KpackMainView {
	return &KpackMainView{
		getter:            getter,
		connectionBuilder: connectionBuilder,
	}
}

type KpackMainView struct {
	getter            ContextGetter
	connectionManager kpackui.KpackConnectionManager
	connectionBuilder func(_ string) (kpackui.KpackConnectionManager, error)
}

func (v *KpackMainView) LoadUI(a fyne.App, context string, onConnectionFailure func()) {
	var err error
	w := a.NewWindow(fmt.Sprintf("Kpack gui - %s", context))

	ShowSpinner(w, fmt.Sprintf("connecting to the %s cluster", context))

	w.Show()

	v.connectionManager, err = v.connectionBuilder(context)
	if err != nil {
		w.SetContent(widget.NewVBox(
			widget.NewLabel(fmt.Sprintf("Not authorized to connect to %s cluster", context)),
			widget.NewButton("Select a different cluster", func() {
				onConnectionFailure()
				w.Hide()
			})))
		return
	}

	builderRepo := kpack.NewBuilderRepo(
		v.connectionManager.GetKpack(),
		v.connectionManager.GetExperimentalKpack(),
	)

	tabs := widget.NewTabContainer(
		widget.NewTabItemWithIcon("Builders", theme.HomeIcon(),
			builderMenu(a, builderRepo),
		),
	)
	tabs.SetTabLocation(widget.TabLocationLeading)
	tabs.SelectTabIndex(a.Preferences().Int(preferenceCurrentTab))
	w.SetContent(tabs)
	a.Preferences().SetInt(preferenceCurrentTab, tabs.CurrentTabIndex())
}

func builderMenu(app fyne.App, builderRepo *kpack.BuilderRepo) *fyne.Container {
	tabs := widget.NewTabContainer(
		widget.NewTabItem("Custom Cluster", NewBuildersScreen(
			builder.NewCustomClusterGetter(
				builderRepo))),
		widget.NewTabItem("Cluster", NewBuildersScreen(
			builder.NewClusterGetter(
				builderRepo))),
	)
	tabs.SetTabLocation(widget.TabLocationLeading)
	tabs.SelectTabIndex(app.Preferences().Int(builderCurrentTab))

	app.Preferences().SetInt(builderCurrentTab, tabs.CurrentTabIndex())
	return fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, nil, nil, nil),
		tabs,
	)
}

//func welcomeScreen(a fyne.App) fyne.CanvasObject {
//	logo := canvas.NewImageFromResource(data.FyneScene)
//	logo.SetMinSize(fyne.NewSize(228, 167))
//
//	link, err := url.Parse("https://fyne.io/")
//	if err != nil {
//		fyne.LogError("Could not parse URL", err)
//	}
//
//	return widget.NewVBox(
//		widget.NewLabelWithStyle("Welcome to the Fyne toolkit demo app", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
//		layout.NewSpacer(),
//		widget.NewHBox(layout.NewSpacer(), logo, layout.NewSpacer()),
//		widget.NewHyperlinkWithStyle("fyne.io", link, fyne.TextAlignCenter, fyne.TextStyle{}),
//		layout.NewSpacer(),
//
//		widget.NewGroup("Theme",
//			fyne.NewContainerWithLayout(layout.NewGridLayout(2),
//				widget.NewButton("Dark", func() {
//					a.Settings().SetTheme(theme.DarkTheme())
//				}),
//				widget.NewButton("Light", func() {
//					a.Settings().SetTheme(theme.LightTheme())
//				}),
//			),
//		),
//	)
//}

//func fyneDemo(a fyne.App) {
//	w := a.NewWindow("Fyne Demo")
//	w.SetMainMenu(fyne.NewMainMenu(fyne.NewMenu("File",
//		fyne.NewMenuItem("New", func() { fmt.Println("Menu New") }),
//		// a quit item will be appended to our first menu
//	), fyne.NewMenu("Edit",
//		fyne.NewMenuItem("Cut", func() { fmt.Println("Menu Cut") }),
//		fyne.NewMenuItem("Copy", func() { fmt.Println("Menu Copy") }),
//		fyne.NewMenuItem("Paste", func() { fmt.Println("Menu Paste") }),
//	)))
//	w.SetMaster()
//
//	tabs := widget.NewTabContainer(
//		widget.NewTabItemWithIcon("Welcome", theme.HomeIcon(), welcomeScreen(a)),
//		widget.NewTabItemWithIcon("Widgets", theme.ContentCopyIcon(), screens.WidgetScreen()),
//		widget.NewTabItemWithIcon("Graphics", theme.DocumentCreateIcon(), screens.GraphicsScreen()),
//		widget.NewTabItemWithIcon("Windows", theme.ViewFullScreenIcon(), screens.DialogScreen(w)),
//		widget.NewTabItemWithIcon("Advanced", theme.SettingsIcon(), screens.AdvancedScreen(w)))
//	tabs.SetTabLocation(widget.TabLocationLeading)
//	tabs.SelectTabIndex(a.Preferences().Int(preferenceCurrentTab))
//	w.SetContent(tabs)
//
//	w.ShowAndRun()
//	a.Preferences().SetInt(preferenceCurrentTab, tabs.CurrentTabIndex())
//}
