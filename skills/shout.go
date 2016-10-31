package skills

import (
	"strings"

	"github.com/sbstjn/hanu"
)

func init() {
	Register(
		"shout <word>",
		"Reply the passed word in uppercase letters",
		hi,
	)
}

func hi(conv hanu.ConversationInterface) {
	str, _ := conv.String("word")
	conv.Reply(strings.ToUpper(str))
}
