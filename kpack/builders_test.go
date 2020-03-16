package kpack_test

import (
	"testing"

	v1 "k8s.io/api/core/v1"

	"github.com/pivotal/kpack/pkg/apis/build/v1alpha1"
	corev1alpha1 "github.com/pivotal/kpack/pkg/apis/core/v1alpha1"
	exp_v1alpha1 "github.com/pivotal/kpack/pkg/apis/experimental/v1alpha1"
	kpack "github.com/pivotal/kpack/pkg/client/clientset/versioned/fake"
	"github.com/sclevine/spec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"

	kpack2 "kpackui/kpack"
)

func TestBuilders(t *testing.T) {
	spec.Run(t, "Test Builders", testBuilders)
}

func testBuilders(t *testing.T, when spec.G, it spec.S) {
	var (
		fakeClient         = kpack.NewSimpleClientset()
		experimentalClient = fakeClient.ExperimentalV1alpha1()
		kpackClient        = fakeClient.BuildV1alpha1()
		subject            = kpack2.NewBuilderRepo(kpackClient, experimentalClient)
	)

	when("#GetAllCustomClusterBuilders", func() {
		it("returns empty list when no custom builders are present", func() {
			builders, err := subject.GetAllCustomClusterBuilders()
			require.NoError(t, err)
			require.Len(t, builders, 0)
		})

		when("custom builders are present", func() {
			it("returns builder with an image when builder was built successfully", func() {
				var _, err = experimentalClient.CustomClusterBuilders().Create(&exp_v1alpha1.CustomClusterBuilder{
					ObjectMeta: v12.ObjectMeta{
						Name: "custom-builder",
					},
					Spec: exp_v1alpha1.CustomClusterBuilderSpec{
						CustomBuilderSpec: exp_v1alpha1.CustomBuilderSpec{
							Tag:   "some/custom-builder:tag",
							Stack: "io.buildpacks.java",
							Store: "some/store:tag",
							Order: []exp_v1alpha1.OrderEntry{
								{
									Group: []exp_v1alpha1.BuildpackRef{
										{
											BuildpackInfo: exp_v1alpha1.BuildpackInfo{
												Id:      "io.buildpack.java",
												Version: "1.0.0",
											},
											Optional: false,
										},
									},
								},
							},
						},
					},
					Status: exp_v1alpha1.CustomBuilderStatus{
						BuilderStatus: v1alpha1.BuilderStatus{
							Status: corev1alpha1.Status{
								ObservedGeneration: 0,
								Conditions: []corev1alpha1.Condition{
									{
										Type:   v1alpha1.ConditionBuilderReady,
										Status: v1.ConditionTrue,
									},
								},
							},
							BuilderMetadata: v1alpha1.BuildpackMetadataList{
								{
									Id:      "io.buildpack.java",
									Version: "1.0.0",
								},
							},
							Stack: v1alpha1.BuildStack{
								RunImage: "some/stack:image",
								ID:       "io.buildpacks.stack",
							},
							LatestImage: "some/custom-builder:tag@098223ad",
						},
					},
				})
				require.NoError(t, err)

				builders, err := subject.GetAllCustomClusterBuilders()
				require.NoError(t, err)
				require.Len(t, builders, 1)
				assert.Equal(t, builders[0].Name(), "custom-builder")
				assert.Equal(t, builders[0].Image, "some/custom-builder:tag@098223ad")
				assert.Equal(t, builders[0].Tag(), "some/custom-builder:tag")
				assert.Equal(t, builders[0].Stack, "some/stack:image")
				assert.Equal(t, builders[0].Store, "some/store:tag")
				assert.Equal(t, builders[0].Buildpacks, []kpack2.Buildpack{
					{
						ID:      "io.buildpack.java",
						Version: "1.0.0",
					},
				})
				assert.True(t, builders[0].BuiltSuccess)
			})

			it("returns builder without an image when builder failed to build", func() {
				var _, err = experimentalClient.CustomClusterBuilders().Create(&exp_v1alpha1.CustomClusterBuilder{
					ObjectMeta: v12.ObjectMeta{
						Name: "custom-builder",
					},
					Spec: exp_v1alpha1.CustomClusterBuilderSpec{
						CustomBuilderSpec: exp_v1alpha1.CustomBuilderSpec{
							Tag:   "some/custom-builder:tag",
							Stack: "io.buildpacks.java",
							Store: "some/store:tag",
							Order: []exp_v1alpha1.OrderEntry{
								{
									Group: []exp_v1alpha1.BuildpackRef{
										{
											BuildpackInfo: exp_v1alpha1.BuildpackInfo{
												Id:      "io.buildpack.java",
												Version: "1.0.0",
											},
											Optional: false,
										},
									},
								},
							},
						},
					},
					Status: exp_v1alpha1.CustomBuilderStatus{
						BuilderStatus: v1alpha1.BuilderStatus{
							Status: corev1alpha1.Status{
								ObservedGeneration: 0,
								Conditions: []corev1alpha1.Condition{
									{
										Type:   v1alpha1.ConditionBuilderReady,
										Status: v1.ConditionFalse,
									},
								},
							},
						},
					},
				})
				require.NoError(t, err)

				builders, err := subject.GetAllCustomClusterBuilders()
				require.NoError(t, err)
				require.Len(t, builders, 1)
				assert.Equal(t, builders[0].Image, "")
				assert.Equal(t, builders[0].Tag(), "some/custom-builder:tag")
				assert.Equal(t, builders[0].Stack, "")
				assert.Equal(t, builders[0].Store, "some/store:tag")
				assert.Nil(t, builders[0].Buildpacks)
				assert.False(t, builders[0].BuiltSuccess)
			})

			it("returns builder without an image when builder still hasnt finished building", func() {
				var _, err = experimentalClient.CustomClusterBuilders().Create(&exp_v1alpha1.CustomClusterBuilder{
					ObjectMeta: v12.ObjectMeta{
						Name: "custom-builder",
					},
					Spec: exp_v1alpha1.CustomClusterBuilderSpec{
						CustomBuilderSpec: exp_v1alpha1.CustomBuilderSpec{
							Tag:   "some/custom-builder:tag",
							Stack: "io.buildpacks.java",
							Store: "some/store:tag",
							Order: []exp_v1alpha1.OrderEntry{
								{
									Group: []exp_v1alpha1.BuildpackRef{
										{
											BuildpackInfo: exp_v1alpha1.BuildpackInfo{
												Id:      "io.buildpack.java",
												Version: "1.0.0",
											},
											Optional: false,
										},
									},
								},
							},
						},
					},
					Status: exp_v1alpha1.CustomBuilderStatus{
						BuilderStatus: v1alpha1.BuilderStatus{
							Status: corev1alpha1.Status{
								ObservedGeneration: 0,
								Conditions: []corev1alpha1.Condition{
									{
										Type:   v1alpha1.ConditionBuilderReady,
										Status: v1.ConditionUnknown,
									},
								},
							},
						},
					},
				})
				require.NoError(t, err)

				builders, err := subject.GetAllCustomClusterBuilders()
				require.NoError(t, err)
				require.Len(t, builders, 1)
				assert.Equal(t, builders[0].Name(), "custom-builder")
				assert.Equal(t, builders[0].Image, "")
				assert.Equal(t, builders[0].Tag(), "some/custom-builder:tag")
				assert.Equal(t, builders[0].Stack, "")
				assert.Equal(t, builders[0].Store, "some/store:tag")
				assert.Nil(t, builders[0].Buildpacks)
				assert.False(t, builders[0].BuiltSuccess)
			})
		})
	})

	when("#GetAllClusterBuilders", func() {
		it("returns empty list when no builders are present", func() {
			builders, err := subject.GetAllClusterBuilders()
			require.NoError(t, err)
			require.Len(t, builders, 0)
		})

		when("cluster builders are present", func() {
			it("returns builder with an image when builder was built successfully", func() {
				var _, err = kpackClient.ClusterBuilders().Create(&v1alpha1.ClusterBuilder{
					ObjectMeta: v12.ObjectMeta{
						Name: "cluster-builder",
					},
					Spec: v1alpha1.BuilderSpec{
						Image:        "some/cluster-builder:tag",
						UpdatePolicy: v1alpha1.External,
					},
					Status: v1alpha1.BuilderStatus{
						Status: corev1alpha1.Status{
							ObservedGeneration: 0,
							Conditions: []corev1alpha1.Condition{
								{
									Type:   v1alpha1.ConditionBuilderReady,
									Status: v1.ConditionTrue,
								},
							},
						},
						BuilderMetadata: v1alpha1.BuildpackMetadataList{
							{
								Id:      "io.buildpack.java",
								Version: "1.0.0",
							},
						},
						Stack: v1alpha1.BuildStack{
							RunImage: "some/stack:image",
							ID:       "io.buildpacks.stack",
						},
						LatestImage: "some/cluster-builder:tag@098223ad",
					},
				})
				require.NoError(t, err)

				builders, err := subject.GetAllClusterBuilders()
				require.NoError(t, err)
				require.Len(t, builders, 1)
				assert.Equal(t, builders[0].Name(), "cluster-builder")
				assert.Equal(t, builders[0].Image, "some/cluster-builder:tag@098223ad")
				assert.Equal(t, builders[0].Tag(), "some/cluster-builder:tag")
				assert.Equal(t, builders[0].Stack, "some/stack:image")
				assert.Equal(t, builders[0].Store, "")
				assert.Equal(t, builders[0].Buildpacks, []kpack2.Buildpack{
					{
						ID:      "io.buildpack.java",
						Version: "1.0.0",
					},
				})
				assert.True(t, builders[0].BuiltSuccess)
			})

			it("returns builder without an image when builder failed to build", func() {
				var _, err = kpackClient.ClusterBuilders().Create(&v1alpha1.ClusterBuilder{
					ObjectMeta: v12.ObjectMeta{
						Name: "cluster-builder",
					},
					Spec: v1alpha1.BuilderSpec{
						Image:        "some/cluster-builder:tag",
						UpdatePolicy: v1alpha1.External,
					},
					Status: v1alpha1.BuilderStatus{
						Status: corev1alpha1.Status{
							ObservedGeneration: 0,
							Conditions: []corev1alpha1.Condition{
								{
									Type:   v1alpha1.ConditionBuilderReady,
									Status: v1.ConditionFalse,
								},
							},
						},
					},
				})
				require.NoError(t, err)

				builders, err := subject.GetAllClusterBuilders()
				require.NoError(t, err)
				require.Len(t, builders, 1)
				assert.Equal(t, builders[0].Image, "")
				assert.Equal(t, builders[0].Tag(), "some/cluster-builder:tag")
				assert.Equal(t, builders[0].Stack, "")
				assert.Equal(t, builders[0].Store, "")
				assert.Nil(t, builders[0].Buildpacks)
				assert.False(t, builders[0].BuiltSuccess)
			})

			it("returns builder without an image when builder still hasnt finished building", func() {
				var _, err = kpackClient.ClusterBuilders().Create(&v1alpha1.ClusterBuilder{
					ObjectMeta: v12.ObjectMeta{
						Name: "cluster-builder",
					},
					Spec: v1alpha1.BuilderSpec{
						Image:        "some/cluster-builder:tag",
						UpdatePolicy: v1alpha1.External,
					},
					Status: v1alpha1.BuilderStatus{
						Status: corev1alpha1.Status{
							ObservedGeneration: 0,
							Conditions: []corev1alpha1.Condition{
								{
									Type:   v1alpha1.ConditionBuilderReady,
									Status: v1.ConditionUnknown,
								},
							},
						},
					},
				})
				require.NoError(t, err)

				builders, err := subject.GetAllClusterBuilders()
				require.NoError(t, err)
				require.Len(t, builders, 1)
				assert.Equal(t, builders[0].Name(), "cluster-builder")
				assert.Equal(t, builders[0].Image, "")
				assert.Equal(t, builders[0].Tag(), "some/cluster-builder:tag")
				assert.Equal(t, builders[0].Stack, "")
				assert.Equal(t, builders[0].Store, "")
				assert.Nil(t, builders[0].Buildpacks)
				assert.False(t, builders[0].BuiltSuccess)
			})
		})
	})
}
