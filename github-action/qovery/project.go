package qovery

import (
	"errors"
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

	return "", errors.New(fmt.Sprintf("Can't find project with name %v! (it's case sensitive)", name))
}
