package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type EnvStatus string
type AppStatus string
type ContStatus string
type DbStatus string

// environments states
const (
	EnvStatusBuilding         = "BUILDING"
	EnvStatusBuildError       = "BUILD_ERROR"
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
	EnvStatusStopped          = "STOPPED"
	EnvStatusStopping         = "STOPPING"
	EnvStatusStopError        = "STOP_ERROR"
	EnvStatusStopQueued       = "STOP_QUEUED"
	EnvStatusRestarted        = "RESTARTED"
	EnvStatusRestartError     = "RESTART_ERROR"
	EnvStatusUnknown          = "UNKNOWN"
)

// application states
const (
	AppStatusBuilding         = "BUILDING"
	AppStatusBuildError       = "BUILD_ERROR"
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
	AppStatusStopped          = "STOPPED"
	AppStatusStopping         = "STOPPING"
	AppStatusStopError        = "STOP_ERROR"
	AppStatusStopQueued       = "STOP_QUEUED"
	AppStatusRestarted        = "RESTARTED"
	AppStatusRestartError     = "RESTART_ERROR"
	AppStatusUnknown          = "UNKNOWN"
)

// container states
const (
	ContStatusBuilding         = "BUILDING"
	ContStatusBuildError       = "BUILD_ERROR"
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
	ContStatusStopped          = "STOPPED"
	ContStatusStopping         = "STOPPING"
	ContStatusStopError        = "STOP_ERROR"
	ContStatusStopQueued       = "STOP_QUEUED"
	ContStatusRestarted        = "RESTARTED"
	ContStatusRestartError     = "RESTART_ERROR"
	ContStatusUnknown          = "UNKNOWN"
)

// database states
const (
	DbStatusBuilding         = "BUILDING"
	DbStatusBuildError       = "BUILD_ERROR"
	DbStatusCanceled         = "CANCELED"
	DbStatusCanceling        = "CANCELING"
	DbStatusDeleted          = "DELETED"
	DbStatusDeleteError      = "DELETE_ERROR"
	DbStatusDeleteQueued     = "DELETE_QUEUED"
	DbStatusDeleting         = "DELETING"
	DbStatusDeployed         = "DEPLOYED"
	DbStatusDeploying        = "DEPLOYING"
	DbStatusDeploymentError  = "DEPLOYMENT_ERROR"
	DbStatusDeploymentQueued = "DEPLOYMENT_QUEUED"
	DbStatusQueued           = "QUEUED"
	DbStatusReady            = "READY"
	DbStatusStopped          = "STOPPED"
	DbStatusStopping         = "STOPPING"
	DbStatusStopError        = "STOP_ERROR"
	DbStatusStopQueued       = "STOP_QUEUED"
	DbStatusRestarted        = "RESTARTED"
	DbtStatusRestartError    = "RESTART_ERROR"
	DbStatusUnknown          = "UNKNOWN"
)

type EnvironmentStatus struct {
	ID                      string    `json:"id"`
	State                   EnvStatus `json:"state"`
	ServiceDeploymentStatus string    `json:"service_deployment_status"`
}
type ApplicationStatus struct {
	ID                      string    `json:"id"`
	State                   AppStatus `json:"state"`
	ServiceDeploymentStatus string    `json:"service_deployment_status"`
}

type ContainerStatus struct {
	ID                      string     `json:"id"`
	State                   ContStatus `json:"state"`
	ServiceDeploymentStatus string     `json:"service_deployment_status"`
}

type DatabaseStatus struct {
	ID                      string   `json:"id"`
	State                   DbStatus `json:"state"`
	ServiceDeploymentStatus string   `json:"service_deployment_status"`
}

func NewUnknownEnvironmentStatus(id string) EnvironmentStatus {
	return EnvironmentStatus{
		ID:                      id,
		State:                   EnvStatusUnknown,
		ServiceDeploymentStatus: "",
	}
}

func NewUnknownApplicationStatus(id string) ApplicationStatus {
	return ApplicationStatus{
		ID:                      id,
		State:                   AppStatusUnknown,
		ServiceDeploymentStatus: "",
	}
}

func NewUnknownContainerStatus(id string) ContainerStatus {
	return ContainerStatus{
		ID:                      id,
		State:                   ContStatusUnknown,
		ServiceDeploymentStatus: "",
	}
}

func NewUnknownDatabaseStatus(id string) DatabaseStatus {
	return DatabaseStatus{
		ID:                      id,
		State:                   DbStatusUnknown,
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
	GetDatabaseStatus(databaseId string) (*DatabaseStatus, error)
	ListOrganizations() ([]Organization, error)
	ListProjects(organizationId string) ([]Project, error)
	ListEnvironments(projectId string) ([]Environment, error)
	ListApplications(environmentId string) ([]Application, error)
	ListContainers(environmentId string) ([]Container, error)
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
		jsonData, err := io.ReadAll(resp.Body)
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
		jsonData, err := io.ReadAll(resp.Body)
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
		jsonData, err := io.ReadAll(resp.Body)
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

func (a qoveryAPIClient) GetDatabaseStatus(databseId string) (*DatabaseStatus, error) {
	req, err := http.NewRequest("GET", a.baseURL+"/database/"+databseId+"/status", nil)
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
		jsonData, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		dbStatus := NewUnknownDatabaseStatus(databseId)
		err = json.Unmarshal(jsonData, &dbStatus)

		if err != nil {
			return nil, err
		}

		return &dbStatus, nil
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
		jsonData, err := io.ReadAll(resp.Body)
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

func (a qoveryAPIClient) ListContainers(environmentId string) ([]Container, error) {
	req, err := http.NewRequest("GET", a.baseURL+"/environment/"+environmentId+"/container", nil)
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
		jsonData, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		res := ContainerResult{}
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
		jsonData, err := io.ReadAll(resp.Body)
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
		jsonData, err := io.ReadAll(resp.Body)
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
		jsonData, err := io.ReadAll(resp.Body)
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
		jsonData, err := io.ReadAll(resp.Body)
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
