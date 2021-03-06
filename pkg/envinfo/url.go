package envinfo

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	devfilev1 "github.com/devfile/api/pkg/apis/workspaces/v1alpha2"
	"github.com/devfile/library/pkg/devfile/parser"
	"github.com/openshift/odo/pkg/localConfigProvider"
	"github.com/openshift/odo/pkg/odo/util/validation"
	"github.com/openshift/odo/pkg/util"
	"github.com/pkg/errors"
)

// GetPorts returns the ports, returns default if nil
func (ei *EnvInfo) GetPorts() []string {

	containerComponents := ei.devfileObj.Data.GetDevfileContainerComponents()
	portMap := make(map[string]bool)

	var portList []string
	for _, component := range containerComponents {
		for _, endpoint := range component.Container.Endpoints {
			portMap[strconv.FormatInt(int64(endpoint.TargetPort), 10)] = true
		}
	}

	for port := range portMap {
		portList = append(portList, port)
	}

	sort.Strings(portList)
	return portList
}

// CompleteURL completes the given URL with default values
func (ei *EnvInfo) CompleteURL(url *localConfigProvider.LocalURL) error {
	if len(url.Path) > 0 && (strings.HasPrefix(url.Path, "/") || strings.HasPrefix(url.Path, "\\")) {
		if len(url.Path) <= 1 {
			url.Path = ""
		} else {
			// remove the leading / or \ from provided path
			url.Path = string([]rune(url.Path)[1:])
		}
	}
	// add leading / to path, if the path provided is empty, it will be set to / which is the default value of path
	url.Path = "/" + url.Path

	// get the port if not provided
	if url.Port == -1 {
		var err error
		url.Port, err = util.GetValidPortNumber(ei.GetName(), url.Port, ei.GetPorts())
		if err != nil {
			return err
		}
	}

	// get the name for the URL if not provided
	if len(url.Name) == 0 {
		url.Name = util.GetURLName(ei.GetName(), url.Port)
	}

	containerComponents := ei.devfileObj.Data.GetDevfileContainerComponents()
	if len(containerComponents) == 0 {
		return fmt.Errorf("no valid components found in the devfile")
	}

	// if a container name is provided for the URL, return
	if len(url.Container) > 0 {
		return nil
	}

	containerPortMap := make(map[int]string)
	portMap := make(map[string]bool)

	// if a container name for the URL is not provided
	// use a container which uses the given URL port in one of it's endpoints
	for _, component := range containerComponents {
		for _, endpoint := range component.Container.Endpoints {
			containerPortMap[endpoint.TargetPort] = component.Name
			portMap[strconv.FormatInt(int64(endpoint.TargetPort), 10)] = true
		}
	}
	if containerName, exist := containerPortMap[url.Port]; exist {
		url.Container = containerName
	}

	// container is not provided, or the specified port is not being used under any containers
	// pick the first container to store the new endpoint
	if len(url.Container) == 0 {
		url.Container = containerComponents[0].Name
	}

	return nil
}

// ValidateURL validates the given URL
func (ei *EnvInfo) ValidateURL(url localConfigProvider.LocalURL) error {

	foundContainer := false
	containerComponents := ei.devfileObj.Data.GetDevfileContainerComponents()

	if len(containerComponents) == 0 {
		return fmt.Errorf("no valid components found in the devfile")
	}

	// map TargetPort with containerName
	containerPortMap := make(map[int]string)
	for _, component := range containerComponents {
		if len(url.Container) > 0 && !foundContainer {
			if component.Name == url.Container {
				foundContainer = true
			}
		}
		for _, endpoint := range component.Container.Endpoints {
			if endpoint.Name == url.Name {
				return fmt.Errorf("url %v already exist in devfile endpoint entry under container %v", url.Name, component.Name)
			}
			containerPortMap[endpoint.TargetPort] = component.Name
		}
	}

	if len(url.Container) > 0 && !foundContainer {
		return fmt.Errorf("the container specified: %v does not exist in devfile", url.Container)
	}
	if containerName, exist := containerPortMap[url.Port]; exist {
		if len(url.Container) > 0 && url.Container != containerName {
			return fmt.Errorf("cannot set URL %v under container %v, TargetPort %v is being used under container %v", url.Name, url.Container, url.Port, containerName)
		}
	}

	errorList := make([]string, 0)
	if url.TLSSecret != "" && (url.Kind != localConfigProvider.INGRESS || !url.Secure) {
		errorList = append(errorList, "TLS secret is only available for secure URLs of Ingress kind")
	}

	// check if a host is provided for route based URLs
	if len(url.Host) > 0 {
		if url.Kind == localConfigProvider.ROUTE {
			errorList = append(errorList, "host is not supported for URLs of Route Kind")
		}
		if err := validation.ValidateHost(url.Host); err != nil {
			errorList = append(errorList, err.Error())
		}
	} else if url.Kind == localConfigProvider.INGRESS {
		errorList = append(errorList, "host must be provided in order to create URLS of Ingress Kind")
	}

	// check the protocol of the URL
	if len(url.Protocol) > 0 {
		switch strings.ToLower(url.Protocol) {
		case string(devfilev1.HTTPEndpointProtocol):
			break
		case string(devfilev1.HTTPSEndpointProtocol):
			break
		case string(devfilev1.WSEndpointProtocol):
			break
		case string(devfilev1.WSSEndpointProtocol):
			break
		case string(devfilev1.TCPEndpointProtocol):
			break
		case string(devfilev1.UDPEndpointProtocol):
			break
		default:
			errorList = append(errorList, fmt.Sprintf("endpoint protocol only supports %v|%v|%v|%v|%v|%v", devfilev1.HTTPEndpointProtocol, devfilev1.HTTPSEndpointProtocol, devfilev1.WSSEndpointProtocol, devfilev1.WSEndpointProtocol, devfilev1.TCPEndpointProtocol, devfilev1.UDPEndpointProtocol))
		}
	}

	for _, localURL := range ei.ListURLs() {
		if url.Name == localURL.Name {
			errorList = append(errorList, fmt.Sprintf("URL %s already exists", url.Name))
		}
	}

	if len(errorList) > 0 {
		return fmt.Errorf(strings.Join(errorList, "\n"))
	}
	return nil
}

// GetURL gets the given url from the env.yaml and devfile
func (ei *EnvInfo) GetURL(name string) *localConfigProvider.LocalURL {
	for _, url := range ei.ListURLs() {
		if name == url.Name {
			return &url
		}
	}
	return nil
}

// CreateURL write the given url to the env.yaml and devfile
func (esi *EnvSpecificInfo) CreateURL(url localConfigProvider.LocalURL) error {
	newEndpointEntry := devfilev1.Endpoint{
		Name:       url.Name,
		Path:       url.Path,
		Secure:     url.Secure,
		Exposure:   devfilev1.PublicEndpointExposure,
		TargetPort: url.Port,
		Protocol:   devfilev1.EndpointProtocol(strings.ToLower(url.Protocol)),
	}

	err := addEndpointInDevfile(esi.devfileObj, newEndpointEntry, url.Container)
	if err != nil {
		return errors.Wrapf(err, "failed to write endpoints information into devfile")
	}
	err = esi.SetConfiguration("url", localConfigProvider.LocalURL{Name: url.Name, Host: url.Host, TLSSecret: url.TLSSecret, Kind: url.Kind})
	if err != nil {
		return errors.Wrapf(err, "failed to persist the component settings to env file")
	}
	return nil
}

// ListURLs returns the urls from the env and devfile, returns default if nil
func (ei *EnvInfo) ListURLs() []localConfigProvider.LocalURL {

	envMap := make(map[string]localConfigProvider.LocalURL)
	if ei.componentSettings.URL != nil {
		for _, url := range *ei.componentSettings.URL {
			envMap[url.Name] = url
		}
	}

	var urls []localConfigProvider.LocalURL

	if ei.devfileObj.Data == nil {
		return urls
	}

	for _, comp := range ei.devfileObj.Data.GetDevfileContainerComponents() {
		for _, localEndpoint := range comp.Container.Endpoints {
			// only exposed endpoint will be shown as a URL in `odo url list`
			if localEndpoint.Exposure == devfilev1.NoneEndpointExposure || localEndpoint.Exposure == devfilev1.InternalEndpointExposure {
				continue
			}

			path := "/"
			if localEndpoint.Path != "" {
				path = localEndpoint.Path
			}

			secure := false
			if localEndpoint.Secure || localEndpoint.Protocol == "https" || localEndpoint.Protocol == "wss" {
				secure = true
			}

			url := localConfigProvider.LocalURL{
				Name:      localEndpoint.Name,
				Port:      localEndpoint.TargetPort,
				Secure:    secure,
				Path:      path,
				Container: comp.Name,
			}

			if envInfoURL, exist := envMap[localEndpoint.Name]; exist {
				url.Host = envInfoURL.Host
				url.TLSSecret = envInfoURL.TLSSecret
				url.Kind = envInfoURL.Kind
			} else {
				url.Kind = localConfigProvider.ROUTE
			}

			urls = append(urls, url)
		}
	}

	return urls
}

// DeleteURL is used to delete environment specific info for url from envinfo and devfile
func (esi *EnvSpecificInfo) DeleteURL(name string) error {
	err := removeEndpointInDevfile(esi.devfileObj, name)
	if err != nil {
		return errors.Wrap(err, "failed to delete URL")
	}

	if esi.componentSettings.URL == nil {
		return nil
	}
	for i, url := range *esi.componentSettings.URL {
		if url.Name == name {
			s := *esi.componentSettings.URL
			s = append(s[:i], s[i+1:]...)
			esi.componentSettings.URL = &s
		}
	}
	return esi.writeToFile()
}

// addEndpointInDevfile writes the provided endpoint information into devfile
func addEndpointInDevfile(devObj parser.DevfileObj, endpoint devfilev1.Endpoint, container string) error {
	components := devObj.Data.GetComponents()
	for _, component := range components {
		if component.Container != nil && component.Name == container {
			component.Container.Endpoints = append(component.Container.Endpoints, endpoint)
			devObj.Data.UpdateComponent(component)
			break
		}
	}
	return devObj.WriteYamlDevfile()
}

// removeEndpointInDevfile deletes the specific endpoint information from devfile
func removeEndpointInDevfile(devObj parser.DevfileObj, urlName string) error {
	found := false
	for _, component := range devObj.Data.GetDevfileContainerComponents() {
		for index, enpoint := range component.Container.Endpoints {
			if enpoint.Name == urlName {
				component.Container.Endpoints = append(component.Container.Endpoints[:index], component.Container.Endpoints[index+1:]...)
				devObj.Data.UpdateComponent(component)
				found = true
				break
			}
		}
		if found {
			break
		}
	}
	if !found {
		return fmt.Errorf("the URL %s does not exist", urlName)
	}
	return devObj.WriteYamlDevfile()
}
