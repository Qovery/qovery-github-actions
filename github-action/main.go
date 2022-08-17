package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github-action/pkg"
	"github-action/qovery"
	"net/http"
	"os"
	"strings"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	kp                  = kingpin.New("Qovery deploy", "A command-line allowing to deploy Qovery application.")
	organizationId      = kingpin.Flag("org-id", "Qovery organization ID").String()
	organizationName    = kingpin.Flag("org-name", "Qovery organization name").String()
	projectId           = kingpin.Flag("project-id", "Qovery project ID").String()
	projectName         = kingpin.Flag("project-name", "Qovery project name").String()
	environmentId       = kingpin.Flag("env-id", "Qovery environment ID").String()
	environmentName     = kingpin.Flag("env-name", "Qovery environment name").String()
	applicationIds      = kingpin.Flag("app-ids", "Qovery application ID(s)").String()
	applicationNames    = kingpin.Flag("app-names", "Qovery application name(s)").String()
	applicationCommitId = kingpin.Flag("app-commit-id", "Application commit ID").String()
	databaseId          = kingpin.Flag("db-id", "Qovery database ID").String()
	databaseName        = kingpin.Flag("db-name", "Qovery database name").String()
	containerIds        = kingpin.Flag("container-ids", "Qovery container ids separated by ,").String()
	containerImageTags  = kingpin.Flag("container-tags", "Qovery container image tags separated by ,").String()
	apiToken            = kingpin.Flag("api-token", "Qovery API token").Required().String()
)

func getOrganizationId(qoveryAPIClient pkg.QoveryAPIClient, id *string, name *string) (string, error) {
	if id != nil && *id != "" {
		return *id, nil
	}

	if name != nil && *name != "" {
		return qovery.GetOrganizationIdByName(qoveryAPIClient, *name)
	}

	return "", errors.New("'org-id' or 'org-name' property must be defined")
}

func getProjectId(qoveryAPIClient pkg.QoveryAPIClient, orgId string, id *string, name *string) (string, error) {
	if id != nil && *id != "" {
		return *id, nil
	}

	if name != nil && *name != "" {
		return qovery.GetProjectIdByName(qoveryAPIClient, orgId, *name)
	}

	return "", errors.New("'project-id' or 'project-name' property must be defined")
}

func getEnvironmentId(qoveryAPIClient pkg.QoveryAPIClient, projectId string, id *string, name *string) (string, error) {
	if id != nil && *id != "" {
		return *id, nil
	}

	if name != nil && *name != "" {
		return qovery.GetEnvironmentIdByName(qoveryAPIClient, projectId, *name)
	}

	return "", errors.New("'env-id' or 'env-name' property must be defined")
}

func getApplicationIds(qoveryAPIClient pkg.QoveryAPIClient, envId string, id *string, name *string) (string, error) {
	if id != nil && *id != "" {
		return *id, nil
	}

	if name != nil && *name != "" {
		var ids []string
		for _, sName := range strings.Split(*name, ",") {
			id, err := qovery.GetApplicationIdByName(qoveryAPIClient, envId, sName)
			handleError(err)

			ids = append(ids, id)
		}

		return strings.Join(ids, ","), nil
	}

	return "", errors.New("'app-ids' or 'app-names' property must be defined")
}

func getDatabaseId(qoveryAPIClient pkg.QoveryAPIClient, envId string, id *string, name *string) (string, error) {
	if id != nil && *id != "" {
		return *id, nil
	}

	if name != nil && *name != "" {
		return qovery.GetDatabaseIdByName(qoveryAPIClient, envId, *name)
	}

	return "", errors.New("'db-id' or 'db-name' property must be defined")
}

func handleError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	kingpin.Parse()

	envCommitID := ""
	if applicationCommitId == nil || *applicationCommitId == "" {
		envCommitID = os.Getenv("GITHUB_SHA")
		applicationCommitId = &envCommitID
	}

	deployApp := (applicationIds != nil && *applicationIds != "") || (applicationNames != nil && *applicationNames != "")
	deployDb := (databaseId != nil && *databaseId != "") || (databaseName != nil && *databaseName != "")
	deployContainer := containerIds != nil && *containerIds != ""

	if deployApp && (applicationCommitId == nil || *applicationCommitId == "") {
		fmt.Println("error: commit ID shouldn't be empty: `app-commit-id` to be set in args or `GITHUB_SHA` env var to be set.")
		os.Exit(1)
	}

	if deployContainer && (containerImageTags == nil || *containerImageTags == "") {
		fmt.Println("error: container-tag shouldn't be empty if you want to deploy a specific container")
		os.Exit(1)
	}

	if !deployApp && !deployDb && !deployContainer {
		fmt.Println("error: 'app-ids' or 'app-names' or 'db-id' or 'db-name' or 'container-ids' property must be defined.")
		os.Exit(1)
	}

	qoveryAPIClient := pkg.NewQoveryAPIClient(
		&http.Client{},
		"https://api.qovery.com",
		*apiToken,
		0,
	)

	organizationId, err := getOrganizationId(qoveryAPIClient, organizationId, organizationName)
	handleError(err)

	projectId, err := getProjectId(qoveryAPIClient, organizationId, projectId, projectName)
	handleError(err)

	environmentId, err := getEnvironmentId(qoveryAPIClient, projectId, environmentId, environmentName)
	handleError(err)

	if deployDb {
		databaseId, err := getDatabaseId(qoveryAPIClient, environmentId, databaseId, databaseName)
		handleError(err)

		fmt.Printf("Qovery database '%s' deployment starting...\n", databaseId)
		err = qovery.DeployDatabase(qoveryAPIClient, databaseId, environmentId)
		handleError(err)
		os.Exit(0)
	}

	if deployApp {
		appsIds, err := getApplicationIds(qoveryAPIClient, environmentId, applicationIds, applicationNames)
		handleError(err)
		applicationIds = &appsIds
	}

	ids := strings.Split(*applicationIds, ",")
	apps := make([]pkg.ApplicationDeployment, 0)
	for _, id := range ids {
		if id == "" {
			continue
		}
		apps = append(apps, pkg.ApplicationDeployment{
			ApplicationId: id,
			GitCommitId:   *applicationCommitId,
		})
	}

	ids = strings.Split(*containerIds, ",")
	tags := strings.Split(*containerImageTags, ",")
	if len(ids) != len(tags) {
		fmt.Println("You don't have the same number of container Ids and image tags.")
		os.Exit(1)
	}

	containers := make([]pkg.ContainerDeployment, 0)
	for ix, id := range ids {
		if id == "" {
			continue
		}
		containers = append(containers, pkg.ContainerDeployment{
			Id:       id,
			ImageTag: tags[ix],
		})
	}
	services := pkg.ServicesDeployment{
		Applications: apps,
		Containers:   containers,
	}

	payload, _ := json.Marshal(services)
	fmt.Printf("Qovery service deployment starting...\n%s\n", payload)
	err = qovery.DeployServices(qoveryAPIClient, environmentId, services)
	handleError(err)
}
