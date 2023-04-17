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
			status.State == pkg.EnvStatusDeployed ||
			status.State == pkg.EnvStatusReady ||
			status.State == pkg.EnvStatusCancelled ||
			status.State == pkg.EnvStatusRestarted ||
			status.State == pkg.EnvStatusRestartError ||
			status.State == pkg.EnvStatusBuildError ||
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
	lastEnvStatus := pkg.EnvStatusUnknown
	for start := time.Now(); time.Since(start) < timeout; {
		status, err := qoveryAPIClient.GetEnvironmentStatus(environmentId)
		if err != nil {
			return fmt.Errorf("⚠️ Error while trying to get environment status: %s", err)
		}

		fmt.Printf("Deployment ongoing: status %s\n", status.State)
		lastEnvStatus = string(status.State)

		if status.State == pkg.EnvStatusDeployed || strings.HasSuffix(string(status.State), "ERROR") {
			break
		}

		time.Sleep(10 * time.Second)
	}

	fmt.Printf("\n####################################\n")
	fmt.Printf("ENVIRONMENT STATUS: %s\n\n", lastEnvStatus)

	// print application status
	appSuccessFullyDeployed := true
	for _, app := range services.Applications {
		status, err := qoveryAPIClient.GetApplicationStatus(app.ApplicationId)
		if err != nil {
			fmt.Errorf("⚠️ Error while trying to get application %s status: %s", app.ApplicationId, err)
		}

		icon := ""
		if status.State == pkg.AppStatusDeployed {
			icon = "✅"
		} else if strings.HasSuffix(string(status.State), "ERROR") {
			icon = "❌"
			appSuccessFullyDeployed = false
		} else {
			icon = "❔"
		}
		fmt.Printf("%s Application %s state: %s\n", icon, app.ApplicationId, status.State)
	}

	// print container status
	containerSuccessFullyDeployed := true
	for _, cont := range services.Containers {
		status, err := qoveryAPIClient.GetContainerStatus(cont.Id)
		if err != nil {
			fmt.Errorf("⚠️ Error while trying to get container %s status: %s", cont.Id, err)
		}

		icon := ""
		if status.State == pkg.AppStatusDeployed {
			icon = "✅"
		} else if strings.HasSuffix(string(status.State), "ERROR") {
			icon = "❌"
			containerSuccessFullyDeployed = false
		} else {
			icon = "❔"
		}
		fmt.Printf("%s Container %s state: %s\n", icon, cont.Id, status.State)
	}

	fmt.Printf("\n####################################")

	if !appSuccessFullyDeployed || !containerSuccessFullyDeployed {
		return fmt.Errorf("error: some application(s) and/or container(s) have not been deployed successfully")
	}
	return nil
}
