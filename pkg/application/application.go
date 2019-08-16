package application

import (
	// "github.com/golang/glog"
	"github.com/pkg/errors"

	applabels "github.com/redhat-developer/odo-fork/pkg/application/labels"
	"github.com/redhat-developer/odo-fork/pkg/kclient"
	"github.com/redhat-developer/odo-fork/pkg/util"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	appPrefixMaxLen   = 12
	appNameMaxRetries = 3
	appAPIVersion     = "odo.openshift.io/v1alpha1"
	appKind           = "app"
	appList           = "List"
)

// List all applications in current project
func List(client *kclient.Client) ([]string, error) {
	return ListInProject(client)
}

// ListInProject lists all applications in given project by Querying the cluster
func ListInProject(client *kclient.Client) ([]string, error) {

	var appNames []string

	// Get all DeploymentConfigs with the "app" label
	deploymentConfigAppNames, err := client.GetDeploymentLabelValues(applabels.ApplicationLabel, applabels.ApplicationLabel)
	if err != nil {
		return nil, errors.Wrap(err, "unable to list applications from deployment config")
	}

	appNames = append(appNames, deploymentConfigAppNames...)

	// Removed following ODO code, need to evaluate whether KDO needs an equivalent lookup
	// // Get all ServiceInstances with the "app" label
	// // Okay, so there is an edge-case here.. if Service Catalog is *not* enabled in the cluster, we shouldn't error out..
	// // however, we should at least warn the user.
	// serviceInstanceAppNames, err := client.GetServiceInstanceLabelValues(applabels.ApplicationLabel, applabels.ApplicationLabel)
	// if err != nil {
	// 	glog.V(4).Infof("Unable to list Service Catalog instances: %s", err)
	// 	log.Warning("Unable to access Service Catalog instances, may not be enabled on cluster")
	// } else {
	// 	appNames = append(deploymentConfigAppNames, serviceInstanceAppNames...)
	// }

	// Filter out any names, as there could be multiple components but within the same application
	return util.RemoveDuplicates(appNames), nil
}

// Exists checks whether the given app exist or not
func Exists(app string, client *kclient.Client) (bool, error) {

	appList, err := List(client)
	if err != nil {
		return false, err
	}
	for _, appName := range appList {
		if appName == app {
			return true, nil
		}
	}
	return false, nil
}

// TODO-KDO: Add when implementing application commands
// // Delete deletes the given application
// func Delete(client *kclient.Client, name string) error {
// 	glog.V(4).Infof("Deleting application %s", name)

// 	labels := applabels.GetLabels(name, false)

// 	// delete application from cluster
// 	err := client.Delete(labels)
// 	if err != nil {
// 		return errors.Wrapf(err, "unable to delete application %s", name)
// 	}

// 	return nil
// }

// TODO-KDO: Add when implementing application commands
// // GetMachineReadableFormat returns resource information in machine readable format
// func GetMachineReadableFormat(client *kclient.Client, appName string, projectName string) App {
// 	componentList, _ := component.List(client, appName)
// 	var compList []string
// 	for _, comp := range componentList.Items {
// 		compList = append(compList, comp.Name)
// 	}
// 	appDef := App{
// 		TypeMeta: metav1.TypeMeta{
// 			Kind:       appKind,
// 			APIVersion: appAPIVersion,
// 		},
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:      appName,
// 			Namespace: projectName,
// 		},
// 		Spec: AppSpec{
// 			Components: compList,
// 		},
// 	}
// 	return appDef
// }

// GetMachineReadableFormatForList returns application list in machine readable format
func GetMachineReadableFormatForList(apps []App) AppList {
	return AppList{
		TypeMeta: metav1.TypeMeta{
			Kind:       appList,
			APIVersion: appAPIVersion,
		},
		ListMeta: metav1.ListMeta{},
		Items:    apps,
	}
}
