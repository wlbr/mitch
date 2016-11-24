package skills

import (
	"fmt"

	"github.com/nlopes/slack"
	"github.com/wlbr/mitch/bot"
)

type EchoSkill struct {
}

func NewEchoSkill() *EchoSkill {
	return &EchoSkill{}
}

func (e *EchoSkill) Keyword() string {
	return "echo"
}

func (e *EchoSkill) Help() string {
	return "`" + e.Keyword() + " <arg>*` echos the arguments."
}

func (r *EchoSkill) Handle(b *bot.Bot, msg string, ev *slack.MessageEvent) {
	user, _ := b.Client.GetUserInfo(ev.User)
	b.Reply(ev, fmt.Sprintf("@%s: %s", user.Name, msg))
}
