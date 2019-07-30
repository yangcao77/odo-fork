package kclient

type KClient interface {
	SetNamespace(namespace string)
	GetCurrentProjectName() string
	SetCurrentProject(projectName string) error
	CreateNewProject(projectName string, wait bool) error
	GetProjectNames() ([]string, error)
	GetProject(projectName string) (interface{}, error)
	GetClusterServiceClass(serviceName string) (interface{}, error)
	GetDeploymentConfigFromName(name string) (interface{}, error)
}

type Client struct {
	namespace string
}

func New(skipConnectionCheck bool) (*Client, error) {
	// TODO: implement for KDO
	client := Client{"fakeNS"}
	return &client, nil
}

func (c *Client) SetNamespace(namespace string) {
	c.namespace = namespace
}

func (c *Client) GetCurrentProjectName() string {
	return "myproject"
}

func (c *Client) SetCurrentProject(projectName string) error {
	c.namespace = projectName
	return nil
}

func (c *Client) CreateNewProject(projectName string, wait bool) error {
	return nil
}

func (c *Client) DeleteProject(name string) error {
	return nil
}

func (c *Client) GetProjectNames() ([]string, error) {
	projectNames := []string{"myproject"}
	return projectNames, nil
}

func (c *Client) GetProject(projectName string) (interface{}, error) {
	return "myproject", nil
}

func (c *Client) GetClusterServiceClass(serviceName string) (interface{}, error) {
	return "clusterServiceClass", nil
}

func (c *Client) GetDeploymentConfigFromName(name string) (interface{}, error) {
	return "clusterServiceClass", nil
}
