package qovery

import (
	"fmt"
	"github-action/pkg"
)

func GetApplicationIdByName(qoveryAPIClient pkg.QoveryAPIClient, environmentId string, name string) (string, error) {
	applications, err := qoveryAPIClient.ListApplications(environmentId)
	if err != nil {
		return "", err
	}

	for _, app := range applications {
		if app.Name == name {
			return app.ID, nil
		}
	}

	return "", fmt.Errorf("can't find application with name %v! (it's case sensitive)", name)
}
