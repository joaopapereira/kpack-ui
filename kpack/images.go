package kpack

import (
	"fmt"
	"log"
	"regexp"

	"github.com/astaxie/beego"
	bv1alpha1 "github.com/pivotal/kpack/pkg/apis/build/v1alpha1"
	"github.com/pivotal/kpack/pkg/client/clientset/versioned/typed/build/v1alpha1"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	duckv1alpha1 "knative.dev/pkg/apis/duck/v1alpha1"
)

func NewImageController(k8sClient corev1.CoreV1Interface, kpackClient v1alpha1.BuildV1alpha1Interface) *ImageController {
	repo := imagesRepo{
		k8sClient: k8sClient,
		imgClient: kpackClient,
	}
	return &ImageController{
		ImagesRepo: repo,
	}
}

type ImageController struct {
	beego.Controller
	ImagesRepo imagesRepo
}

func (c *ImageController) Get() {
	images, err := c.ImagesRepo.GetAll()
	if err != nil {
		log.Printf("error while getting images: %s", err.Error())
		c.Data["json"] = &errorMessage{Error: err.Error()}
	} else {
		c.Data["json"] = &images
	}
	c.ServeJSON()
}

type errorMessage struct {
	Error string `json:"error"`
}

type imagesRepo struct {
	k8sClient corev1.CoreV1Interface
	imgClient v1alpha1.BuildV1alpha1Interface
}

func (i *imagesRepo) GetAll() ([]Project, error) {
	allNamespaces, err := i.k8sClient.Namespaces().List(metav1.ListOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "retrieving from namespace")
	}
	var result []Project
	for _, ns := range allNamespaces.Items {
		if !validNamespace(ns.Name) {
			continue
		}

		allImages, err := i.imgClient.Images(ns.Name).List(metav1.ListOptions{})
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("unable to get images for namespace: %s", ns.Name))
		}
		project := Project{
			Name:   ns.Name,
			Images: nil,
		}
		for _, img := range allImages.Items {
			image := Image{
				Name:         img.Name,
				Tag:          img.Spec.Tag,
				LastBuiltTag: img.Status.LatestImage,
			}

			allBuilds, err := i.imgClient.Builds(ns.Name).List(metav1.ListOptions{
				TypeMeta:      metav1.TypeMeta{},
				LabelSelector: bv1alpha1.ImageLabel + "=" + image.Name,
			})
			if err != nil {
				return nil, errors.Wrap(err, fmt.Sprintf("unable to get builds from image '%s' on namespace: '%s'", img.Name, ns.Name))
			}

			for _, build := range allBuilds.Items {
				image.Builds = append(image.Builds, Build{
					Reason: buildReasons(build),
					Status: buildState(build),
					ID:     build.Labels[bv1alpha1.BuildNumberLabel],
				})
			}

			if img.Generation != img.Status.ObservedGeneration {
				image.Builds = append(image.Builds, Build{
					Status: "Pending",
				})
			}

			project.Images = append(project.Images, image)
		}
		result = append(result, project)
	}

	return result, nil
}

func buildReasons(build bv1alpha1.Build) string {
	if reasons, ok := build.Annotations[bv1alpha1.BuildReasonAnnotation]; ok {
		return reasons
	}
	return ""
}

func buildState(build bv1alpha1.Build) string {
	succeeded := build.Status.GetCondition(duckv1alpha1.ConditionSucceeded)
	switch {
	case succeeded.IsTrue():
		return "Succeeded"
	case succeeded.IsFalse():
		return "Failed"
	case succeeded.IsUnknown():
		return "Building"
	}

	return "Pending"
}

func validNamespace(namespace string) bool {
	invalidNamespaces := []string{"ingress-*", "pks-*", "kpack*", "kube-*", "default"}
	for _, invalidNamespace := range invalidNamespaces {
		match, _ := regexp.Match(invalidNamespace, []byte(namespace))
		if match {
			return false
		}
	}

	return true
}

type Project struct {
	Name   string  `json:"name"`
	Images []Image `json:"images"`
}

type Build struct {
	Reason string `json:"reason"`
	Status string `json:"status"`
	ID     string `json:"id"`
}

type Image struct {
	Name         string  `json:"name"`
	Tag          string  `json:"tag"`
	LastBuiltTag string  `json:"lastBuiltTag"`
	Builds       []Build `json:"builds"`
}
