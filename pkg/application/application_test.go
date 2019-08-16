package application

import (
	"reflect"
	"testing"

	applabels "github.com/redhat-developer/odo-fork/pkg/application/labels"
	"github.com/redhat-developer/odo-fork/pkg/component"
	componentlabels "github.com/redhat-developer/odo-fork/pkg/component/labels"
	"github.com/redhat-developer/odo-fork/pkg/kclient"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ktesting "k8s.io/client-go/testing"
)

func TestGetMachineReadableFormat(t *testing.T) {
	type args struct {
		// client      *kclient.Client
		appName     string
		projectName string
		active      bool
	}
	tests := []struct {
		name string
		args args
		want App
	}{
		{

			name: "Test Case: machine readable output for application",
			args: args{
				appName:     "myapp",
				projectName: "myproject",
				active:      true,
			},
			want: App{
				TypeMeta: metav1.TypeMeta{
					Kind:       appKind,
					APIVersion: appAPIVersion,
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "myapp",
					Namespace: "myproject",
				},
				Spec: AppSpec{
					Components: []string{"frontend"},
				},
			},
		},
	}

	dcList := appsv1.DeploymentList{
		Items: []appsv1.Deployment{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "frontend-myapp",
					Namespace: "myproject",
					Labels: map[string]string{
						applabels.ApplicationLabel:         "myapp",
						componentlabels.ComponentLabel:     "frontend",
						componentlabels.ComponentTypeLabel: "nodejs",
					},
					Annotations: map[string]string{
						component.ComponentSourceTypeAnnotation: "local",
					},
				},
				Spec: appsv1.DeploymentSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name: "dummyContainer",
								},
							},
						},
					},
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "backend-app",
					Namespace: "myproject",
					Labels: map[string]string{
						applabels.ApplicationLabel:         "app",
						componentlabels.ComponentLabel:     "backend",
						componentlabels.ComponentTypeLabel: "java",
					},
					Annotations: map[string]string{
						component.ComponentSourceTypeAnnotation: "local",
					},
				},
				Spec: appsv1.DeploymentSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name: "dummyContainer",
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Fake the client with the appropriate arguments
			client, fakeClientSet := kclient.FakeNew()
			//fake the dcs
			fakeClientSet.Kubernetes.PrependReactor("list", "deploymentconfigs", func(action ktesting.Action) (bool, runtime.Object, error) {
				return true, &dcList, nil
			})

			for i := range dcList.Items {
				fakeClientSet.Kubernetes.PrependReactor("get", "deploymentconfigs", func(action ktesting.Action) (bool, runtime.Object, error) {
					return true, &dcList.Items[i], nil
				})
			}
			if got := GetMachineReadableFormat(client, tt.args.appName, tt.args.projectName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMachineReadableFormat() = %v,\n want %v", got, tt.want)
			}
		})
	}
}

func TestGetMachineReadableFormatForList(t *testing.T) {
	type args struct {
		apps []App
	}
	tests := []struct {
		name string
		args args
		want AppList
	}{
		{
			name: "Test Case: Machine Readable for Application List",
			args: args{
				apps: []App{
					{
						TypeMeta: metav1.TypeMeta{
							Kind:       appKind,
							APIVersion: appAPIVersion,
						},
						ObjectMeta: metav1.ObjectMeta{
							Name: "myapp",
						},
						Spec: AppSpec{
							Components: []string{"frontend"},
						},
						Status: AppStatus{
							Active: true,
						},
					},
				},
			},
			want: AppList{
				TypeMeta: metav1.TypeMeta{
					Kind:       appList,
					APIVersion: appAPIVersion,
				},
				ListMeta: metav1.ListMeta{},
				Items: []App{
					{
						TypeMeta: metav1.TypeMeta{
							Kind:       appKind,
							APIVersion: appAPIVersion,
						},
						ObjectMeta: metav1.ObjectMeta{
							Name: "myapp",
						},
						Spec: AppSpec{
							Components: []string{"frontend"},
						},
						Status: AppStatus{
							Active: true,
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetMachineReadableFormatForList(tt.args.apps); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMachineReadableFormatForList() = %v, want %v", got, tt.want)
			}
		})
	}
}
