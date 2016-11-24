package skills

import (
	"fmt"

	"github.com/nlopes/slack"
	"github.com/wlbr/mitch/bot"
)

/*
func init() {
	Register(
		"version",
		"Reply with the current chatbot version",
		func(conv hanu.ConversationInterface) {
			conv.Reply("Thanks for asking! I'm running with `%s`", Version)
		},
	)
}
*/

type VersionInfo struct {
}

func NewVersionInfo() *VersionInfo {
	return &VersionInfo{}
}

func (v *VersionInfo) Keyword() string {
	return "version"
}

func (v *VersionInfo) Handle(b *bot.Bot, msg string, ev *slack.MessageEvent) {
	user, _ := b.Client.GetUserInfo(ev.User)
	b.Reply(ev, fmt.Sprintf("@%s: Running `%s` built on `%s`.", user.Name,
		b.Config.GitVersion, b.Config.BuildTimeStamp))
}
