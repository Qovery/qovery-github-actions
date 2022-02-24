package app

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github-action/pkg"
)

func DeployApplication(qoveryAPIToken string, qoveryApplicationID string, qoveryEnvironmentID string, applicationCommitID string) error {

	qoveryAPIClient := pkg.NewQoveryAPIClient(
		&http.Client{},
		"https://api.qovery.com",
		qoveryAPIToken,
		0,
	)

	// Launching deployment
	err := qoveryAPIClient.DeployApplication(pkg.Application{ID: qoveryApplicationID, CommitID: applicationCommitID})
	if err != nil {
		return fmt.Errorf("error while trying to deploy application: %s", err)
	}

	// Waiting for deployment to be OK or ERRORED with a timeout
	timeout := time.Second * 3600
	for start := time.Now(); time.Since(start) < timeout; {
		status, err := qoveryAPIClient.GetEnvironmentStatus(qoveryEnvironmentID)
		if err != nil {
			return fmt.Errorf("error while trying to get environment status: %s", err)
		}

		fmt.Printf("Environment is %s\n", status.State)

		if status.State == pkg.EnvStatusRunning {
			break
		} else if strings.HasSuffix(string(status.State), "ERROR") {
			return fmt.Errorf("error: application has not been deployed, environment status is : %s", status.State)
		}

		time.Sleep(10 * time.Second)
	}

	return nil
}
