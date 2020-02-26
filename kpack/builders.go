package kpack

import (
	"github.com/pivotal/kpack/pkg/apis/build/v1alpha1"
	b_v1alpha1 "github.com/pivotal/kpack/pkg/client/clientset/versioned/typed/build/v1alpha1"
	e_v1alpha1 "github.com/pivotal/kpack/pkg/client/clientset/versioned/typed/experimental/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Buildpack struct {
	Id      string
	Version string
}

type CustomClusterBuilder struct {
	BuiltSuccess bool
	Tag          string
	Image        string
	Stack        string
	Store        string
	Buildpacks   []Buildpack
}

func NewBuilderRepo(kpackClient b_v1alpha1.BuildV1alpha1Interface, experimentalKpackClient e_v1alpha1.ExperimentalV1alpha1Interface) *builderRepo {
	repo := builderRepo{
		buildClient:        kpackClient,
		experimentalClient: experimentalKpackClient,
	}
	return &repo
}

type builderRepo struct {
	buildClient        b_v1alpha1.BuildV1alpha1Interface
	experimentalClient e_v1alpha1.ExperimentalV1alpha1Interface
}

func (b builderRepo) GetAllCustomClusterBuilders() ([]CustomClusterBuilder, error) {
	builders, err := b.experimentalClient.CustomClusterBuilders().List(v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var customBuilders []CustomClusterBuilder

	for _, builder := range builders.Items {
		customBuilder := CustomClusterBuilder{
			Tag:   builder.Spec.Tag,
			Store: builder.Spec.Store,
		}
		if builder.Status.GetCondition(v1alpha1.ConditionBuilderReady).IsTrue() {
			var buildpacks []Buildpack
			for _, metadata := range builder.Status.BuilderMetadata {
				buildpacks = append(buildpacks, Buildpack{
					Id:      metadata.Id,
					Version: metadata.Version,
				})
			}
			customBuilder.Buildpacks = buildpacks
			customBuilder.BuiltSuccess = true
			customBuilder.Image = builder.Status.LatestImage
			customBuilder.Stack = builder.Status.Stack.RunImage
		} else {

		}
		customBuilders = append(customBuilders, customBuilder)
	}

	return customBuilders, nil
}
