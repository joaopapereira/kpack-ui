package kpack_test

import (
	"testing"

	"github.com/pivotal/kpack/pkg/apis/build/v1alpha1"
	corev1alpha1 "github.com/pivotal/kpack/pkg/apis/core/v1alpha1"
	v1alpha12 "github.com/pivotal/kpack/pkg/apis/experimental/v1alpha1"
	kpack "github.com/pivotal/kpack/pkg/client/clientset/versioned/fake"
	"github.com/sclevine/spec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"

	kpack2 "kpackui/kpack"
)

func TestBuilders(t *testing.T) {
	spec.Run(t, "Test Builders", testBuilders)
}

func testBuilders(t *testing.T, when spec.G, it spec.S) {
	var (
		fakeClient = kpack.NewSimpleClientset()
		subject    = kpack2.NewBuilderRepo(fakeClient.BuildV1alpha1(), fakeClient.ExperimentalV1alpha1())
	)

	when("#GetAllCustomClusterBuilders", func() {
		it("returns empty list when no custom builders are present", func() {
			builders, err := subject.GetAllCustomClusterBuilders()
			require.NoError(t, err)
			require.Len(t, builders, 0)
		})

		when("custom builders are present", func() {
			const namespace = "nice-namespace"
			it.Before(func() {
				var err error
				_, err = k8sClientFake.CoreV1().Namespaces().Create(&v1.Namespace{
					ObjectMeta: v12.ObjectMeta{
						Name: namespace,
					},
				})
				require.NoError(t, err)
			})

			it("returns projects with an image when namespace as an images", func() {
				var _, err = fakeClient.ExperimentalV1alpha1().CustomClusterBuilders().Create(&v1alpha12.CustomClusterBuilder{
					ObjectMeta: v12.ObjectMeta{
						Name: "custom-builder",
					},
					Spec: v1alpha12.CustomClusterBuilderSpec{
						CustomBuilderSpec: v1alpha12.CustomBuilderSpec{
							Tag:   "some/custom-builder:tag",
							Stack: "io.buildpacks.java",
							Store: "some/store:tag",
							Order: nil,
						},
					},
					Status: v1alpha12.CustomBuilderStatus{
						BuilderStatus: v1alpha1.BuilderStatus{
							Status: corev1alpha1.Status{
								ObservedGeneration: 0,
								Conditions: []corev1alpha1.Condition{
								},
							},
							BuilderMetadata: nil,
							Stack:           v1alpha1.BuildStack{},
							LatestImage:     "",
						},
					},
				})
				require.NoError(t, err)

				projects, err := subject.GetAll()
				require.NoError(t, err)
				require.Len(t, projects, 1)
				require.Len(t, projects[0].Images, 1)
				assert.Equal(t, kpack2.Image{
					Name:         "nice-image",
					Tag:          "some/image:tag",
					LastBuiltTag: "some/image:tag@sha256:123123123",
					Builds: []kpack2.Build{
						{
							Reason: "CONFIG",
							Status: "Succeeded",
							ID:     "1",
						},
					},
				}, projects[0].Images[0])
			})
		})
	})
}
