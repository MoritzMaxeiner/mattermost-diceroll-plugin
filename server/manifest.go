// This file is automatically generated. Do not modify it manually.

package main

import (
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
)

var manifest *model.Manifest

const manifestStr = `
{
  "id": "sh.ucw.mattermost-diceroll-plugin",
  "name": "DiceRoll Plugin",
  "description": "Allows users to roll dice using the /roll command.",
  "homepage_url": "https://github.com/MoritzMaxeiner/mattermost-diceroll-plugin",
  "support_url": "https://github.com/MoritzMaxeiner/mattermost-diceroll-plugin/issues",
  "icon_path": "assets/starter-template-icon.svg",
  "version": "0.2.1",
  "min_server_version": "5.12.0",
  "server": {
    "executables": {
      "linux-amd64": "server/dist/plugin-linux-amd64",
      "darwin-amd64": "server/dist/plugin-darwin-amd64",
      "windows-amd64": "server/dist/plugin-windows-amd64.exe"
    },
    "executable": ""
  },
  "settings_schema": {
    "header": "",
    "footer": "",
    "settings": []
  }
}
`

func init() {
	manifest = model.ManifestFromJson(strings.NewReader(manifestStr))
}
