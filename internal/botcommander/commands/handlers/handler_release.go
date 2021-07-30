package handlers

import (
	"context"
	"fmt"

	"github.com/unanet/eve-bot/internal/botcommander/interfaces"

	"github.com/unanet/eve-bot/internal/botcommander/commands"
	"github.com/unanet/eve-bot/internal/botcommander/params"
	"github.com/unanet/eve-bot/internal/eveapi"
	"github.com/unanet/eve/pkg/eve"
)

// ReleaseHandler is the handler for the ReleaseCmd
type ReleaseHandler struct {
	eveAPIClient interfaces.EveAPI
	chatSvc      interfaces.ChatProvider
}

// NewReleaseHandler creates a ReleaseHandler
func NewReleaseHandler(eveAPIClient interfaces.EveAPI, chatSvc interfaces.ChatProvider) CommandHandler {
	return ReleaseHandler{
		eveAPIClient: eveAPIClient,
		chatSvc:      chatSvc,
	}
}

// Handle handles the ReleaseCmd
func (h ReleaseHandler) Handle(ctx context.Context, cmd commands.EvebotCommand, timestamp string) {

	dynamicOpts := cmd.Options()

	release, err := h.eveAPIClient.Release(ctx, eve.Release{
		Artifact: dynamicOpts[params.ArtifactName].(string),
		Version:  dynamicOpts[params.ArtifactVersionName].(string),
		FromFeed: dynamicOpts[params.FromFeedName].(string),
		ToFeed:   dynamicOpts[params.ToFeedName].(string),
	})
	if err != nil {
		h.chatSvc.UserNotificationThread(ctx, fmt.Sprintf("failed release: %s", err.Error()), cmd.Info().User, cmd.Info().Channel, timestamp)
		return
	}

	h.chatSvc.ReleaseResultsMessageThread(ctx, eveapi.ChatMessage(release), cmd.Info().User, cmd.Info().Channel, timestamp)
}
