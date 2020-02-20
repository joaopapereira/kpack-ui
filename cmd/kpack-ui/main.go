package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/astaxie/beego"
	"github.com/pivotal/kpack/pkg/client/clientset/versioned"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"kpackui"
	"kpackui/kpack"
)

func main() {
	log.Printf("Starting kpackui")
	beego.SetStaticPath("/", "/kpack-ui/ui/dist")

	var (
		clusterConfig *rest.Config
		err           error
	)
	if os.Getenv("LOCAL_START") != "" {
		clusterConfig, err = retrieveLocalConfiguration()
		if err != nil {
			log.Fatalf("unable to retrieve local configuration: %s", err.Error())
		}
		beego.SetStaticPath("/", "ui/dist")
	} else {
		clusterConfig, err = retrieveClusterConfiguration()
		if err != nil {
			log.Fatalf("unable to retrieve cluster configuration: %s", err.Error())
		}
	}

	kpackClientSet, err := versioned.NewForConfig(clusterConfig)
	if err != nil {
		log.Fatalf("could not get cnbbuild clientset: %s", err.Error())
	}
	k8sClient, err := kubernetes.NewForConfig(clusterConfig)
	if err != nil {
		log.Fatalf("could not get cnbbuild clientset: %s", err.Error())
	}

	beego.Router("/images", kpack.NewImageController(k8sClient.CoreV1(), kpackClientSet.BuildV1alpha1()))
	beego.Router("/", &kpackui.HomeController{})
	beego.Run()
}

func retrieveLocalConfiguration() (*rest.Config, error) {
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
	// use the current context in kubeconfig
	clusterConfig, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	return clusterConfig, err
}

func retrieveClusterConfiguration() (*rest.Config, error) {
	clusterConfig, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("could not get kubernetes in-cluster config: %s", err.Error())
	}
	return clusterConfig, err
}
func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
