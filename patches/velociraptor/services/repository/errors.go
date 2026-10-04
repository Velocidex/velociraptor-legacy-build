package repository

import (
	artifacts_proto "www.velocidex.com/golang/velociraptor/artifacts/proto"
)

func reportError(errToReport error,
	artifact *artifacts_proto.Artifact, field string, idx int) error {

	return errToReport
}
