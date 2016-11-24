package skills

import (
	"fmt"

	"net/url"

	"log"
	"net/http"

	"encoding/json"

	"time"

	"github.com/nlopes/slack"
	"github.com/wlbr/mitch/bot"
)

type TimeIn struct {
}

func NewTimeIn() *TimeIn {
	return &TimeIn{}
}

func (t *TimeIn) Keyword() string {
	return "timein"
}

func (t *TimeIn) Handle(b *bot.Bot, msg string, ev *slack.MessageEvent) {
	name := b.GetMessageAuthor(ev)
	city := getCoordinates(msg)
	city = getTimeZone(city)
	now := time.Now()
	tzthere, err := time.LoadLocation(city.timezoneid)
	if err != nil {
		log.Fatal("Time Conversion: ", err)
		return
	}
	b.Reply(ev, fmt.Sprintf("@%s: current time in `%s` is `%s`", name, city.name, now.In(tzthere).Format(time.RFC822)))
}

type City struct {
	search       string
	name         string
	latitude     string
	longitude    string
	timezonename string
	timezoneid   string
}

type googleCityReply struct {
	Results []struct {
		AddressComponents []struct {
			LongName  string   `json:"long_name"`
			ShortName string   `json:"short_name"`
			Types     []string `json:"types"`
		} `json:"address_components"`
		FormattedAddress string `json:"formatted_address"`
		Geometry         struct {
			Bounds struct {
				Northeast struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"northeast"`
				Southwest struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"southwest"`
			} `json:"bounds"`
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
			LocationType string `json:"location_type"`
			Viewport     struct {
				Northeast struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"northeast"`
				Southwest struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"southwest"`
			} `json:"viewport"`
		} `json:"geometry"`
		PlaceID string   `json:"place_id"`
		Types   []string `json:"types"`
	} `json:"results"`
	Status string `json:"status"`
}

func getCoordinates(city string) *City {
	c := url.QueryEscape(city)

	url := fmt.Sprintf("http://maps.googleapis.com/maps/api/geocode/json?address=%s&sensor=false", c)
	// Build the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return nil
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return nil
	}

	defer resp.Body.Close()
	// Fill the record with the data from the JSON
	var record googleCityReply

	// Use json.Decode for reading streams of JSON data
	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		log.Println(err)
	}

	co := &City{}
	co.search = city
	co.name = record.Results[0].AddressComponents[0].LongName
	co.latitude = fmt.Sprintf("%f", record.Results[0].Geometry.Location.Lat)
	co.longitude = fmt.Sprintf("%f", record.Results[0].Geometry.Location.Lng)

	return co
}

type googleTimezoneReply struct {
	DstOffset    int    `json:"dstOffset"`
	RawOffset    int    `json:"rawOffset"`
	Status       string `json:"status"`
	TimeZoneID   string `json:"timeZoneId"`
	TimeZoneName string `json:"timeZoneName"`
}

func getTimeZone(c *City) *City {

	url := fmt.Sprintf("https://maps.googleapis.com/maps/api/timezone/json?location=%s,%s&timestamp=1331161200&sensor=false", c.latitude, c.longitude)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return nil
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return nil
	}

	defer resp.Body.Close()
	var record googleTimezoneReply
	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		log.Println(err)
	}
	c.timezonename = record.TimeZoneName
	c.timezoneid = record.TimeZoneID
	return c
}
