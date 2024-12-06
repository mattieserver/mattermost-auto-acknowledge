package httphelper

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	mattermost "github.com/mattermost/mattermost/server/public/model"
)

type MattermostClient struct {
	client *mattermost.Client4
	userid string
}

func NewMattermostClient(baseurl string, username string, password string) MattermostClient {
	client := mattermost.NewAPIv4Client(baseurl)
	client.Login(context.Background(), username, password)

	e := MattermostClient{client, ""}
	e.userid, _ = e.GetOwnUserId()
	return e
}

func (e *MattermostClient) GetOwnUserId() (string, error) {
	etag := ""
	user_id, _, err := e.client.GetUser(context.Background(), "me", etag)
	if err != nil {
		return "", fmt.Errorf("error gettting teams: %s", err)
	}
	return user_id.Id, nil
}

func (e *MattermostClient) GetTeamId(teamName string) (string, error) {
	etag := ""
	teams_id, _, err := e.client.GetTeamByName(context.Background(), teamName, etag)
	if err != nil {
		return "", fmt.Errorf("error gettting teams: %s", err)
	}
	return teams_id.Id, nil
}
func (e *MattermostClient) GetChannelId(teamId string, channelName string) (string, error) {
	etag := ""
	channel_id, _, err := e.client.GetChannelByName(context.Background(), channelName, teamId, etag)
	if err != nil {
		return "", fmt.Errorf("error gettting teams: %s", err)
	}
	return channel_id.Id, nil
}

func (e *MattermostClient) GetUserId(username string) (string, error) {
	etag := ""
	user_id, _, err := e.client.GetUserByUsername(context.Background(), username, etag)
	if err != nil {
		return "", fmt.Errorf("error gettting teams: %s", err)
	}
	return user_id.Id, nil
}

func (e *MattermostClient) GetPosts(channelId string, userid string, includeSubMessages bool) ([]string, error) {
	etag := ""
	reachedAll := false
	i := 0
	rootPostsByUser := []string{}
	secondaryPostsByUser := make(map[string]string)

	postToReactTo := []string{}

	for !reachedAll {
		posts, _, err := e.client.GetPostsForChannel(context.Background(), channelId, i, 30, etag, false, false)
		if err != nil {
			return []string{}, fmt.Errorf("error gettting teams: %s", err)
		}

		//hasnext is broken
		if posts.PrevPostId == "" && posts.NextPostId == "" {
			reachedAll = true
		}

		for _, post := range posts.Posts {
			if post.UserId == userid {

				if strings.Contains(post.Message,"added to the channel by") {
					continue
				}

				if post.RootId == "" {
					rootPostsByUser = append(rootPostsByUser, post.Id)
					postToReactTo = append(postToReactTo, post.Id)
				} else {
					secondaryPostsByUser[post.Id] = post.RootId
				}
			}
		}

		i += 1
	}

	slog.Info(fmt.Sprintf("Looping %d times to get all posts", i))

	if includeSubMessages {
		for secondaryPostId, secondarypostRooId := range secondaryPostsByUser {
			if contains(rootPostsByUser, secondarypostRooId) {
				postToReactTo = append(postToReactTo, secondaryPostId)
			}
		}
	}	

	return postToReactTo, nil
}

func (e *MattermostClient) LikePost(postid string) {
	reacted := false
	reactions, _, err := e.client.GetReactions(context.Background(), postid)
	if err != nil {
		slog.Error(err.Error())
	}
	for _, reaction := range reactions {
		if reaction.UserId == e.userid {
			reacted = true
			break
		}
	}
	if !reacted {
		slog.Info("Reacting to post")
		reaction := mattermost.Reaction{
			UserId:    e.userid,
			PostId:    postid,
			EmojiName: "+1",
		}
		e.client.SaveReaction(context.Background(), &reaction)
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
