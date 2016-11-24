package skills

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"log"

	"github.com/nlopes/slack"
	"github.com/spf13/viper"
	"github.com/wlbr/mitch/bot"
)

type Archiver struct {
}

func NewArchiver() *Archiver {
	return &Archiver{}
}

func (a *Archiver) Handle(b *bot.Bot, ev *slack.MessageEvent) {
	name := b.GetMessageAuthor(ev)

	shorttime := decodeTimeStamp(ev.Timestamp).Format("15:04")
	appendToArchive(fmt.Sprintf("%s (%s): %s\n", name, shorttime, ev.Text))
}

func decodeTimeStamp(ts string) time.Time {
	tt := strings.Split(ts, ".")
	if len(tt) > 0 {
		ts = tt[0]
	} else {
		ts = ts
	}
	t, _ := strconv.ParseInt(ts, 10, 64)
	tm := time.Unix(t, 0)
	return tm
}

var mu sync.Mutex

func appendToArchive(m string) {
	mu.Lock()
	defer mu.Unlock()
	fname := viper.GetString("archive")
	f, err := os.OpenFile(fname, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			os.Create(fname)
			f, err = os.OpenFile(fname, os.O_APPEND|os.O_WRONLY, 0644)
		} else {
			log.Println(err)
		}
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	fmt.Fprintf(w, "%s", m)
	w.Flush()
}
