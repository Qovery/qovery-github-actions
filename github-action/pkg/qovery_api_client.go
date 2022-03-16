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
	DeployApplication(application Application) error
	GetEnvironmentStatus(environmentID string) (*EnvironmentStatus, error)
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

func (a qoveryAPIClient) DeployApplication(application Application) error {
	values := map[string]string{"git_commit_id": application.CommitID}
	jsonValue, _ := json.Marshal(values)

	req, err := http.NewRequest("POST", a.baseURL+"/application/"+application.ID+"/deploy", bytes.NewBuffer(jsonValue))
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
	case 202:
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
		json.Unmarshal([]byte(jsonData), &envStatus)
		if err != nil {
			return nil, err
		}

		return &envStatus, nil
	default:
		return nil, fmt.Errorf("qovery API error, status code: %s", resp.Status)
	}
}
