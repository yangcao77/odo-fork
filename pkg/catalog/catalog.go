package catalog

import (
	"fmt"
	"strings"

	"github.com/golang/glog"
	// imagev1 "github.com/openshift/api/image/v1"
	"github.com/pkg/errors"
	"github.com/redhat-developer/odo-fork/pkg/kclient"
)

type CatalogImage struct {
	Name          string
	Namespace     string
	AllTags       []string
	NonHiddenTags []string
}

// List lists all the available component types
func List(client *kclient.Client) ([]CatalogImage, error) {

	// TODO: implement for KDO
	fake := CatalogImage{"nodejs", "fakeNS", []string{"latest"}, []string{}}
	catalogList := []CatalogImage{fake}

	if len(catalogList) == 0 {
		return nil, errors.New("unable to retrieve any catalog images from the OpenShift cluster")
	}

	return catalogList, nil
}

// Search searches for the component
func Search(client *kclient.Client, name string) ([]string, error) {
	var result []string
	componentList, err := List(client)
	if err != nil {
		return nil, errors.Wrap(err, "unable to list components")
	}

	// do a partial search in all the components
	for _, component := range componentList {
		// we only show components that contain the search term and that have at least non-hidden tag
		// since a component with all hidden tags is not shown in the odo catalog list components either
		if strings.Contains(component.Name, name) && len(component.NonHiddenTags) > 0 {
			result = append(result, component.Name)
		}
	}

	return result, nil
}

// Exists returns true if the given component type is valid, false if not
func Exists(client *kclient.Client, componentType string) (bool, error) {

	catalogList, err := List(client)
	if err != nil {
		return false, errors.Wrapf(err, "unable to list catalog")
	}

	for _, supported := range catalogList {
		if componentType == supported.Name || componentType == fmt.Sprintf("%s/%s", supported.Namespace, supported.Name) {
			return true, nil
		}
	}
	return false, nil
}

// VersionExists checks if that version exists, returns true if the given version exists, false if not
func VersionExists(client *kclient.Client, componentType string, componentVersion string) (bool, error) {

	// Loading status
	glog.V(4).Info("Checking component version")

	// Retrieve the catalogList
	catalogList, err := List(client)
	if err != nil {
		return false, errors.Wrapf(err, "unable to list catalog")
	}

	// Find the component and then return true if the version has been found
	for _, supported := range catalogList {
		if componentType == supported.Name || componentType == fmt.Sprintf("%s/%s", supported.Namespace, supported.Name) {
			// Now check to see if that version matches that components tag
			// here we use the AllTags, because if the user somehow got hold of a version that was hidden
			// then it's safe to assume that this user went to a lot of trouble to actually use that version,
			// so let's allow it
			for _, tag := range supported.AllTags {
				if componentVersion == tag {
					return true, nil
				}
			}
		}
	}

	// Else return false if nothing is found
	return false, nil
}
