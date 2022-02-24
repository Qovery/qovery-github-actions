package main

import (
	"fmt"
	"os"

	"github-action/app"
)

func main() {
	applicationCommitID := os.Getenv("GITHUB_SHA")
	// qoveryOrganizationID := os.Args[1]
	qoveryEnvironmentID := os.Args[2]
	qoveryApplicationID := os.Args[3]
	qoveryAPIToken := os.Args[4]

	fmt.Printf("Qovery deployment starting for commit: %s ...\n", applicationCommitID)

	err := app.DeployApplication(qoveryAPIToken, qoveryApplicationID, qoveryEnvironmentID, applicationCommitID)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
