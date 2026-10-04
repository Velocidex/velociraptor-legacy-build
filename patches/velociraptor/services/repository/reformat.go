package repository

import (
	"context"

	"www.velocidex.com/golang/velociraptor/utils"
)

func (self *RepositoryManager) ReformatVQL(
	ctx context.Context, artifact_yaml string) (string, error) {

	return "", utils.NotImplementedError
}
