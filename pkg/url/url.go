package url

import (
	"fmt"
	"strings"

	// routev1 "github.com/openshift/api/route/v1"
	"github.com/pkg/errors"
	applabels "github.com/redhat-developer/odo-fork/pkg/application/labels"
	componentlabels "github.com/redhat-developer/odo-fork/pkg/component/labels"
	"github.com/redhat-developer/odo-fork/pkg/kclient"
	urlLabels "github.com/redhat-developer/odo-fork/pkg/url/labels"
	"github.com/redhat-developer/odo-fork/pkg/util"
	iextensionsv1 "k8s.io/api/extensions/v1beta1"

	"github.com/golang/glog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// Get returns URL defination for given URL name
func (urls UrlList) Get(urlName string) Url {
	for _, url := range urls.Items {
		if url.Name == urlName {
			return url
		}
	}
	return Url{}

}

// Delete deletes a URL
func Delete(client *kclient.Client, urlName string, applicationName string) error {

	// Namespace the URL name
	namespaceKubernetesObject, err := util.NamespaceKubernetesObject(urlName, applicationName)
	if err != nil {
		return errors.Wrapf(err, "unable to create namespaced name")
	}

	return client.DeleteIngress(namespaceKubernetesObject)
}

// Create creates a URL and returns url string and error if any
// portNumber is the target port number for the ingress and is -1 in case no port number is specified in which case it is automatically detected for components which expose only one service port)
func Create(client *kclient.Client, urlName string, portNumber int, ingressDomain string, componentName, applicationName string) (string, error) {
	fmt.Println("IN Create function")
	labels := urlLabels.GetLabels(urlName, componentName, applicationName, false)

	serviceName, err := util.NamespaceKubernetesObject(componentName, applicationName)
	if err != nil {
		return "", errors.Wrapf(err, "unable to create namespaced name")
	}

	urlName, err = util.NamespaceKubernetesObject(urlName, applicationName)
	if err != nil {
		return "", errors.Wrapf(err, "unable to create namespaced name")
	}
	// Pass in the namespace name, link to the service (componentName) and labels to create a ingress
	ingress, err := client.CreateIngress(urlName, serviceName, ingressDomain, intstr.FromInt(portNumber), labels)
	if err != nil {
		return "", errors.Wrap(err, "unable to create ingress")
	}

	return GetURLString(getProtocol(*ingress), ingressDomain), nil
}

// List lists the URLs in an application. The results can further be narrowed
// down if a component name is provided, which will only list URLs for the
// given component
func List(client *kclient.Client, componentName string, applicationName string) (UrlList, error) {

	labelSelector := fmt.Sprintf("%v=%v", applabels.ApplicationLabel, applicationName)

	if componentName != "" {
		labelSelector = labelSelector + fmt.Sprintf(",%v=%v", componentlabels.ComponentLabel, componentName)
	}

	glog.V(4).Infof("Listing ingresses with label selector: %v", labelSelector)
	ingresses, err := client.ListIngresses(labelSelector)
	if err != nil {
		return UrlList{}, errors.Wrap(err, "unable to list ingress names")
	}

	var urls []Url
	for _, i := range ingresses {
		a := getMachineReadableFormat(i)
		urls = append(urls, a)
	}

	urlList := getMachineReadableFormatForList(urls)
	return urlList, nil
}

func getProtocol(ingress iextensionsv1.Ingress) string {
	if ingress.Spec.TLS != nil {
		return "https"
	}
	return "http"

}

// GetURLString returns a string representation of given url
func GetURLString(protocol, URL string) string {
	return protocol + "://" + URL
}

// Exists checks if the url exists in the component or not
// urlName is the name of the url for checking
// componentName is the name of the component to which the url's existence is checked
// applicationName is the name of the application to which the url's existence is checked
func Exists(client *kclient.Client, urlName string, componentName string, applicationName string) (bool, error) {
	urls, err := List(client, componentName, applicationName)
	if err != nil {
		return false, errors.Wrap(err, "unable to list the urls")
	}

	for _, url := range urls.Items {
		if url.Name == urlName {
			return true, nil
		}
	}
	return false, nil
}

// GetComponentServicePortNumbers returns the port numbers exposed by the service of the component
// componentName is the name of the component
// applicationName is the name of the application
func GetComponentServicePortNumbers(client *kclient.Client, componentName string, applicationName string) ([]int, error) {
	componentLabels := componentlabels.GetLabels(componentName, applicationName, false)
	componentSelector := util.ConvertLabelsToSelector(componentLabels)

	services, err := client.GetServicesFromSelector(componentSelector)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get the service")
	}

	var ports []int

	for _, service := range services {
		for _, port := range service.Spec.Ports {
			ports = append(ports, int(port.Port))
		}
	}

	return ports, nil
}

// GetURLName returns a url name from the component name and the given port number
func GetURLName(componentName string, componentPort int) string {
	if componentPort == -1 {
		return componentName
	}
	return fmt.Sprintf("%v-%v", componentName, componentPort)
}

// GetValidPortNumber checks if the given port number is a valid component port or not
// if port number is not provided and the component is a single port component, the component port is returned
// port number is -1 if the user does not specify any port
func GetValidPortNumber(client *kclient.Client, portNumber int, componentName string, applicationName string) (int, error) {
	componentPorts, err := GetComponentServicePortNumbers(client, componentName, applicationName)
	if err != nil {
		return portNumber, errors.Wrapf(err, "unable to get exposed ports for component %s", componentName)
	}

	// port number will be -1 if the user doesn't specify any port
	if portNumber == -1 {

		switch {
		case len(componentPorts) > 1:
			return portNumber, errors.Errorf("port for the component %s is required as it exposes %d ports: %s", componentName, len(componentPorts), strings.Trim(strings.Replace(fmt.Sprint(componentPorts), " ", ",", -1), "[]"))
		case len(componentPorts) == 1:
			return componentPorts[0], nil
		default:
			return portNumber, errors.Errorf("no port is exposed by the component %s", componentName)
		}
	} else {
		for _, port := range componentPorts {
			if portNumber == port {
				return portNumber, nil
			}
		}
	}

	return portNumber, nil
}

// getMachineReadableFormat gives machine readable URL definition
func getMachineReadableFormat(i iextensionsv1.Ingress) Url {
	return Url{
		TypeMeta:   metav1.TypeMeta{Kind: "url", APIVersion: "odo.openshift.io/v1alpha1"},
		ObjectMeta: metav1.ObjectMeta{Name: i.Labels[urlLabels.URLLabel]},
		Spec:       UrlSpec{Host: i.Spec.Backend.ServiceName, Port: i.Spec.Backend.ServicePort.IntValue(), Protocol: getProtocol(i)},
	}

}

func getMachineReadableFormatForList(urls []Url) UrlList {
	return UrlList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "List",
			APIVersion: "odo.openshift.io/v1alpha1",
		},
		ListMeta: metav1.ListMeta{},
		Items:    urls,
	}
}
