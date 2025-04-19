package tools

// The survey package relies on incompatible console management
// function. This wont be used in legacy builds anyway to create the
// config wizard.

import (
	config_proto "www.velocidex.com/golang/velociraptor/config/proto"
	"www.velocidex.com/golang/velociraptor/utils"
)

func GetAPIClientDecryptPassword() (string, error) {
	return "", utils.NotImplementedError
}

func GenerateNewKeys(config_obj *config_proto.Config) error {
	return utils.NotImplementedError
}

func GetAPIClientPassword() (string, error) {
	return "", utils.NotImplementedError
}

func GetInteractiveConfig() (*config_proto.Config, error) {
	return nil, utils.NotImplementedError
}

func StoreServerConfig(config_obj *config_proto.Config) error {
	return utils.NotImplementedError
}

func GenerateFrontendPackages(config_obj *config_proto.Config) error {
	return utils.NotImplementedError
}
