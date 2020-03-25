package k8s

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

func NewNamespaceGetter(k8sClient corev1.CoreV1Interface) *NamespaceGetter {
	return &NamespaceGetter{
		k8sClient: k8sClient,
	}
}

type NamespaceGetter struct {
	k8sClient corev1.CoreV1Interface
}

func (g *NamespaceGetter) GetNamespaces() ([]string, error) {
	ns, err := g.k8sClient.Namespaces().List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var namespaces []string
	for _, namespace := range ns.Items {
		namespaces = append(namespaces, namespace.Name)
	}

	return namespaces, nil
}
