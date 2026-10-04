package launcher

import (
	"context"
	"regexp"

	api_proto "www.velocidex.com/golang/velociraptor/api/proto"
	artifacts_proto "www.velocidex.com/golang/velociraptor/artifacts/proto"
	config_proto "www.velocidex.com/golang/velociraptor/config/proto"
	"www.velocidex.com/golang/velociraptor/services"
)

type Suppression struct {
	Name    string
	Subject string

	subjectRegex *regexp.Regexp
}

type AnalysisState struct {
	ArtifactName string
	Suppressions []Suppression
	Errors       []*VerifierError
	Warnings     []*VerifierError
}

func (self *AnalysisState) Done() error {
	return nil
}

func NewAnalysisState(artifact string) *AnalysisState {
	return &AnalysisState{}
}

type VerifierError struct {
	Name, Message string
}

func (self *VerifierError) AsProto() *api_proto.VerifierError {
	res := &api_proto.VerifierError{}
	return res
}

func (self *VerifierError) Error() string {
	return "Error"
}

func (self *AnalysisState) SetError(
	name string, message string, args ...interface{}) {
}

func VerifyArtifact(
	ctx context.Context, config_obj *config_proto.Config,
	repository services.Repository,
	artifact *artifacts_proto.Artifact,
	state *AnalysisState) {
}
