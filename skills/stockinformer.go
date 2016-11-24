package skills

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strings"

	"github.com/nlopes/slack"
	"github.com/wlbr/mitch/bot"
)

type StockInformer struct {
}

func NewStockInformer() *StockInformer {
	return &StockInformer{}
}

func (si *StockInformer) Keyword() string {
	return "stock"
}

func (si *StockInformer) Handle(b *bot.Bot, msg string, ev *slack.MessageEvent) {
	user, _ := b.Client.GetUserInfo(ev.User)
	cleanmsg := b.StripLinkMarkup(msg)
	b.Reply(ev, fmt.Sprintf("@%s: %s", user.Name, getQuote(cleanmsg)))
}

func (si *StockInformer) Help() string {
	return "`" + si.Keyword() + " <arg>` shows the currents stock price for the stock id `arg`. " +
		"Try `stock AAPL` or `stock UTDI.de`"
}

func getQuote(sym string) string {

	sym = strings.ToUpper(sym)
	url := fmt.Sprintf("http://download.finance.yahoo.com/d/quotes.csv?s=%s&f=nsl1op&e=.csv", sym)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	rows, err := csv.NewReader(resp.Body).ReadAll()
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	if len(rows) >= 1 && len(rows[0]) == 5 {
		return fmt.Sprintf("`%s (%s)` is trading at `$%s`", rows[0][0], rows[0][1], rows[0][2])
	}
	return fmt.Sprintf("unknown response format (symbol was \"%s\")", sym)
}
