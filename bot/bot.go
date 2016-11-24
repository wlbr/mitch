package bot

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/nlopes/slack"
)

type Config struct {
	SlackToken     string
	ArchiveFile    string
	BuildTimeStamp string
	GitVersion     string
	Upstart        time.Time
}

type AnyHandler interface {
	Handle(b *Bot, i interface{})
}

type AnyMessageHandler interface {
	Handle(b *Bot, ev *slack.MessageEvent)
}

type SkillHandler interface {
	Keyword() string
	Handle(b *Bot, msg string, ev *slack.MessageEvent)
}

type Bot struct {
	MyName          string
	MyId            string
	Config          *Config
	Client          *slack.Client
	Rtm             *slack.RTM
	AnyHandlers     []AnyHandler
	MessageHandlers []AnyMessageHandler
	Skills          []SkillHandler
	skillHandlers   map[string][]SkillHandler
}

func (b *Bot) RegisterAnyHandler(h AnyHandler) {
	b.AnyHandlers = append(b.AnyHandlers, h)
}

func (b *Bot) RegisterMessageHandler(h AnyMessageHandler) {
	b.MessageHandlers = append(b.MessageHandlers, h)
}

func (b *Bot) RegisterSkillHandler(s SkillHandler) {
	b.Skills = append(b.Skills, s)
	if nil == b.skillHandlers {
		b.skillHandlers = make(map[string][]SkillHandler)
	}
	b.skillHandlers[s.Keyword()] = append(b.skillHandlers[s.Keyword()], s)
}

func (b *Bot) HandleAny(msg interface{}) {
	for _, h := range b.AnyHandlers {
		h.Handle(b, msg)
	}
}

func (b *Bot) HandleAnyMessage(ev *slack.MessageEvent) {
	for _, h := range b.MessageHandlers {
		h.Handle(b, ev)
	}
}

func (b *Bot) HandleSkill(ev *slack.MessageEvent) {
	res, msg := b.IsRelevant(ev)
	if res {
		kword := strings.Fields(msg)
		handlers := b.skillHandlers[kword[0]]
		for _, h := range handlers {
			h.Handle(b, strings.Join(kword[1:], " "), ev)
		}
	}
}

func (b *Bot) MainLoop() {
Loop:
	for msg := range b.Rtm.IncomingEvents {
		b.Rtm.GetUsers()
		b.HandleAny(msg)
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
			// Ignore hello

		case *slack.ConnectedEvent:
			//log.Printf("\nInfos: %+v /n", ev.Info)
			// fmt.Println("Connection counter:", ev.ConnectionCount)
			// Replace #general/C2TQVGY5V with your Channel ID
			/* bot.Rtm.SendMessage(bot.Rtm.NewOutgoingMessage(
			fmt.Sprintf("Moin, ich bin wieder da.\nVersion: %s vom %s.",
				bot.Config.GitVersion, bot.Config.BuildTimeStamp), "C2TQVGY5V")) */
			b.MyName = ev.Info.User.Name
			b.MyId = ev.Info.User.ID
			fmt.Println(" done.\nConnected, listening...")

		case *slack.MessageEvent:
			b.HandleAnyMessage(ev)
			b.HandleSkill(ev)

		case *slack.PresenceChangeEvent:
			// fmt.Printf("Presence Change: %v\n", ev)

		case *slack.LatencyReport:
			// fmt.Printf("Current latency: %v\n", ev.Value)

		case *slack.RTMError:
			fmt.Printf("Error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			fmt.Printf("Invalid credentials")
			break Loop

		default:

			// Ignore other events..
			// fmt.Printf("Unexpected: %v\n", msg.Data)
		}
	}
}

// StripLinkMarkup converts <http://google.com|google.com> into google.com etc.
// https://api.slack.com/docs/message-formatting#how_to_display_formatted_messages
func (b *Bot) StripLinkMarkup(m string) string {
	re := regexp.MustCompile("<(.*?)>")
	result := re.FindAllStringSubmatch(m, -1)

	var link string
	for _, c := range result {
		link = c[len(c)-1]

		// Done change Channel, User or Specials tags
		if link[:2] == "#C" || link[:2] == "@U" || link[:1] == "!" {
			continue
		}

		url := link
		if strings.Contains(link, "|") {
			splits := strings.Split(link, "|")
			url = splits[1]
		}

		m = strings.Replace(m, "<"+link+">", url, -1)
	}

	return m
}

func (b *Bot) GetMessageAuthor(ev *slack.MessageEvent) string {
	user, err := b.Client.GetUserInfo(ev.User)
	name := ""
	if err == nil {
		name = user.Name
	} else {
		bot, err := b.Client.GetBotInfo(ev.BotID)
		if err != nil {
			log.Printf("bot.GetMessageAuthor(): could not find user nor bot!")
		} else {
			name = bot.Name
		}
	}
	return name
}

func (b *Bot) DefaultMessageParameters() slack.PostMessageParameters {
	return slack.PostMessageParameters{
		Username:    b.MyName,
		AsUser:      true,
		Parse:       "full",
		LinkNames:   slack.DEFAULT_MESSAGE_LINK_NAMES,
		Attachments: nil,
		UnfurlLinks: slack.DEFAULT_MESSAGE_UNFURL_LINKS,
		UnfurlMedia: slack.DEFAULT_MESSAGE_UNFURL_MEDIA,
		IconURL:     slack.DEFAULT_MESSAGE_ICON_URL,
		IconEmoji:   slack.DEFAULT_MESSAGE_ICON_EMOJI,
		Markdown:    slack.DEFAULT_MESSAGE_MARKDOWN,
		EscapeText:  slack.DEFAULT_MESSAGE_ESCAPE_TEXT,
	}
}

func (b *Bot) Reply(ev *slack.MessageEvent, msg string) {
	b.Rtm.PostMessage(ev.Channel, msg, b.DefaultMessageParameters())
}

/*// IsHelpRequest checks if the user requests the help command
func (m Message) IsHelpRequest() bool {
	return strings.HasSuffix(m.Message, "help") || strings.HasPrefix(m.Message, "help")
} */

// IsDirectMessage checks if the message is received using a direct messaging channel
func (b *Bot) IsDirectChannelMessage(ev *slack.MessageEvent) (bool, string) {
	var result bool

	search := "D"

	if strings.HasPrefix(ev.Channel, search) {
		result = true
	} else {
		result = false
	}
	return result, ev.Text
}

func (b *Bot) IsDirectMessage(ev *slack.MessageEvent) (bool, string) {
	var result bool
	var rest string

	search := "<@" + b.MyId + ">"

	if strings.HasPrefix(ev.Text, search) {
		result = true
		rest = strings.Trim(ev.Text, search)
	} else {
		result = false
		rest = ev.Text
	}
	return result, rest
}

func (b *Bot) IsRelevant(ev *slack.MessageEvent) (bool, string) {
	result, remtext := b.IsDirectMessage(ev)
	if result == false {
		result, remtext = b.IsDirectChannelMessage(ev)
	}

	return result, remtext
}
