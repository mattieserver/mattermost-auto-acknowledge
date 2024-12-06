package main

import (
	"fmt"
	"log/slog"

	"github.com/mattieserver/mattermost-auto-acknowledge/internal/confighelper"
	"github.com/mattieserver/mattermost-auto-acknowledge/internal/httphelper"
)

func acknowledgeMessage(mattermosthttp *httphelper.MattermostClient, teamName string, channelName string, userName string, includeSubMessages bool) {
	slog.Info("Acknowledging all messages")

	teamid, _ := mattermosthttp.GetTeamId(teamName)
	slog.Info(teamid)
	channelid, _ := mattermosthttp.GetChannelId(teamid, channelName)
	slog.Info(channelid)
	userid,_ := mattermosthttp.GetUserId(userName)
	slog.Info(userid)

	posts_ids,_ := mattermosthttp.GetPosts(channelid, userid, includeSubMessages)

	for _, post := range posts_ids {
		mattermosthttp.LikePost(post)
	}

}

func main() {
	slog.Info("Starting")

	conf := confighelper.ReadConfig()
	slog.Info(fmt.Sprintf("Using Mattermost: %s", conf.Mattermost.BaseURL))

	mattermosthttp := httphelper.NewMattermostClient(conf.Mattermost.BaseURL, conf.Mattermost.Username, conf.Mattermost.Password)
	acknowledgeMessage(&mattermosthttp, conf.Mattermost.TeamName, conf.Mattermost.ChannelName, conf.Mattermost.LookupUserName, conf.Mattermost.IncludeSubMessages)

	slog.Info("Done")
}
