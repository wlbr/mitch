package skills

import (
	"fmt"

	"github.com/nlopes/slack"
	"github.com/wlbr/mitch/bot"
)

type Help struct {
}

func NewHelp() *Help {
	return &Help{}
}

func (h *Help) Keyword() string {
	return "help"
}

func (h *Help) Help() string {
	return "`" + h.Keyword() + "` creates this help."
}

func (h *Help) Handle(b *bot.Bot, msg string, ev *slack.MessageEvent) {
	name := b.GetMessageAuthor(ev)
	var res string

	res = "This is the bot `" + b.MyName + "` based on http://github.com/wlbr/mitch."

	if len(b.AnyHandlers)+len(b.Skills) > 0 {
		res = res + "\nGeneral functionality:\n"
	}
	for _, handler := range b.AnyHandlers {
		res = res + "\n" + handler.Help()
	}
	for _, handler := range b.MessageHandlers {
		res = res + "\n" + handler.Help()
	}

	if len(b.Skills) > 0 {
		if len(b.AnyHandlers)+len(b.Skills) > 0 {
			res += "\n"
		}
		res = res + "\nCommands:\n"
	}
	for _, handler := range b.Skills {
		res = res + "\n" + handler.Help()
	}

	b.Reply(ev, fmt.Sprintf("@%s: \n %s", name, res))
}
