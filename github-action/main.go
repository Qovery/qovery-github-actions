package main

import (
	"fmt"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"

	"github-action/app"
)

var (
	application = kingpin.New("Qovery deploy", "A command-line allowing to deploy Qovery application.")

	qoveryOrganizationID = kingpin.Arg("qovery-org-id", "Qovery organization ID").Required().String()
	qoveryEnvironmentID  = kingpin.Arg("qovery-env-id", "Qovery environment ID").Required().String()
	qoveryApplicationID  = kingpin.Arg("qovery-app-id", "Qovery application ID").Required().String()
	qoveryAPIToken       = kingpin.Arg("qovery-api-token", "Qovery API token").Required().String()
)

func main() {
	kingpin.Parse()

	applicationCommitID := os.Getenv("GITHUB_SHA")
	if applicationCommitID == "" {
		fmt.Println("error: GITHUB_SHA shouldn't be empty")
		os.Exit(1)
	}

	fmt.Printf("Qovery deployment starting for commit: %s ...\n", applicationCommitID)

	err := app.DeployApplication(*qoveryAPIToken, *qoveryApplicationID, *qoveryEnvironmentID, applicationCommitID)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
