package kpack_test

import (
	"fmt"
	"testing"

	"github.com/pivotal/kpack/pkg/apis/build/v1alpha1"
	kpack "github.com/pivotal/kpack/pkg/client/clientset/versioned/fake"
	"github.com/sclevine/spec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8s "k8s.io/client-go/kubernetes/fake"
	duckv1alpha1 "knative.dev/pkg/apis/duck/v1alpha1"

	kpack2 "kpackui/kpack"
)

func TestProjects(t *testing.T) {
	spec.Run(t, "Test Projects", test)
}

func test(t *testing.T, when spec.G, it spec.S) {
	var (
		k8sClientFake = k8s.NewSimpleClientset()
		fakeClient    = kpack.NewSimpleClientset(&v1alpha1.Image{})
		subject       = kpack2.NewProjectsRepo(k8sClientFake.CoreV1(), fakeClient.BuildV1alpha1())
	)

	when("#GetAll", func() {
		it("returns empty list when no namespaces exist", func() {
			projects, err := subject.GetAll()
			require.NoError(t, err)
			require.Len(t, projects, 0)
		})
		var invalidNamespaces = []string{"ingress-nginx", "pks-something", "kpack", "kpack-ui", "kube-node", "default"}
		for _, namespace := range invalidNamespaces {
			namespace := namespace
			it(fmt.Sprintf("returns empty when namespace %s is the only namespace", namespace), func() {
				k8sClientFake.CoreV1().Namespaces().Create(&v1.Namespace{
					ObjectMeta: v12.ObjectMeta{
						Name: namespace,
					},
				})
				projects, err := subject.GetAll()
				require.NoError(t, err)
				require.Len(t, projects, 0)
			})
		}

		when("valid namespace exists", func() {
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

			it("returns projects with no images when namespace as no images", func() {
				projects, err := subject.GetAll()
				require.NoError(t, err)
				require.Len(t, projects, 1)
				require.Len(t, projects[0].Images, 0)
			})

			it("returns projects with an image when namespace as an images", func() {
				_, err := fakeClient.BuildV1alpha1().Images(namespace).Create(&v1alpha1.Image{
					ObjectMeta: v12.ObjectMeta{
						Name:      "nice-image",
						Namespace: namespace,
					},
					Spec: v1alpha1.ImageSpec{
						Tag: "some/image:tag",
					},
					Status: v1alpha1.ImageStatus{
						LatestImage: "some/image:tag@sha256:123123123",
					},
				})
				require.NoError(t, err)
				fakeClient.BuildV1alpha1().Builds(namespace).Create(&v1alpha1.Build{
					ObjectMeta: v12.ObjectMeta{
						Name:      "some-build",
						Namespace: namespace,
						Annotations: map[string]string{
							v1alpha1.BuildReasonAnnotation: "CONFIG",
						},
						Labels: map[string]string{
							v1alpha1.ImageLabel:       "nice-image",
							v1alpha1.BuildNumberLabel: "1",
						},
					},
					Status: v1alpha1.BuildStatus{
						Status: duckv1alpha1.Status{
							Conditions: duckv1alpha1.Conditions{
								{
									Type:   duckv1alpha1.ConditionSucceeded,
									Status: v1.ConditionTrue,
								},
							},
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
