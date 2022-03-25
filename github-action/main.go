package main

import (
	"fmt"
	"github-action/pkg"
	"github-action/qovery"
	"net/http"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	kp                  = kingpin.New("Qovery deploy", "A command-line allowing to deploy Qovery application.")
	organizationId      = kingpin.Flag("org-id", "Qovery organization ID").Required().String()
	environmentId       = kingpin.Flag("env-id", "Qovery environment ID").Required().String()
	applicationIds      = kingpin.Flag("app-ids", "Qovery application ID").String()
	applicationCommitId = kingpin.Flag("app-commit-id", "Application commit ID").String()
	databaseId          = kingpin.Flag("db-id", "Qovery database ID").String()
	apiToken            = kingpin.Flag("api-token", "Qovery API token").Required().String()
)

func main() {
	kingpin.Parse()

	envCommitID := ""
	if applicationCommitId == nil || *applicationCommitId == "" {
		envCommitID = os.Getenv("GITHUB_SHA")
		applicationCommitId = &envCommitID
	}

	deployApp := applicationIds != nil && *applicationIds != ""
	deployDb := databaseId != nil && *databaseId != ""

	if deployApp && (applicationCommitId == nil || *applicationCommitId == "") {
		fmt.Println("error: commit ID shouldn't be empty: `app-commit-id` to be set in args or `GITHUB_SHA` env var to be set.")
		os.Exit(1)
	}

	if !deployApp && !deployDb {
		fmt.Println("error: 'app-ids' or 'db-id' property must be defined.")
		os.Exit(1)
	}

	qoveryAPIClient := pkg.NewQoveryAPIClient(
		&http.Client{},
		"https://api.qovery.com",
		*apiToken,
		0,
	)

	var err error = nil
	if deployApp {
		fmt.Printf("Qovery application(s) '%s' deployment starting with commit: %s ...\n", *applicationIds, *applicationCommitId)
		err = qovery.DeployApplication(qoveryAPIClient, *applicationIds, *environmentId, *applicationCommitId)
	} else if deployDb {
		fmt.Printf("Qovery database '%s' deployment starting...\n", *databaseId)
		err = qovery.DeployDatabase(qoveryAPIClient, *databaseId, *environmentId)
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
