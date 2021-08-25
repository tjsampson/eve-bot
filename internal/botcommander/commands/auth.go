package commands

import (
	"context"
	"strings"

	"github.com/unanet/eve-bot/internal/botcommander/params"
	"github.com/unanet/go/pkg/log"
	"go.uber.org/zap"
)

type UserItem struct {
	UserID string
	Email  string
	Name   string
	Roles  []string
	Groups []string
}

// TODO: Implement me
// const userTableName = "eve-bot-users"

// func validUserRoleCheck(cmd EvebotCommand, chatUser *chatmodels.ChatUser, db *dynamodb.DynamoDB) bool {
// 	log.Logger.Info("fetched slack user", zap.Any("user", chatUser), zap.String("cmd", cmd.Info().CommandName))
// 	_, err := db.GetItem(&dynamodb.GetItemInput{
// 		Key: map[string]*dynamodb.AttributeValue{
// 			"UserID": {
// 				S: aws.String(chatUser.FullyQualifiedName()),
// 			},
// 		},
// 		TableName: aws.String(userTableName),
// 	})
// 	if err != nil {
// 		log.Logger.Error("failed to get fully qualified user auth record", zap.Error(err))
// 		return false
// 	}
// 	return true
// }

// validChannelAuthCheck validates/confirm if the incoming channel matches one of the "approved" channels
// approved channels configured via Environment Variable: EVEBOT_SLACK_CHANNELS_AUTH
func validChannelAuthCheck(channel string, channelMap map[string]interface{}, fn ChatChannelInfoFn) bool {
	incomingChannelInfo, err := fn(context.TODO(), channel)
	if err != nil {
		log.Logger.Error("failed to get channel info auth check", zap.Error(err))
		return false
	}

	// Coming from an Elevated/Approved Channel
	// let them pass
	if _, ok := channelMap[incomingChannelInfo.Name]; ok {
		return true
	}
	return false
}

// lowerEnvAuthCheck checks if the incoming command is against a lowe-environment (int/qa/dev)
func lowerEnvAuthCheck(options CommandOptions) bool {
	if options == nil {
		return false
	}

	if env, ok := options[params.EnvironmentName].(string); ok {
		// Let's see if they are performing an action to something in the lower environments (int,qa,dev)
		// Most actions can be taken against resources in the lower environments
		// the only action that can't is the `release` command
		switch {
		case strings.Contains(env, "int"), strings.Contains(env, "qa"), strings.Contains(env, "dev"):
			return true
		}
		return false
	}

	log.Logger.Warn("environment not set")
	return false
}
