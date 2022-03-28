package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type EnvStatus string

const (
	EnvStatusInitialized      string = "INITIALIZED"
	EnvStatusBuildingQueued          = "BUILDING_QUEUED"
	EnvStatusBuilding                = "BUILDING"
	EnvStatusBuildError              = "BUILD_ERROR"
	EnvStatusBuilt                   = "BUILT"
	EnvStatusDeploymentQueued        = "DEPLOYMENT_QUEUED"
	EnvStatusDeploying               = "DEPLOYING"
	EnvStatusDeploymentError         = "DEPLOYMENT_ERROR"
	EnvStatusDeployed                = "DEPLOYED"
	EnvStatusStopQueued              = "STOP_QUEUED"
	EnvStatusStopping                = "STOPPING"
	EnvStatusStopError               = "STOP_ERROR"
	EnvStatusStopped                 = "STOPPED"
	EnvStatusDeleteQueued            = "DELETE_QUEUED"
	EnvStatusDeleting                = "DELETING"
	EnvStatusDeleteError             = "DELETE_ERROR"
	EnvStatusDeleted                 = "DELETED"
	EnvStatusRunning                 = "RUNNING"
	EnvStatusRunningError            = "RUNNING_ERROR"
	EnvStatusCancelQueued            = "CANCEL_QUEUED"
	EnvStatusCancelling              = "CANCELLING"
	EnvStatusCancelError             = "CANCEL_ERROR"
	EnvStatusCancelled               = "CANCELLED"
	EnvStatusUnknown                 = "UNKNOWN"
)

type EnvironmentStatus struct {
	ID                      string    `json:"id"`
	State                   EnvStatus `json:"state"`
	Message                 string    `json:"message"`
	ServiceDeploymentStatus string    `json:"service_deployment_status"`
}

func NewUnknownEnvironmentStatus(id string) EnvironmentStatus {
	return EnvironmentStatus{
		ID:                      id,
		State:                   EnvStatusUnknown,
		Message:                 "",
		ServiceDeploymentStatus: "",
	}
}

func NewEnvironmentStatus(id string, state EnvStatus, message string, serviceDeploymentStatus string) EnvironmentStatus {
	return EnvironmentStatus{
		ID:                      id,
		State:                   state,
		Message:                 message,
		ServiceDeploymentStatus: serviceDeploymentStatus,
	}
}

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

type QoveryAPIClient interface {
	DeployApplications(environmentId string, applications Applications) error
	DeployDatabase(database Database) error
	GetEnvironmentStatus(environmentID string) (*EnvironmentStatus, error)
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

func (a qoveryAPIClient) DeployApplications(environmentId string, applications Applications) error {
	appIds := strings.Split(applications.IDS, ",")
	values := make([]map[string]string, 0, len(appIds))

	for _, appId := range appIds {
		values = append(values, map[string]string{"git_commit_id": applications.CommitID, "application_id": appId})
	}

	jsonValue, _ := json.Marshal(map[string]interface{}{"applications": values})

	req, err := http.NewRequest("POST", a.baseURL+"/environment/"+environmentId+"/application/deploy", bytes.NewBuffer(jsonValue))

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

func (a qoveryAPIClient) GetEnvironmentStatus(environmentID string) (*EnvironmentStatus, error) {
	req, err := http.NewRequest("GET", a.baseURL+"/environment/"+environmentID+"/status", nil)
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

		envStatus := NewUnknownEnvironmentStatus(environmentID)
		err = json.Unmarshal([]byte(jsonData), &envStatus)

		if err != nil {
			return nil, err
		}

		return &envStatus, nil
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

		result := ApplicationResult{}
		err = json.Unmarshal([]byte(jsonData), &result)
		if err != nil {
			return nil, err
		}

		return result.results, nil
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

		result := DatabaseResult{}
		err = json.Unmarshal([]byte(jsonData), &result)
		if err != nil {
			return nil, err
		}

		return result.results, nil
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

		result := EnvironmentResult{}
		err = json.Unmarshal([]byte(jsonData), &result)
		if err != nil {
			return nil, err
		}

		return result.results, nil
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

		result := ProjectResult{}
		err = json.Unmarshal([]byte(jsonData), &result)
		if err != nil {
			return nil, err
		}

		return result.results, nil
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

		result := OrganizationResult{}
		err = json.Unmarshal([]byte(jsonData), &result)
		if err != nil {
			return nil, err
		}

		return result.results, nil
	default:
		return nil, fmt.Errorf("qovery API error, status code: %s", resp.Status)
	}
}
