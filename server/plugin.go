package main

import (
	"fmt"
	"strings"
	"sync"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin struct {
	plugin.MattermostPlugin

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration

	diceRoller *DiceRoller
}

func NewPlugin() *Plugin {
	pluginInstance := &Plugin{}

	pluginInstance.diceRoller = NewDiceRoller()

	return pluginInstance
}

func (p *Plugin) OnActivate() (err error) {
	err = p.API.RegisterCommand(&model.Command{
		Trigger:          "roll",
		DisplayName:      "DiceRoller",
		Description:      "Roll a number of dice using dice algebra.",
		AutoComplete:     true,
		AutoCompleteDesc: "🎲 Roll some dice. See `/roll help` for usage.",
	})
	return
}

func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	if strings.HasPrefix(args.Command, "/roll") {
		return p.ExecuteRoll(c, args)
	} else {
		return nil, nil
	}
}

func (p *Plugin) ExecuteRoll(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {

	commandArgStr := strings.TrimSpace(strings.TrimPrefix(args.Command, "/roll"))

	if strings.HasPrefix(commandArgStr, "help") {
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text: "Usage: `/roll FORMULA...`\n" +
				"Roll at most 10 [dice algebra](https://en.wikipedia.org/wiki/Dice_notation) `FORMULA`(s).\n\n" +
				"A single `FORMULA` has the canonical form `[N]dT[EXPLODE][FILTER...][TOTAL][SUCCESS]` and is evaluated from left to right, where\n" +
				"- `N` is the *optional* number of dice to roll (default: 1)\n" +
				"- `T` is the type of dice to roll:\n" +
				"  - A number: Roll `T`-sided dice (`T >= 2`)\n" +
				"  - `%`: Roll d100 (*percentile*) dice\n" +
				"  - `F`: Roll [Fudge](https://en.wikipedia.org/wiki/Fudge_%28role-playing_game_system%29) dice (equiprobable die outcomes {`plus`, `minus`, `blank`}) and aggregate the total\n" +
				"  - `AE`: Roll Aetherium dice (d12 with outcomes {`switch` on (1-5), `chip` on (6-9), `short` on (10-11), `crash` on (12)} × {`disruption` on (5,9,11,12), `blank` otherwise}) and aggregate the symbols\n" +
				"- `EXPLODE` enables *optional* die explosion:\n" +
				"  - `e>=K`: Roll 1 additional die for each die outcome greater than or equal to `K`\n" +
				"  - `e<=K`: Roll 1 additional die for each die outcome less than or equal to `K`\n" +
				"- Each *optional* `FILTER` (sub)selects the dice used for aggregation:\n" +
				"  - `dlK`: Drops the `K` lowest dice\n" +
				"  - `klK`: Keeps only the `K` lowest dice\n" +
				"  - `dhK`: Drops the `K` highest dice\n" +
				"  - `khK`: Keeps only the `K` highest dice\n" +
				"- `TOTAL` *optionally* sums up the dice outcomes and *optionally* applies a modifier on the result:\n" +
				"  - `=`: Calculates the total without a modifier\n" +
				"  - `+K`: Calculates the total and adds `K`\n" +
				"  - `-K`: Calculates the total and subtracts `K`\n" +
				"  - `*K`: Calculates the total and multiplies it by `K`\n" +
				"  - `/K`: Calculates the total and divides it by `K`\n" +
				"- `SUCCESS` *optionally* determines the number of dice meeting a target number (successes)\n" +
				"  - `s>=K`: Die outcomes greater than or equal to `K` are successes\n" +
				"  - `s<=K`: Die outcomes less than or equal to `K` are successes\n",
			Props: map[string]interface{}{
				"from_webhook": "true",
			},
			Username: "DiceRoller",
		}, nil
	}

	responseText := ""

	commandArgs := strings.Fields(commandArgStr)
	if len(commandArgs) > 10 {
		commandArgs = commandArgs[0:10]
		responseText += fmt.Sprintf("⚠️ Limited to 10 rolls at once.\n")
	}

	responseText += fmt.Sprintf("*throws the dice…* ")

	attachments := []*model.SlackAttachment{}
	for idx := 0; idx < len(commandArgs); idx++ {
		attachments = append(attachments, p.diceRoller.RollNotation(commandArgs[idx]))
	}

	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_IN_CHANNEL,
		Text:         responseText,
		Attachments:  attachments,
		Props: map[string]interface{}{
			"from_webhook": "true",
		},
		Username: "DiceRoller",
	}, nil
}
