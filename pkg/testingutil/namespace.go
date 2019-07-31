package testingutil

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func getFakeNamespace(namespace string) corev1.Namespace {
	return corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}
}

// FakeNamespaces returns fake namespacelist for use by API mock functions for Unit tests
func FakeNamespaces() *corev1.NamespaceList {
	return &corev1.NamespaceList{
		Items: []corev1.Namespace{
			getFakeNamespace("testing"),
			getFakeNamespace("prj1"),
			getFakeNamespace("prj2"),
		},
	}
}

// FakeNamespaceStatus returns fake namespace status for use by mock watch on namespace
func FakeNamespaceStatus(prjStatus corev1.NamespacePhase, prjName string) *corev1.Namespace {
	return &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: prjName,
		},
		Status: corev1.NamespaceStatus{Phase: prjStatus},
	}
}

// FakeOnlyOneExistingNamespaces returns fake namespacelist with single namespace for use by API mock functions for Unit tests testing delete of the only available namespace
func FakeOnlyOneExistingNamespaces() *corev1.NamespaceList {
	return &corev1.NamespaceList{
		Items: []corev1.Namespace{
			getFakeNamespace("testing"),
		},
	}
}

// FakeRemoveNamespace removes the delete requested namespace from the list of namespace passed
func FakeRemoveNamespace(namespace string, namespaces *corev1.NamespaceList) *corev1.NamespaceList {
	for index, ns := range namespaces.Items {
		if ns.Name == namespace {
			namespaces.Items = append(namespaces.Items[:index], namespaces.Items[index+1:]...)
		}
	}
	return namespaces
}
