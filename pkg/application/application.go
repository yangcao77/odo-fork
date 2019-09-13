package application

import (
	"github.com/pkg/errors"

	applabels "github.com/redhat-developer/odo-fork/pkg/application/labels"
	"github.com/redhat-developer/odo-fork/pkg/kclient"
	"github.com/redhat-developer/odo-fork/pkg/util"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	appPrefixMaxLen   = 12
	appNameMaxRetries = 3
	appAPIVersion     = "apps.udo.io/v1alpha1"
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

	// Get all Deployments with the "app" label
	deploymentAppNames, err := client.GetDeploymentLabelValues(applabels.ApplicationLabel, applabels.ApplicationLabel)
	if err != nil {
		return nil, errors.Wrap(err, "unable to list applications from deployment")
	}

	appNames = append(appNames, deploymentAppNames...)

	// Filter out any names, as there could be multiple components but within the same application
	return util.RemoveDuplicates(appNames), nil
}

// Exists checks whether the given app exists or not
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
