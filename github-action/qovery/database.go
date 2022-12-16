package qovery

import (
	"fmt"
	"github-action/pkg"
)

func GetDatabaseIdByName(qoveryAPIClient pkg.QoveryAPIClient, environmentId string, name string) (string, error) {
	databases, err := qoveryAPIClient.ListDatabases(environmentId)
	if err != nil {
		return "", err
	}

	for _, db := range databases {
		if db.Name == name {
			return db.ID, nil
		}
	}

	return "", fmt.Errorf("can't find database with name %v! (it's case sensitive)", name)
}
