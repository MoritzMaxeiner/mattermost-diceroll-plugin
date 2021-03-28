package main

import (
	"fmt"
	"os"

	"github.com/mattermost/mattermost-server/v5/plugin"
)

func main() {
	p := NewPlugin()

	if len(os.Args) == 1 {
		plugin.ClientMain(p)
	} else {
		for idx := 1; idx < len(os.Args); idx++ {
			r := p.diceRoller.RollNotation(os.Args[idx])
			fmt.Println(r.Fields[0].Value)
		}
	}
}
