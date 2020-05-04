package eveapi

type CallbackState struct {
	User    string `url:"user"`
	Channel string `url:"channel"`
}

type EveParams struct {
	State CallbackState `url:"state"`
}

type APIError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type ArtifactDefinitions []*ArtifactDefinition

type DeploymentPlanOptions struct {
	Artifacts   ArtifactDefinitions `json:"artifacts"`
	ForceDeploy bool                `json:"force_deploy"`
	DryRun      bool                `json:"dry_run"`
	CallbackURL string              `json:"callback_url"`
	Environment string              `json:"environment"`
	Namespaces  []string            `json:"namespaces,omitempty"`
	Messages    []string            `json:"messages,omitempty"`
	Type        string              `json:"type"`
	User        string              `json:"user"`
}

type ArtifactDefinition struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	RequestedVersion string `json:"requested_version,omitempty"`
	AvailableVersion string `json:"available_version"`
	ArtifactoryFeed  string `json:"artifactory_feed"`
	ArtifactoryPath  string `json:"artifactory_path"`
	Matched          bool   `json:"-"`
}
