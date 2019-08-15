package skills

import (
	"fmt"

	"time"

	"github.com/nlopes/slack"
	"github.com/wlbr/mitch/bot"
)

type VersionInfo struct {
}

func NewVersionInfo() *VersionInfo {
	return &VersionInfo{}
}

func (v *VersionInfo) Keyword() string {
	return "version"
}

func (v *VersionInfo) Help() string {
	return "`" + v.Keyword() + "` version and build timestamp of the bot."
}

func (v *VersionInfo) Handle(b *bot.Bot, msg string, ev *slack.MessageEvent) {
	user, _ := b.Client.GetUserInfo(ev.User)
	b.Reply(ev, fmt.Sprintf("@%s: Running `%s` built on `%s`.", user.Name,
		b.Config.GitVersion, b.Config.BuildTimeStamp.Format(time.RFC822)))
}
