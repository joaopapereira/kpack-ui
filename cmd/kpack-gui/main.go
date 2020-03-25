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
	//mainView := gui.NewKpackMainView(contextGetter, stubbedConnectionManager())

	//gui.FyneDemo(a)
	var contextView fyne.Window
	contextView = gui.SelectContext(a, contextGetter, func(context string) {
		mainView.LoadUI(a, context, func() {
			contextView.Show() // SegFault........
		})
	}, gui.ErrorContainer)
}

//func stubbedConnectionManager() func(context string) (manager kpackui.KpackConnectionManager, err error) {
//	return func(context string) (manager kpackui.KpackConnectionManager, err error) {
//		connectionManager := kpackui.DummyKpackConnectionManager{}
//		connectionManager.GetExperimentalKpack().CustomClusterBuilders().Create(&v1alpha1.CustomClusterBuilder{
//			ObjectMeta: v1.ObjectMeta{
//				Name: "custom-builder",
//			},
//			Spec: v1alpha1.CustomClusterBuilderSpec{
//				CustomBuilderSpec: v1alpha1.CustomBuilderSpec{
//					Tag:   "some/custom-builder:tag",
//					Stack: "io.buildpacks.java",
//					Store: "some/store:tag",
//					Order: []v1alpha1.OrderEntry{
//						{
//							Group: []v1alpha1.BuildpackRef{
//								{
//									BuildpackInfo: v1alpha1.BuildpackInfo{
//										Id:      "io.buildpack.java",
//										Version: "1.0.0",
//									},
//									Optional: false,
//								},
//							},
//						},
//					},
//				},
//			},
//			Status: v1alpha1.CustomBuilderStatus{
//				BuilderStatus: k_v1alpha1.BuilderStatus{
//					Status: corev1alpha1.Status{
//						ObservedGeneration: 0,
//						Conditions: []corev1alpha1.Condition{
//							{
//								Type:   k_kpack_v1.ConditionReady,
//								Status: core_v1.ConditionTrue,
//							},
//						},
//					},
//					BuilderMetadata: k_v1alpha1.BuildpackMetadataList{
//						{
//							Id:      "io.buildpack.java",
//							Version: "1.0.0",
//						},
//					},
//					Stack: k_v1alpha1.BuildStack{
//						RunImage: "some/stack:image",
//						ID:       "io.buildpacks.stack",
//					},
//					LatestImage: "some/custom-builder:tag@098223ad",
//				},
//			},
//		})
//		return &connectionManager, nil
//	}
//}
