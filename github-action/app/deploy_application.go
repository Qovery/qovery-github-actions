package app

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github-action/pkg"
)

func DeployApplication(qoveryAPIToken string, qoveryApplicationIDS string, qoveryEnvironmentID string, applicationCommitID string) error {
	qoveryAPIClient := pkg.NewQoveryAPIClient(
		&http.Client{},
		"https://api.qovery.com",
		qoveryAPIToken,
		0,
	)

	timeout := time.Second * 900 // 15 minutes

	// Checking deployment is not QUEUED or DEPLOYING already
	// if so, wait for it to be ready
	state_is_ok := false
	var status *pkg.EnvironmentStatus
	for start := time.Now(); time.Since(start) < timeout; {
		status, err := qoveryAPIClient.GetEnvironmentStatus(qoveryEnvironmentID)
		if err != nil {
			return fmt.Errorf("error while trying to get environment status: %s", err)
		}

		// Statuses ok to start a deployment
		if status.State == pkg.EnvStatusDeploymentError ||
			status.State == pkg.EnvStatusStopError ||
			status.State == pkg.EnvStatusRunning ||
			status.State == pkg.EnvStatusRunningError ||
			status.State == pkg.EnvStatusCancelError ||
			status.State == pkg.EnvStatusCancelled ||
			status.State == pkg.EnvStatusUnknown {
			state_is_ok = true
			break
		}

		fmt.Printf("Environment cannot accept deploy yet, state: %s\n", status.State)

		time.Sleep(10 * time.Second)
	}

	// Environment state is not valid even after timeout, cannot deploy the application
	if !state_is_ok {
		return fmt.Errorf("error: application cannot be deployed, environment status is : %s", status.State)
	}

	// Launching deployment
	err := qoveryAPIClient.DeployApplications(qoveryEnvironmentID, pkg.Applications{IDS: qoveryApplicationIDS, CommitID: applicationCommitID})
	if err != nil {
		return fmt.Errorf("error while trying to deploy application: %s", err)
	}

	// Waiting for deployment to be OK or ERRORED with a timeout
	for start := time.Now(); time.Since(start) < timeout; {
		status, err := qoveryAPIClient.GetEnvironmentStatus(qoveryEnvironmentID)
		if err != nil {
			return fmt.Errorf("error while trying to get environment status: %s", err)
		}

		fmt.Printf("Deployment ongoing: status %s\n", status.State)

		if status.State == pkg.EnvStatusRunning {
			break
		} else if strings.HasSuffix(string(status.State), "ERROR") {
			return fmt.Errorf("error: application has not been deployed, environment status is : %s", status.State)
		}

		time.Sleep(10 * time.Second)
	}

	return nil
}
