package kpackui

import (
	"github.com/pivotal/kpack/pkg/client/clientset/versioned/typed/build/v1alpha1"
	v1alpha12 "github.com/pivotal/kpack/pkg/client/clientset/versioned/typed/experimental/v1alpha1"
)

type KpackConnectionManager interface {
	GetKpack() v1alpha1.BuildV1alpha1Interface
	GetExperimentalKpack() v1alpha12.ExperimentalV1alpha1Interface
}
