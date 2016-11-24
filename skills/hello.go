package skills

import (
	"fmt"

	"github.com/nlopes/slack"
	"github.com/wlbr/mitch/bot"
)

type Hello struct {
}

func NewHello() *Hello {
	return &Hello{}
}

func (h *Hello) Keyword() string {
	return "hi"
}

func (h *Hello) Help() string {
	return "`" + h.Keyword() + "` is a kind of HelloWorld."
}

func (h *Hello) Handle(b *bot.Bot, msg string, ev *slack.MessageEvent) {
	name := b.GetMessageAuthor(ev)
	b.Reply(ev, fmt.Sprintf("@%s: Hi yourself!", name))
}
