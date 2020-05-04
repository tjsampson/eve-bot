package eveapi

import (
	"fmt"

	"gitlab.unanet.io/devops/eve/pkg/eve"
)

type CallbackState struct {
	User    string               `json:"user"`
	Channel string               `json:"channel"`
	Payload eve.NSDeploymentPlan `json:"payload"`
}

//func newBlockMsgOpt(text string) slack.MsgOption {
//	return slack.MsgOptionBlocks(
//		slack.NewSectionBlock(
//			slack.NewTextBlockObject(
//				slack.MarkdownType,
//				text,
//				false,
//				false),
//			nil,
//			nil),
//		slack.NewDividerBlock())
//}

func artifactResultMsg(services eve.DeployServices) string {
	successfulResultsMsg := ""
	successfulResultsHeader := "*Successful:*\n"
	successfulResults := ""
	failedResultsMsg := ""
	failedResultsHeader := "*Failed:*\n"
	failedResults := ""
	noopResultsMsg := ""
	//noopResultsHeader := "*No Change:*\n"
	noopResults := ""
	for _, svc := range services {
		switch svc.Result {
		case eve.DeployArtifactResultFailed:
			if len(failedResults) == 0 {
				failedResults = fmt.Sprintf("`%s:%s`\n", svc.ArtifactName, svc.AvailableVersion)
			} else {
				failedResults = failedResults + fmt.Sprintf("`%s:%s`\n", svc.ArtifactName, svc.AvailableVersion)
			}
		case eve.DeployArtifactResultSucceeded:
			if len(successfulResults) == 0 {
				successfulResults = fmt.Sprintf("`%s:%s`\n", svc.ArtifactName, svc.DeployedVersion)
			} else {
				successfulResults = successfulResults + fmt.Sprintf("`%s:%s`\n", svc.ArtifactName, svc.DeployedVersion)
			}
		case eve.DeployArtifactResultNoop:
			if len(noopResults) == 0 {
				noopResults = fmt.Sprintf("`%s:%s`\n", svc.ArtifactName, svc.AvailableVersion)
			} else {
				noopResults = noopResults + fmt.Sprintf("`%s:%s`\n", svc.ArtifactName, svc.AvailableVersion)
			}
		}
	}

	if len(successfulResults) > 0 {
		successfulResultsMsg = successfulResultsHeader + successfulResults + "\n"
	}

	if len(failedResults) > 0 {
		failedResultsMsg = failedResultsHeader + failedResults + "\n"
	}

	if len(noopResults) > 0 {
		noopResultsMsg = noopResults + "\n"
	}

	return successfulResultsMsg + failedResultsMsg + noopResultsMsg
}

func apiMessages(msgs []string) string {
	infoHeader := "*Info:*\n"
	infoMsgs := ""
	for _, msg := range msgs {
		if len(infoMsgs) == 0 {
			infoMsgs = "```\n- " + msg + "\n"
		} else {
			infoMsgs = infoMsgs + "- " + msg + "\n"
		}
	}
	if len(infoMsgs) == 0 {
		return ""
	}
	return infoHeader + infoMsgs + "```\n"
}

func (cbs *CallbackState) SlackMsgHeader() string {
	switch cbs.Payload.Status {
	case eve.DeploymentPlanStatusComplete:
		return fmt.Sprintf("<@%s>, *%s* deployment in *%s* is complete! Here are your results...\n\n", cbs.User, cbs.Payload.Namespace.Alias, cbs.Payload.EnvironmentName)
	case eve.DeploymentPlanStatusErrors:
		return fmt.Sprintf("<@%s>, we've encountered some errors while deploying *%s* in *%s*! Here are your results...\n\n", cbs.User, cbs.Payload.Namespace.Alias, cbs.Payload.EnvironmentName)
	case eve.DeploymentPlanStatusDryrun:
		return fmt.Sprintf("<@%s>, here's the *%s* `dryrun` results for *%s*...\n\n", cbs.User, cbs.Payload.Namespace.Alias, cbs.Payload.EnvironmentName)
	case eve.DeploymentPlanStatusPending:
		return fmt.Sprintf("<@%s>, your *%s* deployment in *%s* is pending! Here's the plan...\n\n", cbs.User, cbs.Payload.Namespace.Alias, cbs.Payload.EnvironmentName)
	default:
		return ""
	}
}

func (cbs *CallbackState) SlackMsgResults() string {
	return artifactResultMsg(cbs.Payload.Services) + "\n\n" + apiMessages(cbs.Payload.Messages)
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
