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

type NamespacedBuilders struct {
	ClusterBuilder
	Namespace string
}

type ClusterBuilder struct {
	BuiltSuccess bool
	name         string
	tag          string
	Image        string
	Stack        string
	Store        string
	Buildpacks   []Buildpack
}

func (b *ClusterBuilder) Name() string {
	return b.name
}

func (b *ClusterBuilder) Tag() string {
	return b.tag
}

func (b *ClusterBuilder) BuiltSuccessful() bool {
	return b.BuiltSuccess
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

func (b BuilderRepo) GetAllCustomClusterBuilders() ([]ClusterBuilder, error) {
	builders, err := b.experimentalClient.CustomClusterBuilders().List(v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var customBuilders []ClusterBuilder

	for _, builder := range builders.Items {
		customBuilder := ClusterBuilder{
			tag:   builder.Spec.Tag,
			Store: builder.Spec.Store,
			name:  builder.Name,
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

func (b BuilderRepo) GetAllClusterBuilders() ([]ClusterBuilder, error) {
	builders, err := b.buildClient.ClusterBuilders().List(v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var clusterBuilders []ClusterBuilder

	for _, builder := range builders.Items {
		clusterBuilder := ClusterBuilder{
			tag:   builder.Spec.Image,
			Store: "",
			name:  builder.Name,
		}

		if builder.Status.GetCondition(v1alpha1.ConditionBuilderReady).IsTrue() {
			var buildpacks []Buildpack
			for _, metadata := range builder.Status.BuilderMetadata {
				buildpacks = append(buildpacks, Buildpack{
					ID:      metadata.Id,
					Version: metadata.Version,
				})
			}
			clusterBuilder.Buildpacks = buildpacks
			clusterBuilder.BuiltSuccess = true
			clusterBuilder.Image = builder.Status.LatestImage
			clusterBuilder.Stack = builder.Status.Stack.RunImage
		}

		clusterBuilders = append(clusterBuilders, clusterBuilder)
	}

	return clusterBuilders, nil
}

func (b BuilderRepo) GetAllCustomBuilders(namespace string) ([]NamespacedBuilders, error) {
	builders, err := b.experimentalClient.CustomBuilders(namespace).List(v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var customBuilders []NamespacedBuilders

	for _, builder := range builders.Items {
		customBuilder := NamespacedBuilders{
			ClusterBuilder: ClusterBuilder{
				tag:   builder.Spec.Tag,
				Store: builder.Spec.Store,
				name:  builder.Name,
			},
			Namespace: namespace,
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

func (b BuilderRepo) GetAllNamespacedBuilders(namespace string) ([]NamespacedBuilders, error) {
	builders, err := b.buildClient.Builders(namespace).List(v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var clusterBuilders []NamespacedBuilders

	for _, builder := range builders.Items {
		clusterBuilder := NamespacedBuilders{
			ClusterBuilder: ClusterBuilder{
				tag:   builder.Spec.Image,
				Store: "",
				name:  builder.Name,
			},
			Namespace: namespace,
		}

		if builder.Status.GetCondition(v1alpha1.ConditionBuilderReady).IsTrue() {
			var buildpacks []Buildpack
			for _, metadata := range builder.Status.BuilderMetadata {
				buildpacks = append(buildpacks, Buildpack{
					ID:      metadata.Id,
					Version: metadata.Version,
				})
			}
			clusterBuilder.Buildpacks = buildpacks
			clusterBuilder.BuiltSuccess = true
			clusterBuilder.Image = builder.Status.LatestImage
			clusterBuilder.Stack = builder.Status.Stack.RunImage
		}

		clusterBuilders = append(clusterBuilders, clusterBuilder)
	}

	return clusterBuilders, nil
}