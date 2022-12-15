package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type EnvStatus string
type AppStatus string
type ContStatus string

// environments states
const (
	EnvStatusBuilding         = "BUILDING"
	EnvStatusCancelled        = "CANCELED"
	EnvStatusCancelling       = "CANCELING"
	EnvStatusDeleted          = "DELETED"
	EnvStatusDeleteError      = "DELETE_ERROR"
	EnvStatusDeleteQueued     = "DELETE_QUEUED"
	EnvStatusDeleting         = "DELETING"
	EnvStatusDeployed         = "DEPLOYED"
	EnvStatusDeploying        = "DEPLOYING"
	EnvStatusDeploymentError  = "DEPLOYMENT_ERROR"
	EnvStatusDeploymentQueued = "DEPLOYMENT_QUEUED"
	EnvStatusQueued           = "QUEUED"
	EnvStatusReady            = "READY"
	EnvStatusRunning          = "RUNNING"
	EnvStatusStopped          = "STOPPED"
	EnvStatusStopping         = "STOPPING"
	EnvStatusStopError        = "STOP_ERROR"
	EnvStatusStopQueued       = "STOP_QUEUED"
	EnvStatusUnknown          = "UNKNOWN"
)

// application states
const (
	AppStatusBuilding         = "BUILDING"
	AppStatusCanceled         = "CANCELED"
	AppStatusCanceling        = "CANCELING"
	AppStatusDeleted          = "DELETED"
	AppStatusDeleteError      = "DELETE_ERROR"
	AppStatusDeleteQueued     = "DELETE_QUEUED"
	AppStatusDeleting         = "DELETING"
	AppStatusDeployed         = "DEPLOYED"
	AppStatusDeploying        = "DEPLOYING"
	AppStatusDeploymentError  = "DEPLOYMENT_ERROR"
	AppStatusDeploymentQueued = "DEPLOYMENT_QUEUED"
	AppStatusQueued           = "QUEUED"
	AppStatusReady            = "READY"
	AppStatusRunning          = "RUNNING"
	AppStatusStopped          = "STOPPED"
	AppStatusStopping         = "STOPPING"
	AppStatusStopError        = "STOP_ERROR"
	AppStatusStopQueued       = "STOP_QUEUED"
	AppStatusUnknown          = "UNKNOWN"
)

// container states
const (
	ContStatusBuilding         = "BUILDING"
	ContStatusCanceled         = "CANCELED"
	ContStatusCanceling        = "CANCELING"
	ContStatusDeleted          = "DELETED"
	ContStatusDeleteError      = "DELETE_ERROR"
	ContStatusDeleteQueued     = "DELETE_QUEUED"
	ContStatusDeleting         = "DELETING"
	ContStatusDeployed         = "DEPLOYED"
	ContStatusDeploying        = "DEPLOYING"
	ContStatusDeploymentError  = "DEPLOYMENT_ERROR"
	ContStatusDeploymentQueued = "DEPLOYMENT_QUEUED"
	ContStatusQueued           = "QUEUED"
	ContStatusReady            = "READY"
	ContStatusRunning          = "RUNNING"
	ContStatusStopped          = "STOPPED"
	ContStatusStopping         = "STOPPING"
	ContStatusStopError        = "STOP_ERROR"
	ContStatusStopQueued       = "STOP_QUEUED"
	ContStatusUnknown          = "UNKNOWN"
)

type EnvironmentStatus struct {
	ID                      string    `json:"id"`
	State                   EnvStatus `json:"state"`
	Message                 string    `json:"message"`
	ServiceDeploymentStatus string    `json:"service_deployment_status"`
}
type ApplicationStatus struct {
	ID                      string    `json:"id"`
	State                   AppStatus `json:"state"`
	Message                 string    `json:"message"`
	ServiceDeploymentStatus string    `json:"service_deployment_status"`
}

type ContainerStatus struct {
	ID                      string     `json:"id"`
	State                   ContStatus `json:"state"`
	Message                 string     `json:"message"`
	ServiceDeploymentStatus string     `json:"service_deployment_status"`
}

func NewUnknownEnvironmentStatus(id string) EnvironmentStatus {
	return EnvironmentStatus{
		ID:                      id,
		State:                   EnvStatusUnknown,
		Message:                 "",
		ServiceDeploymentStatus: "",
	}
}

func NewUnknownApplicationStatus(id string) ApplicationStatus {
	return ApplicationStatus{
		ID:                      id,
		State:                   AppStatusUnknown,
		Message:                 "",
		ServiceDeploymentStatus: "",
	}
}

func NewUnknownContainerStatus(id string) ContainerStatus {
	return ContainerStatus{
		ID:                      id,
		State:                   ContStatusUnknown,
		Message:                 "",
		ServiceDeploymentStatus: "",
	}
}

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

type QoveryAPIClient interface {
	DeployServices(environmentId string, services ServicesDeployment) error
	DeployDatabase(database Database) error
	GetEnvironmentStatus(environmentId string) (*EnvironmentStatus, error)
	GetApplicationStatus(applicationId string) (*ApplicationStatus, error)
	GetContainerStatus(containerId string) (*ContainerStatus, error)
	ListOrganizations() ([]Organization, error)
	ListProjects(organizationId string) ([]Project, error)
	ListEnvironments(projectId string) ([]Environment, error)
	ListApplications(environmentId string) ([]Application, error)
	ListDatabases(environmentId string) ([]Database, error)
}

type qoveryAPIClient struct {
	c        HTTPClient
	baseURL  string
	apiToken string
	timeout  time.Duration
}

func NewQoveryAPIClient(c HTTPClient, baseURL string, apiToken string, timeout time.Duration) QoveryAPIClient {
	return &qoveryAPIClient{
		c:        c,
		baseURL:  baseURL,
		apiToken: apiToken,
		timeout:  timeout,
	}
}

func (a qoveryAPIClient) DeployServices(environmentId string, services ServicesDeployment) error {
	jsonValue, _ := json.Marshal(services)

	req, err := http.NewRequest("POST", a.baseURL+"/environment/"+environmentId+"/service/deploy", bytes.NewBuffer(jsonValue))

	req.Header.Set("Authorization", "Token "+a.apiToken)
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		return err
	}

	resp, err := a.c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200:
		return nil // deployment launched
	default:
		return fmt.Errorf("qovery API error, status code: %s", resp.Status)
	}
}

func (a qoveryAPIClient) DeployDatabase(database Database) error {
	req, err := http.NewRequest("POST", a.baseURL+"/database/"+database.ID+"/deploy", nil)

	req.Header.Set("Authorization", "Token "+a.apiToken)
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		return err
	}

	resp, err := a.c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200:
		return nil // deployment launched
	default:
		return fmt.Errorf("qovery API error, status code: %s", resp.Status)
	}
}

func (a qoveryAPIClient) GetEnvironmentStatus(environmentId string) (*EnvironmentStatus, error) {
	req, err := http.NewRequest("GET", a.baseURL+"/environment/"+environmentId+"/status", nil)
	req.Header.Set("Authorization", "Token "+a.apiToken)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}

	resp, err := a.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200:
		jsonData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		envStatus := NewUnknownEnvironmentStatus(environmentId)
		err = json.Unmarshal(jsonData, &envStatus)

		if err != nil {
			return nil, err
		}

		return &envStatus, nil
	default:
		return nil, fmt.Errorf("qovery API error, status code: %s", resp.Status)
	}
}

func (a qoveryAPIClient) GetContainerStatus(containerId string) (*ContainerStatus, error) {
	req, err := http.NewRequest("GET", a.baseURL+"/container/"+containerId+"/status", nil)
	req.Header.Set("Authorization", "Token "+a.apiToken)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}

	resp, err := a.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200:
		jsonData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		contStatus := NewUnknownContainerStatus(containerId)
		err = json.Unmarshal(jsonData, &contStatus)

		if err != nil {
			return nil, err
		}

		return &contStatus, nil
	default:
		return nil, fmt.Errorf("qovery API error, status code: %s", resp.Status)
	}
}

func (a qoveryAPIClient) GetApplicationStatus(applicationId string) (*ApplicationStatus, error) {
	req, err := http.NewRequest("GET", a.baseURL+"/application/"+applicationId+"/status", nil)
	req.Header.Set("Authorization", "Token "+a.apiToken)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}

	resp, err := a.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200:
		jsonData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		appStatus := NewUnknownApplicationStatus(applicationId)
		err = json.Unmarshal(jsonData, &appStatus)

		if err != nil {
			return nil, err
		}

		return &appStatus, nil
	default:
		return nil, fmt.Errorf("qovery API error, status code: %s", resp.Status)
	}
}

func (a qoveryAPIClient) ListApplications(environmentId string) ([]Application, error) {
	req, err := http.NewRequest("GET", a.baseURL+"/environment/"+environmentId+"/application", nil)
	req.Header.Set("Authorization", "Token "+a.apiToken)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}

	resp, err := a.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200:
		jsonData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		res := ApplicationResult{}
		err = json.Unmarshal(jsonData, &res)
		if err != nil {
			return nil, err
		}

		return res.Results, nil
	default:
		return nil, fmt.Errorf("qovery API error, status code: %s", resp.Status)
	}
}

func (a qoveryAPIClient) ListDatabases(environmentId string) ([]Database, error) {
	req, err := http.NewRequest("GET", a.baseURL+"/environment/"+environmentId+"/database", nil)
	req.Header.Set("Authorization", "Token "+a.apiToken)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}

	resp, err := a.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200:
		jsonData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		res := DatabaseResult{}
		err = json.Unmarshal(jsonData, &res)
		if err != nil {
			return nil, err
		}

		return res.Results, nil
	default:
		return nil, fmt.Errorf("qovery API error, status code: %s", resp.Status)
	}
}

func (a qoveryAPIClient) ListEnvironments(projectId string) ([]Environment, error) {
	req, err := http.NewRequest("GET", a.baseURL+"/project/"+projectId+"/environment", nil)
	req.Header.Set("Authorization", "Token "+a.apiToken)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}

	resp, err := a.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200:
		jsonData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		res := EnvironmentResult{}
		err = json.Unmarshal(jsonData, &res)
		if err != nil {
			return nil, err
		}

		return res.Results, nil
	default:
		return nil, fmt.Errorf("qovery API error, status code: %s", resp.Status)
	}
}

func (a qoveryAPIClient) ListProjects(organizationId string) ([]Project, error) {
	req, err := http.NewRequest("GET", a.baseURL+"/organization/"+organizationId+"/project", nil)
	req.Header.Set("Authorization", "Token "+a.apiToken)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}

	resp, err := a.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200:
		jsonData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		res := ProjectResult{}
		err = json.Unmarshal(jsonData, &res)
		if err != nil {
			return nil, err
		}

		return res.Results, nil
	default:
		return nil, fmt.Errorf("qovery API error, status code: %s", resp.Status)
	}
}

func (a qoveryAPIClient) ListOrganizations() ([]Organization, error) {
	req, err := http.NewRequest("GET", a.baseURL+"/organization", nil)
	req.Header.Set("Authorization", "Token "+a.apiToken)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}

	resp, err := a.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200:
		jsonData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		res := OrganizationResult{}
		err = json.Unmarshal(jsonData, &res)
		if err != nil {
			return nil, err
		}

		return res.Results, nil
	default:
		return nil, fmt.Errorf("qovery API error, status code: %s", resp.Status)
	}
}
