package main

import (
	"fmt"
	"strings"
    "io/ioutil"
	"github.com/bluele/slack"
)

type (
	Repo struct {
		Owner string
		Name  string
	}

	Build struct {
		Event  string
		Number int
		Commit string
		Branch string
		Author string
		Status string
		Link   string
	}

	Config struct {
		Webhook   []string
		Channel   string
		Recipient string
		Username  string
		Template  string
		ImageURL  string
		IconURL   string
		IconEmoji string
		Letter    bool
	}

	Plugin struct {
		Repo   Repo
		Build  Build
		Config Config
	}
)

func (p Plugin) Exec() error {
	attachment := slack.Attachment{
		Text:       message(p.Repo, p.Build),
		Fallback:   fallback(p.Repo, p.Build),
		Color:      color(p.Build),
		MarkdownIn: []string{"text", "fallback"},
		ImageURL:   p.Config.ImageURL,
	}

	payload := slack.WebHookPostPayload{}
	payload.Username = p.Config.Username
	payload.Attachments = []*slack.Attachment{&attachment}
	payload.IconUrl = p.Config.IconURL
	payload.IconEmoji = p.Config.IconEmoji

	if p.Config.Recipient == "" {
		payload.Channel = prepend("#", p.Config.Channel)
	} else {
		payload.Channel = prepend("@", p.Config.Recipient)
	}

	if p.Config.Template != "" {
		txt, err := RenderTrim(p.Config.Template, p)
		if err != nil {
			return err
		}
		attachment.Text = txt
	}

	var letter bool
	letterPayload := slack.WebHookPostPayload{}

	if p.Config.Letter {
	    b, err := ioutil.ReadFile(".Pipeline-Letter")
	    if err == nil {
			letterPayload.Username = p.Config.Username
			letterPayload.Attachments = []*slack.Attachment{&attachment}
			letterPayload.IconUrl = p.Config.IconURL
			letterPayload.IconEmoji = p.Config.IconEmoji
			letterPayload.Text = string(b)
			letter = true
	    }
	}

	for _, webhook := range p.Config.Webhook {

		client := slack.NewWebHook(webhook)
		if letter {
			err := client.PostMessage(&letterPayload)
			if err != nil {
				return err
			}
		}

		err := client.PostMessage(&payload)
		if err != nil {
			return err
		}
	}

	return nil
}

func message(repo Repo, build Build) string {
	return fmt.Sprintf("*%s* <%s|%s/%s#%s> (%s) by %s",
		build.Status,
		build.Link,
		repo.Owner,
		repo.Name,
		build.Commit[:8],
		build.Branch,
		build.Author,
	)
}

func fallback(repo Repo, build Build) string {
	return fmt.Sprintf("%s %s/%s#%s (%s) by %s",
		build.Status,
		repo.Owner,
		repo.Name,
		build.Commit[:8],
		build.Branch,
		build.Author,
	)
}

func color(build Build) string {
	switch build.Status {
	case "success":
		return "good"
	case "failure", "error", "killed":
		return "danger"
	default:
		return "warning"
	}
}

func prepend(prefix, s string) string {
	if !strings.HasPrefix(s, prefix) {
		return prefix + s
	}
	return s
}
