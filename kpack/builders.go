package kpack

import (
	"github.com/pivotal/kpack/pkg/apis/build/v1alpha1"
	b_v1alpha1 "github.com/pivotal/kpack/pkg/client/clientset/versioned/typed/build/v1alpha1"
	e_v1alpha1 "github.com/pivotal/kpack/pkg/client/clientset/versioned/typed/experimental/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Buildpack struct {
	ID      string
	Version string
}

type CustomClusterBuilder struct {
	BuiltSuccess bool
	Name         string
	Tag          string
	Image        string
	Stack        string
	Store        string
	Buildpacks   []Buildpack
}

func NewBuilderRepo(kpackClient b_v1alpha1.BuildV1alpha1Interface, experimentalKpackClient e_v1alpha1.ExperimentalV1alpha1Interface) *BuilderRepo {
	repo := BuilderRepo{
		buildClient:        kpackClient,
		experimentalClient: experimentalKpackClient,
	}
	return &repo
}

type BuilderRepo struct {
	buildClient        b_v1alpha1.BuildV1alpha1Interface
	experimentalClient e_v1alpha1.ExperimentalV1alpha1Interface
}

func (b BuilderRepo) GetAllCustomClusterBuilders() ([]CustomClusterBuilder, error) {
	builders, err := b.experimentalClient.CustomClusterBuilders().List(v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var customBuilders []CustomClusterBuilder

	for _, builder := range builders.Items {
		customBuilder := CustomClusterBuilder{
			Tag:   builder.Spec.Tag,
			Store: builder.Spec.Store,
			Name:  builder.Name,
		}
		if builder.Status.GetCondition(v1alpha1.ConditionBuilderReady).IsTrue() {
			var buildpacks []Buildpack
			for _, metadata := range builder.Status.BuilderMetadata {
				buildpacks = append(buildpacks, Buildpack{
					ID:      metadata.Id,
					Version: metadata.Version,
				})
			}
			customBuilder.Buildpacks = buildpacks
			customBuilder.BuiltSuccess = true
			customBuilder.Image = builder.Status.LatestImage
			customBuilder.Stack = builder.Status.Stack.RunImage
		}

		customBuilders = append(customBuilders, customBuilder)
	}

	return customBuilders, nil
}
