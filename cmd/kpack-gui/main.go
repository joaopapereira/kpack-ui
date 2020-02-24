package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/cmd/fyne_demo/data"
	"fyne.io/fyne/cmd/fyne_demo/screens"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/pivotal/kpack/pkg/client/clientset/versioned"
	"github.com/pkg/errors"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	"kpackui/gui"
	"kpackui/k8s"
)

const preferenceCurrentTab = "currentTab"

func welcomeScreen(a fyne.App) fyne.CanvasObject {
	logo := canvas.NewImageFromResource(data.FyneScene)
	logo.SetMinSize(fyne.NewSize(228, 167))

	link, err := url.Parse("https://fyne.io/")
	if err != nil {
		fyne.LogError("Could not parse URL", err)
	}

	return widget.NewVBox(
		widget.NewLabelWithStyle("Welcome to the Fyne toolkit demo app", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		layout.NewSpacer(),
		widget.NewHBox(layout.NewSpacer(), logo, layout.NewSpacer()),
		widget.NewHyperlinkWithStyle("fyne.io", link, fyne.TextAlignCenter, fyne.TextStyle{}),
		layout.NewSpacer(),

		widget.NewGroup("Theme",
			fyne.NewContainerWithLayout(layout.NewGridLayout(2),
				widget.NewButton("Dark", func() {
					a.Settings().SetTheme(theme.DarkTheme())
				}),
				widget.NewButton("Light", func() {
					a.Settings().SetTheme(theme.LightTheme())
				}),
			),
		),
	)
}

func main() {
	a := app.NewWithID("com.github.kpack-gui")
	a.SetIcon(theme.FyneLogo())

	selectContext(a, k8s.NewContextGetter())

	//fyneDemo(a)
}

func selectContext(a fyne.App, getter gui.ContextGetter) {
	w := a.NewWindow("Kpack gui - Context Selector")
	gui.ContextSelector(w, getter, func(name string) {
		displayKpackForContext(a, getter, name)
		w.Close()
	})
	w.ShowAndRun()
}

func displaySpinner(w fyne.Window, title string) {
	progressBar := widget.NewProgressBarInfinite()

	w.SetContent(widget.NewVBox(
		widget.NewLabel(title), progressBar))
}

func displayKpackForContext(a fyne.App, getter gui.ContextGetter, context string) {
	w := a.NewWindow(fmt.Sprintf("Kpack gui - %s", context))

	displaySpinner(w, fmt.Sprintf("connecting to the %s cluster", context))

	go func() {
		time.Sleep(3 * time.Second)
		if connectToCluster(context) != nil {
			w.SetContent(widget.NewVBox(
				widget.NewLabel(fmt.Sprintf("Not authorized to connect to %s cluster", context)),
				widget.NewButton("Select a different cluster", func() {
					selectContext(a, getter)
					w.Close()
				})))
			return
		}

		tabs := widget.NewTabContainer(
			widget.NewTabItemWithIcon("Welcome", theme.HomeIcon(), welcomeScreen(a)),
		)
		tabs.SetTabLocation(widget.TabLocationLeading)
		tabs.SelectTabIndex(a.Preferences().Int(preferenceCurrentTab))
		w.SetContent(tabs)
		a.Preferences().SetInt(preferenceCurrentTab, tabs.CurrentTabIndex())
	}()

	w.Show()
}

type unauthorized struct {}

func (u unauthorized) Error() string {
	return "not authorized"
}

func connectToCluster(context string) error {
	cfg, err := retrieveLocalConfiguration()
	if err != nil {
		log.Fatalf("unable to retrieve local config: %s", err.Error())
	}

	clientCfg := clientcmd.NewNonInteractiveClientConfig(*cfg, context, &clientcmd.ConfigOverrides{}, nil)
	cliCfg, err := clientCfg.ClientConfig()
	if err != nil {
		log.Fatalf("unable to retrieve client config: %s", err.Error())
	}
	clientCfg.Namespace()

	kpackClientSet, err := versioned.NewForConfig(cliCfg)
	if err != nil {
		log.Fatalf("unable to create kpack client set: %s", err.Error())
	}

	builders, err := kpackClientSet.ExperimentalV1alpha1().CustomClusterBuilders().List(v1.ListOptions{})
	if err != nil {
		if k8serrors.IsUnauthorized(errors.Cause(err)) {
			log.Printf("cluster unauthorized: %s", err.Error())
			return unauthorized{}
		} else {
			log.Fatalf("unable to retrieve builders: %s", err.Error())
		}
	}
	for _, builder := range builders.Items {
		log.Printf("builder: %s", builder.Name)
	}

	return nil
}

func retrieveLocalConfiguration() (*clientcmdapi.Config, error) {
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	return clientcmd.LoadFromFile(*kubeconfig)
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func fyneDemo(a fyne.App) {
	w := a.NewWindow("Fyne Demo")
	w.SetMainMenu(fyne.NewMainMenu(fyne.NewMenu("File",
		fyne.NewMenuItem("New", func() { fmt.Println("Menu New") }),
		// a quit item will be appended to our first menu
	), fyne.NewMenu("Edit",
		fyne.NewMenuItem("Cut", func() { fmt.Println("Menu Cut") }),
		fyne.NewMenuItem("Copy", func() { fmt.Println("Menu Copy") }),
		fyne.NewMenuItem("Paste", func() { fmt.Println("Menu Paste") }),
	)))
	w.SetMaster()

	tabs := widget.NewTabContainer(
		widget.NewTabItemWithIcon("Welcome", theme.HomeIcon(), welcomeScreen(a)),
		widget.NewTabItemWithIcon("Widgets", theme.ContentCopyIcon(), screens.WidgetScreen()),
		widget.NewTabItemWithIcon("Graphics", theme.DocumentCreateIcon(), screens.GraphicsScreen()),
		widget.NewTabItemWithIcon("Windows", theme.ViewFullScreenIcon(), screens.DialogScreen(w)),
		widget.NewTabItemWithIcon("Advanced", theme.SettingsIcon(), screens.AdvancedScreen(w)))
	tabs.SetTabLocation(widget.TabLocationLeading)
	tabs.SelectTabIndex(a.Preferences().Int(preferenceCurrentTab))
	w.SetContent(tabs)

	w.ShowAndRun()
	a.Preferences().SetInt(preferenceCurrentTab, tabs.CurrentTabIndex())
}
