package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
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
    events := []EV_Day{}
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
            eve := EV_Day{
                Host: "Fluc Wanne",
                Day: ev_day,
                Event: strings.Split(info, "\n"),
            }
            events = append(events, eve)
        }
	})
	coll.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error on '%s': %s", r.Request.URL, err.Error())
	})
	coll.Visit(fmt.Sprintf("https://fluc.at/programm/2023_Flucwoche%d.html", week))
    return events
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
    events := []EV_Day{}

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
                for i, ev := range events {
                    if ev.Day == date.Weekday().String() {
                        events[i].Event = append(events[i].Event, title)
                        found = true
                        break
                    }
                }
                if !found {
                    eve := EV_Day{
                        Host:  "Grelle Forelle",
                        Day:   date.Weekday().String(),
                        Event: []string{title},
                    }
                    events = append(events, eve)
                }
            }
        }
	})
	coll.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error on '%s': %s", r.Request.URL, err.Error())
	})
	coll.Visit("https://www.grelleforelle.com/programm/")
    return events
}


func get_flex() ([]EV_Day){
    events := []EV_Day{}

    weekendDates := getWeekendDates()

	coll := colly.NewCollector()
	coll.OnRequest(func(req *colly.Request) {
		// fmt.Println(fmt.Printf("Visiting %s", req.URL))
	})

    coll.OnHTML("div.tribe-events-calendar-month__day", func(h *colly.HTMLElement) {
        selection := h.DOM
        for _, date := range weekendDates {
            tmp_date := date.Format("2006-01-02")
            selection.Find(fmt.Sprintf("div#tribe-events-calendar-day-%s", tmp_date)).Each(
                func(_ int, sel_day *goquery.Selection) {
                    sel_day.Find("article.tribe-events-calendar-month__calendar-event").Each(
                        func(_ int, sel_art *goquery.Selection) {
                        time := strings.TrimSpace(sel_art.Find("div.tribe-events-calendar-month__calendar-event-datetime").Text())
                        time = strings.ReplaceAll(time, "\t", "")
                        time = strings.ReplaceAll(time, "\n", "")
                        title := strings.TrimSpace(sel_art.Find("h3.tribe-events-calendar-month__calendar-event-title").Text())
                        found := false
                        for i, ev := range events {
                            if ev.Day == date.Weekday().String() {
                                events[i].Event = append(events[i].Event, fmt.Sprintf("%s %s", time, title))
                                found = true
                                break
                            }
                        }
                        if !found {
                            eve := EV_Day{
                                Host:  "Flex",
                                Day:   date.Weekday().String(),
                                Event: []string{fmt.Sprintf("%s %s", time, title)},
                            }
                            events = append(events, eve)
                        }
                    })

                })
            }
    })
	coll.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error on '%s': %s", r.Request.URL, err.Error())
	})

	coll.Visit("https://flex.at/events/monat/")
    return events
}


func main() {
    cur_events := events{}
    cur_events.Events = append(cur_events.Events, get_fluc()...)
    cur_events.Events = append(cur_events.Events, get_fish()...)
    cur_events.Events = append(cur_events.Events, get_flex()...)

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
