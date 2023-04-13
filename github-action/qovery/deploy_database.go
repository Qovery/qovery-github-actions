package qovery

import (
	"fmt"
	"strings"
	"time"

	"github-action/pkg"
)

func DeployDatabase(qoveryAPIClient pkg.QoveryAPIClient, databaseId string, qoveryEnvironmentId string) error {
	timeout := time.Hour * 24 // high timeout we should never reach, API wil timeout before

	// Checking deployment is not QUEUED or DEPLOYING already
	// if so, wait for it to be ready
	stateIsOk := false
	var status *pkg.EnvironmentStatus
	for start := time.Now(); time.Since(start) < timeout; {
		status, err := qoveryAPIClient.GetEnvironmentStatus(qoveryEnvironmentId)
		if err != nil {
			fmt.Printf("error while trying to get environment status: %s\n", err)
			time.Sleep(10 * time.Second)
			continue
		}

		// Statuses ok to start a deployment
		if status.State == pkg.EnvStatusDeploymentError ||
			status.State == pkg.EnvStatusStopError ||
			status.State == pkg.EnvStatusRunning ||
      status.State == pkg.EnvStatusDeployed ||
			status.State == pkg.EnvStatusReady ||
			status.State == pkg.EnvStatusCancelled ||
			status.State == pkg.EnvStatusUnknown {
			stateIsOk = true
			break
		}

		fmt.Printf("Environment cannot accept deploy yet, state: %s\n", status.State)

		time.Sleep(10 * time.Second)
	}

	// Environment state is not valid even after timeout, cannot deploy the database
	if !stateIsOk {
		return fmt.Errorf("error: database cannot be deployed, environment status is : %s", status.State)
	}

	// Launching deployment
	err := qoveryAPIClient.DeployDatabase(pkg.Database{ID: databaseId})
	if err != nil {
		return fmt.Errorf("error while trying to deploy database: %s", err)
	}

	// Waiting for deployment to be OK or ERRORED with a timeout
	lastEnvStatus := pkg.EnvStatusUnknown
	for start := time.Now(); time.Since(start) < timeout; {
		status, err := qoveryAPIClient.GetEnvironmentStatus(qoveryEnvironmentId)
		if err != nil {
			return fmt.Errorf("⚠️ error while trying to get environment status: %s", err)
		}

		fmt.Printf("Deployment ongoing: status %s\n", status.State)
		lastEnvStatus = string(status.State)

		if status.State == pkg.EnvStatusRunning || status.State == pkg.EnvStatusDeployed || strings.HasSuffix(string(status.State), "ERROR") {
			break
		}

		time.Sleep(10 * time.Second)
	}

	fmt.Printf("\n####################################\n")
	fmt.Printf("ENVIRONMENT STATUS: %s\n\n", lastEnvStatus)

	// print database status
	dbStatus, dbErr := qoveryAPIClient.GetDatabaseStatus(databaseId)
	if err != nil {
		return fmt.Errorf("⚠️ Error while trying to get database %s status: %s", databaseId, dbErr)
	}

	dbSuccessFullyDeployed := true
	icon := ""
	if dbStatus.State == pkg.DbStatusRunning || status.State == pkg.DbStatusDeployed {
		icon = "✅"
	} else if strings.HasSuffix(string(dbStatus.State), "ERROR") {
		dbSuccessFullyDeployed = false
		icon = "❌"
	} else {
		dbSuccessFullyDeployed = false
		icon = "❔"
	}
	fmt.Printf("%s Database %s state: %s\n", icon, databaseId, dbStatus.State)
	fmt.Printf("\n####################################")

	if !dbSuccessFullyDeployed {
		return fmt.Errorf("error: database have not been deployed successfully")
	}
	return nil
}
