package docs

import (
	"context"
	"sync"

	api_proto "www.velocidex.com/golang/velociraptor/api/proto"
	config_proto "www.velocidex.com/golang/velociraptor/config/proto"
	"www.velocidex.com/golang/velociraptor/services"
	"www.velocidex.com/golang/velociraptor/utils"
)

type NullDocsManager struct{}

func (self NullDocsManager) Search(
	ctx context.Context,
	query_str string, start, len int) (*api_proto.DocSearchResponses, error) {
	return nil, utils.NotImplementedError
}

func NewDocManager(
	ctx context.Context,
	wg *sync.WaitGroup,
	config_obj *config_proto.Config) (services.DocManager, error) {

	return &NullDocsManager{}, nil
}
