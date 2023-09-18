package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type events struct {
    Events []EV_Day `json:"host_events"`
}

type EV_Day struct {
    Host string `json:"host"`
    Day string `json:"day"`
    Event []string `json:"event"`
}

func get_fluc() ([]EV_Day){
    fluc_events := []EV_Day{}
	tn := time.Now().UTC()
	_, week := tn.ISOWeek()

	coll := colly.NewCollector()
	coll.OnRequest(func(req *colly.Request) {
		// fmt.Println(fmt.Printf("Visiting %s", req.URL))
	})
	coll.OnHTML("li.datum", func(h *colly.HTMLElement) {
        is_weekend := false
        ev_day := ""

        selection := h.DOM
        day := strings.TrimSpace(selection.Find("span.tag").Text())
        info := strings.TrimSpace(selection.Find("ul.info").Find("li.wanne").Text())
        switch day {
            case "Freitag":
                ev_day = "Friday"
                is_weekend = true
            case "Samstag":
                ev_day = "Saturday"
                is_weekend = true
            default:
        }

        if is_weekend && info != ""{
            fluc_ev := EV_Day{
                Host: "Fluc Wanne",
                Day: ev_day,
                Event: strings.Split(info, "\n"),
            }
            fluc_events = append(fluc_events, fluc_ev)
        }
	})
	coll.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error on '%s': %s", r.Request.URL, err.Error())
	})
	coll.Visit(fmt.Sprintf("https://fluc.at/programm/2023_Flucwoche%d.html", week))
    return fluc_events
}

func getWeekendDates() []time.Time {
    var weekendDates []time.Time

    // Aktuelles Datum und Wochentag abrufen
    now := time.Now()
    weekday := now.Weekday()

    // Anzahl der Tage, um zum X zu gelangen
    daysUntilFriday := time.Friday - weekday
    daysUntilSaturday := time.Saturday - weekday
    // fix sunday starting with 0
    daysUntilSunday := time.Sunday + 7 - weekday

    // Datum des nächsten X berechnen
    friday := now.AddDate(0, 0, int(daysUntilFriday))
    saturday := now.AddDate(0, 0, int(daysUntilSaturday))
    sunday := now.AddDate(0, 0, int(daysUntilSunday))

    weekendDates = append(weekendDates, friday, saturday, sunday)

    return weekendDates
}

func get_fish() ([]EV_Day){
    fish_events := []EV_Day{}

    weekendDates := getWeekendDates()

	coll := colly.NewCollector()
	coll.OnRequest(func(req *colly.Request) {
		// fmt.Println(fmt.Printf("Visiting %s", req.URL))
	})
	coll.OnHTML("div.project", func(h *colly.HTMLElement) {
        selection := h.DOM
        title := strings.TrimSpace(selection.Find("h2").Text())
        for _, date := range weekendDates {
            tmp_date := date.Format("02/01")
            if strings.HasPrefix(title, tmp_date) {
                found := false
                for i, ev := range fish_events {
                    if ev.Day == tmp_date {
                        fish_events[i].Event = append(fish_events[i].Event, title)
                        found = true
                        break
                    }
                }
                if !found {
                    fish_ev := EV_Day{
                        Host:  "Grelle Forelle",
                        Day:   date.Weekday().String(),
                        Event: []string{title},
                    }
                    fish_events = append(fish_events, fish_ev)
                }
            }
        }
	})
	coll.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error on '%s': %s", r.Request.URL, err.Error())
	})
	coll.Visit("https://www.grelleforelle.com/programm/")
    return fish_events
}


func main() {
    cur_events := events{}
    cur_events.Events = append(cur_events.Events, get_fluc()...)
    cur_events.Events = append(cur_events.Events, get_fish()...)

	content, err := json.MarshalIndent(cur_events, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	err = os.WriteFile("events.json", content, 0644)
	if err != nil {
        fmt.Println("Error writing file")
	}

    fmt.Println(string(content))
}
