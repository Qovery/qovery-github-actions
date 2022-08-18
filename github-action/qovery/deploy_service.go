package qovery

import (
	"fmt"
	"strings"
	"time"

	"github-action/pkg"
)

func DeployServices(qoveryAPIClient pkg.QoveryAPIClient, environmentId string, services pkg.ServicesDeployment) error {
	timeout := time.Hour * 24 // high timeout we should never reach, API wil timeout before

	// Checking deployment is not QUEUED or DEPLOYING already
	// if so, wait for it to be ready
	stateIsOk := false
	var status *pkg.EnvironmentStatus
	for start := time.Now(); time.Since(start) < timeout; {
		status, err := qoveryAPIClient.GetEnvironmentStatus(environmentId)
		if err != nil {
			fmt.Printf("error while trying to get environment status: %s\n", err)
			time.Sleep(10 * time.Second)
			continue
		}

		// Statuses ok to start a deployment
		if status.State == pkg.EnvStatusDeploymentError ||
			status.State == pkg.EnvStatusStopError ||
			status.State == pkg.EnvStatusRunning ||
			status.State == pkg.EnvStatusRunningError ||
			status.State == pkg.EnvStatusCancelError ||
			status.State == pkg.EnvStatusCancelled ||
			status.State == pkg.EnvStatusUnknown {
			stateIsOk = true
			break
		}

		fmt.Printf("Environment cannot accept deploy yet, state: %s\n", status.State)

		time.Sleep(10 * time.Second)
	}

	// Environment state is not valid even after timeout, cannot deploy the application
	if !stateIsOk {
		return fmt.Errorf("error: services cannot be deployed, environment status is : %s", status.State)
	}

	// Launching deployment
	err := qoveryAPIClient.DeployServices(environmentId, services)
	if err != nil {
		return fmt.Errorf("error while trying to deploy services: %s", err)
	}

	// Waiting for deployment to be OK or ERRORED with a timeout
	for start := time.Now(); time.Since(start) < timeout; {
		status, err := qoveryAPIClient.GetEnvironmentStatus(environmentId)
		if err != nil {
			return fmt.Errorf("error while trying to get environment status: %s", err)
		}

		fmt.Printf("Deployment ongoing: status %s\n", status.State)

		if status.State == pkg.EnvStatusRunning {
			return nil
		} else if strings.HasSuffix(string(status.State), "ERROR") {
			return fmt.Errorf("error: services has not been deployed, environment status is : %s", status.State)
		}

		time.Sleep(10 * time.Second)
	}

	return fmt.Errorf("error: timeout reached, deployment appears to be still ongoing, please check Qovery console.")
}
