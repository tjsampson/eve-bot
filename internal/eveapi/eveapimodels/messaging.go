package eveapimodels

import (
	"fmt"
	"strings"

	"gitlab.unanet.io/devops/eve/pkg/eve"
)

func serviceLabel(svc *eve.DeployService) string {
	if svc.ArtifactName == svc.ServiceName {
		return fmt.Sprintf("\n%s:%s", svc.ServiceName, svc.AvailableVersion)
	} else {
		return fmt.Sprintf("\n%s (%s):%s", svc.ServiceName, svc.ArtifactName, svc.AvailableVersion)
	}
}

func migrationLabel(mig *eve.DeployMigration) string {
	if mig.ArtifactName == mig.DatabaseName {
		return fmt.Sprintf("\n%s:%s", mig.DatabaseName, mig.AvailableVersion)
	} else {
		return fmt.Sprintf("\n%s (%s):%s", mig.DatabaseName, mig.ArtifactName, mig.AvailableVersion)
	}
}

func HeaderMsg(val string) string {
	return fmt.Sprintf("\n*%s*", strings.Title(strings.ToLower(val)))
}

func ServicesResultBlock(svcs eve.DeployServices) string {
	result := ""

	if svcs == nil || len(svcs) == 0 {
		return ""
	}

	for _, svc := range svcs {
		if len(result) == 0 {
			result = serviceLabel(svc)
		} else {
			result = result + serviceLabel(svc)
		}
	}

	return result
}

func MigrationResultBlock(migs eve.DeployMigrations) string {
	result := ""

	if migs == nil || len(migs) == 0 {
		return ""
	}

	for _, mig := range migs {
		if len(result) == 0 {
			result = migrationLabel(mig)
		} else {
			result = result + migrationLabel(mig)
		}
	}

	return result
}

func APIMessages(msgs []string) string {
	infoMsgs := ""
	for _, msg := range msgs {
		if len(infoMsgs) == 0 {
			infoMsgs = "\n- " + msg
		} else {
			infoMsgs = infoMsgs + "\n- " + msg
		}
	}
	if len(infoMsgs) == 0 {
		return ""
	}
	return infoMsgs
}

func EnvironmentNamespaceMsg(deploymentResponsePayload *eve.NSDeploymentPlan) string {
	return fmt.Sprintf("```Namespace: %s\nEnvironment: %s\nCluster: %s```", deploymentResponsePayload.Namespace.Alias, deploymentResponsePayload.EnvironmentName, deploymentResponsePayload.Namespace.ClusterName)
}