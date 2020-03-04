package kpackui

import (
	kpack "github.com/pivotal/kpack/pkg/client/clientset/versioned/fake"
	"github.com/pivotal/kpack/pkg/client/clientset/versioned/typed/build/v1alpha1"
	v1alpha12 "github.com/pivotal/kpack/pkg/client/clientset/versioned/typed/experimental/v1alpha1"
)

type KpackConnectionManager interface {
	GetKpack() v1alpha1.BuildV1alpha1Interface
	GetExperimentalKpack() v1alpha12.ExperimentalV1alpha1Interface
}

type DummyKpackConnectionManager struct {
	fakeClient *kpack.Clientset
	expKpack   v1alpha12.ExperimentalV1alpha1Interface
	buildKpack v1alpha1.BuildV1alpha1Interface
}

func (d DummyKpackConnectionManager) GetKpack() v1alpha1.BuildV1alpha1Interface {
	if d.fakeClient == nil {
		d.fakeClient = kpack.NewSimpleClientset()
		d.expKpack = d.fakeClient.ExperimentalV1alpha1()
		d.buildKpack = d.fakeClient.BuildV1alpha1()
	}

	return d.buildKpack
}

func (d DummyKpackConnectionManager) GetExperimentalKpack() v1alpha12.ExperimentalV1alpha1Interface {
	if d.fakeClient == nil {
		d.fakeClient = kpack.NewSimpleClientset()
		d.expKpack = d.fakeClient.ExperimentalV1alpha1()
		d.buildKpack = d.fakeClient.BuildV1alpha1()
	}

	return d.expKpack
}
