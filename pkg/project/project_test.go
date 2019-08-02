package project

import (
	"os"
	"testing"

	"github.com/redhat-developer/odo-fork/pkg/kclient"
	"github.com/redhat-developer/odo-fork/pkg/testingutil"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	ktesting "k8s.io/client-go/testing"
	"k8s.io/client-go/tools/clientcmd"
)

func TestDelete(t *testing.T) {
	tests := []struct {
		name          string
		wantErr       bool
		namespaceName string
	}{
		{
			name:          "Test namespace delete for multiple namespaces",
			wantErr:       false,
			namespaceName: "prj2",
		},
		{
			name:          "Test delete the only remaining namespace",
			wantErr:       false,
			namespaceName: "testing",
		},
	}

	odoConfigFile, kubeConfigFile, err := testingutil.SetUp(
		testingutil.ConfigDetails{
			FileName:      "odo-test-config",
			Config:        testingutil.FakeOdoConfig("odo-test-config", false, ""),
			ConfigPathEnv: "GLOBALODOCONFIG",
		}, testingutil.ConfigDetails{
			FileName:      "kube-test-config",
			Config:        testingutil.FakeKubeClientConfig(),
			ConfigPathEnv: "KUBECONFIG",
		},
	)
	defer testingutil.CleanupEnv([]*os.File{odoConfigFile, kubeConfigFile}, t)
	if err != nil {
		t.Errorf("failed to create mock odo and kube config files. Error %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Fake the client with the appropriate arguments
			client, fakeClientSet := kclient.FakeNew()

			loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
			configOverrides := &clientcmd.ConfigOverrides{}
			client.KubeConfig = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

			client.Namespace = "testing"
			fkWatch := watch.NewFake()

			fakeClientSet.Kubernetes.PrependReactor("list", "namespaces", func(action ktesting.Action) (bool, runtime.Object, error) {
				if tt.name == "Test delete the only remaining namespace" {
					return true, testingutil.FakeOnlyOneExistingNamespaces(), nil
				}
				return true, testingutil.FakeNamespaces(), nil
			})

			fakeClientSet.Kubernetes.PrependReactor("delete", "namespaces", func(action ktesting.Action) (bool, runtime.Object, error) {
				return true, nil, nil
			})

			go func() {
				fkWatch.Delete(testingutil.FakeNamespaceStatus(corev1.NamespacePhase(""), tt.namespaceName))
			}()
			fakeClientSet.Kubernetes.PrependWatchReactor("namespaces", func(action ktesting.Action) (handled bool, ret watch.Interface, err error) {
				return true, fkWatch, nil
			})

			// The function we are testing
			err := Delete(client, tt.namespaceName)

			if err == nil && !tt.wantErr {
				if len(fakeClientSet.Kubernetes.Actions()) != 2 {
					t.Errorf("expected 2 ProjClientSet.Actions() in Project Delete, got: %v", len(fakeClientSet.Kubernetes.Actions()))
				}
			}

			// Checks for error in positive cases
			if !tt.wantErr == (err != nil) {
				t.Errorf("namespace Delete() unexpected error %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
