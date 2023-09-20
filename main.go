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

var weekendDates []time.Time

func add_event_info(events []EV_Day, host string, day string, event_info []string) ([]EV_Day){
    found := false
    for i, ev := range events {
        if ev.Day == day {
            events[i].Event = append(events[i].Event, event_info...)
            found = true
            break
        }
    }
    if !found {
        eve := EV_Day{
            Host:  host,
            Day:   day,
            Event: event_info,
        }
        events = append(events, eve)
    }
    return events
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
            events = add_event_info(events, "Fluc Wanne", ev_day, strings.Split(info, "\n"))
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
                events = add_event_info(events, "Grelle Forelle", date.Weekday().String(), []string{title})
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
                        events = add_event_info(events, "Flex", date.Weekday().String(),
                            []string{fmt.Sprintf("%s: %s", time, title)})
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


func get_exil() ([]EV_Day){
    events := []EV_Day{}

	coll := colly.NewCollector()
	coll.OnRequest(func(req *colly.Request) {
		// fmt.Println(fmt.Printf("Visiting %s", req.URL))
	})

    coll.OnHTML("tr.container", func(h *colly.HTMLElement) {
        selection := h.DOM
        title := strings.TrimSpace(selection.Find("h3").Text())
        day := selection.Find("span:not([class])").First().Text()
        time := selection.Find("span:not([class])").Eq(1).Text()
        for _, date := range weekendDates {
            tmp_date := date.Format("02/01/2006")
            if strings.Contains(day, tmp_date){
                events = add_event_info(events, "Exil", date.Weekday().String(),
                    []string{fmt.Sprintf("%s: %s", time, title)})
            }
        }
    })

	coll.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error on '%s': %s", r.Request.URL, err.Error())
	})

	coll.Visit("https://exil1.ticket.io/")
    return events
}

func get_werk() ([]EV_Day){
    events := []EV_Day{}

	coll := colly.NewCollector()
	coll.OnRequest(func(req *colly.Request) {
		// fmt.Println(fmt.Printf("Visiting %s", req.URL))
	})

    coll.OnHTML("div.events--preview-item", func(h *colly.HTMLElement) {
        ev_day := ""
        selection := h.DOM
        title := strings.TrimSpace(selection.Find(".preview-item--headline").Text())
        day := selection.Find("ul.preview-item--information").Find("li:not([class])").First().Text()
        time := selection.Find("ul.preview-item--information").Find("li:not([class])").Eq(1).Text()
        time = strings.ReplaceAll(time, " Uhr", "")
        location := selection.Find("p:not([class])").Text()
        if location == "CLUB" {
            switch {
                case strings.HasPrefix(day, "Freitag"):
                    ev_day = "Friday"
                case strings.HasPrefix(day, "Samstag"):
                    ev_day = "Saturday"
                default:
            }
            for _, date := range weekendDates {
                tmp_date := date.Format("2. January")
                if ev_day == date.Weekday().String() && strings.Contains(day, tmp_date){
                    events = add_event_info(events, "dasWerk", date.Weekday().String(),
                        []string{fmt.Sprintf("%s: %s", time, title)})
                }
            }

        }

    })

	coll.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error on '%s': %s", r.Request.URL, err.Error())
	})

	coll.Visit("https://www.daswerk.org/programm/")
    return events
}


func main() {
    weekendDates = getWeekendDates()
    cur_events := events{}
    cur_events.Events = append(cur_events.Events, get_fluc()...)
    cur_events.Events = append(cur_events.Events, get_fish()...)
    cur_events.Events = append(cur_events.Events, get_flex()...)
    cur_events.Events = append(cur_events.Events, get_exil()...)
    cur_events.Events = append(cur_events.Events, get_werk()...)

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
