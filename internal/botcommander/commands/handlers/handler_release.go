package handlers

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/params"
	"gitlab.unanet.io/devops/eve/pkg/eve"
	"gitlab.unanet.io/devops/eve/pkg/log"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/commands"
	"gitlab.unanet.io/devops/eve-bot/internal/chatservice"
	"gitlab.unanet.io/devops/eve-bot/internal/eveapi"
)

type ReleaseHandler struct {
	eveAPIClient eveapi.Client
	chatSvc      chatservice.Provider
}

func NewReleaseHandler(eveAPIClient *eveapi.Client, chatSvc *chatservice.Provider) CommandHandler {
	return ReleaseHandler{
		eveAPIClient: *eveAPIClient,
		chatSvc:      *chatSvc,
	}
}

func (h ReleaseHandler) Handle(ctx context.Context, cmd commands.EvebotCommand, timestamp string) {

	payload := eve.Release{
		Artifact: cmd.APIOptions()[params.ArtifactName].(string),
		Version:  cmd.APIOptions()[params.ArtifactVersionName].(string),
		FromFeed: cmd.APIOptions()[params.FromFeedName].(string),
		ToFeed:   cmd.APIOptions()[params.ToFeedName].(string),
	}

	resp, err := h.eveAPIClient.Release(ctx, payload)
	if err != nil {
		h.chatSvc.UserNotificationThread(ctx, fmt.Sprintf("failed release: %s", err.Error()), cmd.User(), cmd.Channel(), timestamp)
		return
	}

	log.Logger.Debug("release response", zap.String("message", resp.Message))
	h.chatSvc.UserNotificationThread(ctx, message(payload), cmd.User(), cmd.Channel(), timestamp)
}

func message(payload eve.Release) string {
	if len(payload.Version) > 0 && len(payload.ToFeed) > 0 {
		return fmt.Sprintf("successfully released `%s:%s` from `%s` to `%s`", payload.Artifact, payload.Version, payload.FromFeed, payload.ToFeed)
	}
	if len(payload.Version) == 0 && len(payload.ToFeed) == 0 {
		return fmt.Sprintf("successfully released `%s` from `%s`", payload.Artifact, payload.FromFeed)
	}
	if len(payload.Version) > 0 && len(payload.ToFeed) == 0 {
		return fmt.Sprintf("successfully released `%s:%s` from `%s`", payload.Artifact, payload.Version, payload.FromFeed)
	}
	if len(payload.Version) == 0 && len(payload.ToFeed) > 0 {
		return fmt.Sprintf("successfully released `%s` from `%s` to %s", payload.Artifact, payload.FromFeed, payload.ToFeed)
	}
	return fmt.Sprintf("successfully released `%s` from `%s`", payload.Artifact, payload.FromFeed)
}
