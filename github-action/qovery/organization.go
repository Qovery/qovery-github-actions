package qovery

import (
	"errors"
	"fmt"
	"github-action/pkg"
)

func GetOrganizationIdByName(qoveryAPIClient pkg.QoveryAPIClient, name string) (string, error) {
	organizations, err := qoveryAPIClient.ListOrganizations()
	if err != nil {
		return "", err
	}

	for _, org := range organizations {
		if org.Name == name {
			return org.ID, nil
		}
	}

	return "", errors.New(fmt.Sprintf("Can't find organization with name %v! (it's case sensitive)", name))
}
