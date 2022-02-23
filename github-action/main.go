package main

import (
	"fmt"
	"os"

	"github.com/Qovery/qovery-github-action/app"
)

func main() {
	applicationCommitID := os.Getenv("GITHUB_SHA")
	// qoveryOrganizationID := os.Getenv("QOVERY_ORGANIZATION_ID")
	qoveryEnvironmentID := os.Getenv("INPUT_QOVERY_ENVIRONMENT_ID")
	qoveryApplicationID := os.Getenv("INPUT_QOVERY_APPLICATION_ID")
	qoveryAPIToken := os.Getenv("QOVERY_API_TOKEN")

	fmt.Printf("Qovery deployment starting for commit: %s ...\n", applicationCommitID)

	err := app.DeployApplication(qoveryAPIToken, qoveryApplicationID, qoveryEnvironmentID, applicationCommitID)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
