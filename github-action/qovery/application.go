package qovery

import (
	"errors"
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

	return "", errors.New(fmt.Sprintf("Can't find application with name %v! (it's case sensitive)", name))
}
