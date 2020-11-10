package commands

import (
	"fmt"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/help"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/params"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/resources"
)

type deleteCmd struct {
	baseCommand
}

const (
	// DeleteCmdName used as key/id for the delete command
	DeleteCmdName = "delete"
)

var (
	deleteCmdHelpSummary = help.Summary("The `delete` command is used to delete resource values (metadata)")
	deleteCmdHelpUsage   = help.Usage{
		"delete {{ resources }} for {{ service }} in {{ namespace }} {{ environment }}",
	}
	deleteCmdHelpExample = help.Examples{
		"delete metadata for unaneta in current una-int key",
		"delete metadata for unaneta in current una-int key key2 key3 keyN",
		"delete version for unaneta in current una-int",
	}
)

// NewDeleteCommand creates a New DeleteCmd that implements the EvebotCommand interface
func NewDeleteCommand(cmdFields []string, channel, user string) EvebotCommand {
	cmd := deleteCmd{baseCommand{
		input:  cmdFields,
		info:   ChatInfo{User: user, Channel: channel, CommandName: DeleteCmdName},
		opts:   make(CommandOptions),
		bounds: InputLengthBounds{Min: 7, Max: -1},
	}}
	cmd.resolveDynamicOptions()
	return cmd
}

// AckMsg satisfies the EveBotCommand Interface and returns the acknowledgement message
func (cmd deleteCmd) AckMsg() (string, bool) {
	return cmd.BaseAckMsg(help.New(
		help.HeaderOpt(deleteCmdHelpSummary.String()),
		help.UsageOpt(deleteCmdHelpUsage.String()),
		help.ExamplesOpt(deleteCmdHelpExample.String()),
	).String())
}

// IsAuthorized satisfies the EveBotCommand Interface and checks the auth
func (cmd deleteCmd) IsAuthorized(allowedChannelMap map[string]interface{}, fn chatChannelInfoFn) bool {
	return cmd.IsHelpRequest() || validChannelAuthCheck(cmd.info.Channel, allowedChannelMap, fn) || lowerEnvAuthCheck(cmd.opts)
}

// Options satisfies the EveBotCommand Interface and returns the dynamic options
func (cmd deleteCmd) Options() CommandOptions {
	return cmd.opts
}

// Info satisfies the EveBotCommand Interface and returns the Chat Info
func (cmd deleteCmd) Info() ChatInfo {
	return cmd.info
}

func (cmd *deleteCmd) resolveDynamicOptions() {
	if cmd.ValidInputLength() == false {
		cmd.errs = append(cmd.errs, fmt.Errorf("invalid delete command: %v", cmd.input))
		return
	}

	if resources.IsValidDelete(cmd.input[1]) {
		cmd.opts["resource"] = cmd.input[1]
	} else {
		cmd.errs = append(cmd.errs, fmt.Errorf("invalid delete resource: %v", cmd.input))
		return
	}

	if cmd.opts["resource"] == nil {
		cmd.errs = append(cmd.errs, fmt.Errorf("invalid resource: %v", cmd.input))
		return
	}

	if len(cmd.errs) > 0 {
		return
	}

	switch cmd.opts["resource"] {
	case resources.MetadataName:
		// delete metadata for unaneta in current una-int key key2 key3
		// delete metadata for {{ service }} in {{ namespace }} {{ environment }} key key2 key3
		if cmd.ValidInputLength() == false {
			cmd.errs = append(cmd.errs, fmt.Errorf("invalid delete metadata: %v", cmd.input))
			return
		}

		cmd.opts[params.ServiceName] = cmd.input[3]
		cmd.opts[params.NamespaceName] = cmd.input[5]
		cmd.opts[params.EnvironmentName] = cmd.input[6]
		cmd.opts[params.MetadataName] = cmd.input[7:]
		return
	case resources.VersionName:
		// delete version for unaneta in current una-int
		if cmd.ValidInputLength() == false {
			cmd.errs = append(cmd.errs, fmt.Errorf("invalid delete version: %v", cmd.input))
			return
		}
		cmd.opts[params.ServiceName] = cmd.input[3]
		cmd.opts[params.NamespaceName] = cmd.input[5]
		cmd.opts[params.EnvironmentName] = cmd.input[6]
		return
	default:
		cmd.errs = append(cmd.errs, fmt.Errorf("invalid resource supplied: %v", cmd.opts["resource"]))
		return
	}
}
