package skills

import (
	"fmt"

	"time"

	"github.com/nlopes/slack"
	"github.com/wlbr/mitch/bot"
)

/*
func init() {
	Register(
		"uptime",
		"Reply with the uptime",
		func(conv hanu.ConversationInterface) {
			conv.Reply("Thanks for asking! I'm running since `%s`", time.Since(Start))
		},
	)
}
*/

type UptimeInfo struct {
}

func NewUptimeInfo() *UptimeInfo {
	return &UptimeInfo{}
}

func (u *UptimeInfo) Keyword() string {
	return "uptime"
}

func (up *UptimeInfo) Handle(b *bot.Bot, msg string, ev *slack.MessageEvent) {
	user, _ := b.Client.GetUserInfo(ev.User)
	rt := time.Since(b.Config.Upstart)

	b.Reply(ev, fmt.Sprintf("@%s: Running since `%s`.",
		user.Name, round(rt, time.Second)))
}

/*
	samples := []time.Duration{9.63e6, 1.23456789e9, 1.5e9, 1.4e9, -1.4e9, -1.5e9, 8.91234e9, 34.56789e9, 12345.6789e9}
	format := "% 13s % 13s % 13s % 13s % 13s % 13s % 13s\n"
	fmt.Printf(format, "duration", "ms", "0.5s", "s", "10s", "m", "h")
	for _, d := range samples {
		fmt.Printf(
			format,
			d,
			Round(d, time.Millisecond),
			Round(d, 0.5e9),
			Round(d, time.Second),
			Round(d, 10*time.Second),
			Round(d, time.Minute),
			Round(d, time.Hour),
		)
	}

     duration            ms          0.5s             s           10s             m             h
       9.63ms          10ms            0s            0s            0s            0s            0s
  1.23456789s        1.235s            1s            1s            0s            0s            0s
         1.5s          1.5s          1.5s            2s            0s            0s            0s
         1.4s          1.4s          1.5s            1s            0s            0s            0s
        -1.4s         -1.4s         -1.5s           -1s            0s            0s            0s
        -1.5s         -1.5s         -1.5s           -2s            0s            0s            0s
     8.91234s        8.912s            9s            9s           10s            0s            0s
    34.56789s       34.568s         34.5s           35s           30s          1m0s            0s
3h25m45.6789s  3h25m45.679s    3h25m45.5s      3h25m46s      3h25m50s       3h26m0s        3h0m0s
*/

func round(d, r time.Duration) time.Duration {
	if r <= 0 {
		return d
	}
	neg := d < 0
	if neg {
		d = -d
	}
	if m := d % r; m+m < r {
		d = d - m
	} else {
		d = d + r - m
	}
	if neg {
		return -d
	}
	return d
}
