package qovery

import (
	"fmt"
	"github-action/pkg"
)

func GetContainerIdByName(qoveryAPIClient pkg.QoveryAPIClient, environmentId string, name string) (string, error) {
	containers, err := qoveryAPIClient.ListContainers(environmentId)
	if err != nil {
		return "", err
	}

	for _, container := range containers {
		if container.Name == name {
			return container.ID, nil
		}
	}

	return "", fmt.Errorf("can't find container with name %v! (it's case sensitive)", name)
}
