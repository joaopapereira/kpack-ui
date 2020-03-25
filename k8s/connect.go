package k8s

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/pivotal/kpack/pkg/client/clientset/versioned"
	kbuild "github.com/pivotal/kpack/pkg/client/clientset/versioned/typed/build/v1alpha1"
	kexp "github.com/pivotal/kpack/pkg/client/clientset/versioned/typed/experimental/v1alpha1"
	"github.com/pkg/errors"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	core_v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	"kpackui"
)

func ConnectToCluster(context string) (kpackui.KpackConnectionManager, error) {
	cfg, err := retrieveLocalConfiguration()
	if err != nil {
		return nil, errors.Wrap(err, "unable to retrieve local config")
	}

	clientCfg := clientcmd.NewNonInteractiveClientConfig(*cfg, context, &clientcmd.ConfigOverrides{}, nil)
	cliCfg, err := clientCfg.ClientConfig()
	if err != nil {
		return nil, errors.Wrap(err, "unable to retrieve client config")
	}
	clientCfg.Namespace()

	kpackClientSet, err := versioned.NewForConfig(cliCfg)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create kpack client set")
	}

	_, err = kpackClientSet.ExperimentalV1alpha1().CustomClusterBuilders().List(v1.ListOptions{})
	if err != nil {
		if k8serrors.IsUnauthorized(errors.Cause(err)) {
			log.Printf("cluster unauthorized: %s", err.Error())
			return nil, errors.Wrap(err, "retrieve cluster cluster info")
		}
		return nil, errors.Wrap(err, "error retrieving information from the cluster")
	}
	k8sClientSet, err := kubernetes.NewForConfig(cliCfg)

	if err != nil {
		log.Printf("unable to create configuration: %s", err)
		return nil, errors.Wrap(err, "generate new config")
	}

	return &connectionManager{
		kpackExperimental: kpackClientSet.ExperimentalV1alpha1(),
		kpackInterface:    kpackClientSet.BuildV1alpha1(),
		k8sClient:         k8sClientSet.CoreV1(),
	}, nil
}

type connectionManager struct {
	kpackExperimental kexp.ExperimentalV1alpha1Interface
	kpackInterface    kbuild.BuildV1alpha1Interface
	k8sClient         core_v1.CoreV1Interface
}

func (c connectionManager) GetExperimentalKpack() kexp.ExperimentalV1alpha1Interface {
	return c.kpackExperimental
}

func (c connectionManager) GetKpack() kbuild.BuildV1alpha1Interface {
	return c.kpackInterface
}

func (c connectionManager) GetCorev1() core_v1.CoreV1Interface {
	return c.k8sClient
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
