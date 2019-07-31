package kclient

import (
	fakeServiceCatalogClientSet "github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset/fake"
	fakeKubeClientset "k8s.io/client-go/kubernetes/fake"
)

// FakeClientset holds fake ClientSets
// this is returned by FakeNew to access methods of fake client sets
type FakeClientset struct {
	Kubernetes              *fakeKubeClientset.Clientset
	ServiceCatalogClientSet *fakeServiceCatalogClientSet.Clientset
}

// FakeNew creates new fake client for testing
// returns Client that is filled with fake clients and
// FakeClientSet that holds fake Clientsets to access Actions, Reactors etc... in fake client
func FakeNew() (*Client, *FakeClientset) {
	var client Client
	var fkclientset FakeClientset

	fkclientset.Kubernetes = fakeKubeClientset.NewSimpleClientset()
	client.KubeClient = fkclientset.Kubernetes

	fkclientset.ServiceCatalogClientSet = fakeServiceCatalogClientSet.NewSimpleClientset()

	return &client, &fkclientset
}
