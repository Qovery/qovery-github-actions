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
	applicationCommitID  = kingpin.Arg("application-commit-id", "Application commit ID").String()
)

func main() {
	kingpin.Parse()

	envCommitID := ""
	if applicationCommitID == nil || *applicationCommitID == "" {
		envCommitID = os.Getenv("GITHUB_SHA")
		applicationCommitID = &envCommitID
	}

	if applicationCommitID == nil || *applicationCommitID == "" {
		fmt.Println("error: commit ID shouldn't be empty: `application-commit-id` to be set in args or `GITHUB_SHA` env var to be set.")
		os.Exit(1)
	}

	fmt.Printf("Qovery deployment starting for commit: %s ...\n", *applicationCommitID)

	err := app.DeployApplication(*qoveryAPIToken, *qoveryApplicationID, *qoveryEnvironmentID, *applicationCommitID)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
