package confighelper

import (
	"encoding/json"
	"log/slog"
	"os"
)

type Config struct {
	Mattermost struct {
		BaseURL string `json:"base_url"`
		Username  string `json:"username"`
		Password  string `json:"password"`
		TeamName string `json:"team"`
		ChannelName string `json:"channel"`
		LookupUserName string `json:"user"`
		IncludeSubMessages bool `json:"sub_messages"`
	} `json:"mattermost"`
}

func ReadConfig() Config {
	f, err := os.ReadFile("configs/settings.json")
	if err != nil {
		slog.Error(err.Error())
	}

	var data Config
	json.Unmarshal([]byte(f), &data)

	return data

}
