package qovery

import (
	"fmt"
	"github-action/pkg"
)

func GetProjectIdByName(qoveryAPIClient pkg.QoveryAPIClient, orgId string, name string) (string, error) {
	projects, err := qoveryAPIClient.ListProjects(orgId)
	if err != nil {
		return "", err
	}

	for _, project := range projects {
		if project.Name == name {
			return project.ID, nil
		}
	}

	return "", fmt.Errorf("Can't find project with name %v! (it's case sensitive)", name)
}
