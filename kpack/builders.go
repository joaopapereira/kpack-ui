package kpack

import (
	b_v1alpha1 "github.com/pivotal/kpack/pkg/client/clientset/versioned/typed/build/v1alpha1"
	e_v1alpha1 "github.com/pivotal/kpack/pkg/client/clientset/versioned/typed/experimental/v1alpha1"
)

type CustomClusterBuilder struct {
	Tag   string `json:"tag,omitempty"`
	Stack string `json:"stack,omitempty"`
	Store string `json:"store,omitempty"`
	// +listType
	Order []struct {
		Group []struct {
			Id       string `json:"id"`
			Version  string `json:"version,omitempty"`
			Optional bool   `json:"optional,omitempty"`
		} `json:"group,omitempty"`
	} `json:"order,omitempty"`
}

type BuilderRepo interface {
	GetAllCustomClusterBuilders() (CustomClusterBuilder, error)
}

func NewBuilderRepo(kpackClient b_v1alpha1.BuildV1alpha1Interface, experimentalKpackClient e_v1alpha1.ExperimentalV1alpha1Interface) BuilderRepo {
	repo := builderRepo{
		buildClient: kpackClient,
		experimentalClient: experimentalKpackClient,
	}
	return &repo
}

type builderRepo struct {
	buildClient        b_v1alpha1.BuildV1alpha1Interface
	experimentalClient e_v1alpha1.ExperimentalV1alpha1Interface
}

func (b builderRepo) GetAllCustomClusterBuilders() (CustomClusterBuilder, error) {
	panic("implement me")
}
