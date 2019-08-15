package skills

import (
	"fmt"

	"net/url"

	"log"
	"net/http"

	"encoding/json"

	"github.com/nlopes/slack"
	"github.com/wlbr/mitch/bot"
)

type WeatherIn struct {
}

func NewWeatherIn() *WeatherIn {
	return &WeatherIn{}
}

func (w *WeatherIn) Keyword() string {
	return "weatherin"
}

func (w *WeatherIn) Help() string {
	return "`" + w.Keyword() + " <arg>` shows the weather forecast in city `arg`. " +
		"Try 'weatherin honolulu' or 'weatherin würzburg'"
}

func (w *WeatherIn) Handle(b *bot.Bot, msg string, ev *slack.MessageEvent) {
	name := b.GetMessageAuthor(ev)
	city := weatherGetCoordinates(b, msg)

	b.Reply(ev, fmt.Sprintf("@%s: weather in `%s` will be `%+v`", name, city.name, city))
}

type WeatherCity struct {
	search        string
	name          string
	latitude      string
	longitude     string
	weather       string
	population    int
	syspopulation int
	temperature   float64
	rain          float64
	wind          float64
}

type weatherResponse struct {
	City struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Coord struct {
			Lon float64 `json:"lon"`
			Lat float64 `json:"lat"`
		} `json:"coord"`
		Country    string `json:"country"`
		Population int    `json:"population"`
		Sys        struct {
			Population int `json:"population"`
		} `json:"sys"`
	} `json:"city"`
	Cod     string  `json:"cod"`
	Message float64 `json:"message"`
	Cnt     int     `json:"cnt"`
	List    []struct {
		Dt   int `json:"dt"`
		Main struct {
			Temp      float64 `json:"temp"`
			TempMin   float64 `json:"temp_min"`
			TempMax   float64 `json:"temp_max"`
			Pressure  float64 `json:"pressure"`
			SeaLevel  float64 `json:"sea_level"`
			GrndLevel float64 `json:"grnd_level"`
			Humidity  int     `json:"humidity"`
			TempKf    float64 `json:"temp_kf"`
		} `json:"main"`
		Weather []struct {
			ID          int    `json:"id"`
			Main        string `json:"main"`
			Description string `json:"description"`
			Icon        string `json:"icon"`
		} `json:"weather"`
		Clouds struct {
			All int `json:"all"`
		} `json:"clouds"`
		Wind struct {
			Speed float64 `json:"speed"`
			Deg   float64 `json:"deg"`
		} `json:"wind"`
		Rain struct {
			ThreeH float64 `son:"3h"`
		} `json:"rain"`
		Sys struct {
			Pod string `json:"pod"`
		} `json:"sys"`
		DtTxt string `json:"dt_txt"`
	} `json:"list"`
}

func weatherGetCoordinates(b *bot.Bot, city string) *WeatherCity {
	c := url.QueryEscape(city)

	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/forecast?q=%s&mode=json&appid=%s", c, b.Config.OpenWeatherMapToken)
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
	var record weatherResponse

	// Use json.Decode for reading streams of JSON data
	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		log.Println(err)
	}

	co := &WeatherCity{}
	co.search = city
	co.name = fmt.Sprintf("%s, %s", record.City.Name, record.City.Country)
	co.latitude = fmt.Sprintf("%f", record.City.Coord.Lat)
	co.longitude = fmt.Sprintf("%f", record.City.Coord.Lon)
	co.syspopulation = record.City.Sys.Population
	co.population = record.City.Population
	co.temperature = record.List[0].Main.Temp
	co.wind = record.List[0].Wind.Speed
	co.rain = record.List[0].Rain.ThreeH

	return co
}
