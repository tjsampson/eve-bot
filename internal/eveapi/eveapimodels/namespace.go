package eveapimodels

import (
	"gitlab.unanet.io/devops/eve/pkg/log"
	"go.uber.org/zap"

	"gitlab.unanet.io/devops/eve/pkg/json"
)

type Namespace struct {
	ID                 int       `json:"id"`
	Name               string    `json:"name"`
	Alias              string    `json:"alias"`
	EnvironmentID      int       `json:"environment_id"`
	EnvironmentName    string    `json:"environment_name"`
	RequestedVersion   string    `json:"requested_version"`
	ExplicitDeployOnly bool      `json:"explicit_deploy_only"`
	ClusterID          int       `json:"cluster_id"`
	Metadata           json.Text `json:"metadata"`
	CreatedAt          string    `json:"created_at"`
	UpdatedAt          string    `json:"updated_at"`
}

type Namespaces []Namespace

func (e Namespaces) ToChatMessage() string {
	if e == nil || len(e) == 0 {
		return "no environments"
	}

	msg := ""

	for _, v := range e {
		log.Logger.Info("WTF:", zap.Any("namespace", e))
		msg += "*Name:*" + v.Alias + "\n" + "*Version:*" + "_" + v.RequestedVersion + "_" + "\n\n"
	}

	return msg
}