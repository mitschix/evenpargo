package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
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
    Events []string `json:"events"`
}

var weekendDates []time.Time

func add_event_info(events []EV_Day, host string, day string, event_info []string) ([]EV_Day){
    found := false
    for i, ev := range events {
        if ev.Day == day {
            events[i].Events = append(events[i].Events, event_info...)
            found = true
            break
        }
    }
    if !found {
        eve := EV_Day{
            Host:  host,
            Day:   day,
            Events: event_info,
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
        ev_time := ""
        selection := h.DOM
        ev_text := selection.Find("h2")
        full_text := ev_text.Text()
        for _, date := range weekendDates {
            tmp_date := date.Format("02/01")
            if strings.HasPrefix(full_text, tmp_date) {

                ev_link := ev_text.Find("a[href]")
                link, exists := ev_link.Attr("href")
                if exists {
                    coll.OnHTML("div.et_pb_text",func(h *colly.HTMLElement) {
                        link_sel := h.DOM
                        link_sel.Find("p:not([class])").Each(func(_ int, s *goquery.Selection) {
                            cur_text := s.Text()
                            if strings.Contains(cur_text, "DOORS") {
                                ev_time = cur_text
                            }
                        })
                    })
                    coll.Visit(h.Request.AbsoluteURL(link))
                }

                title := strings.TrimSpace(strings.Trim(full_text, tmp_date))
                ev_time = strings.Trim(ev_time, "DORS \n")
                ev_title := fmt.Sprintf("%s: %s", ev_time, title)
                events = add_event_info(events, "Grelle Forelle", date.Weekday().String(), []string{ev_title})
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

func get_hot() ([]EV_Day){
    events := []EV_Day{}

	coll := colly.NewCollector()
	coll.OnRequest(func(req *colly.Request) {
		// fmt.Println(fmt.Printf("Visiting %s", req.URL))
	})

    coll.OnHTML("div.header_txt", func(h *colly.HTMLElement) {
        selection := h.DOM
        text := strings.TrimSpace(selection.Find("div.marquee").Text())
        ev_days := strings.Split(text, "/////")
        for _, date := range weekendDates {
            tmp_date := date.Format("02.01.2006")
            for _, event := range ev_days {
                ev_splitted := strings.Split(event, "//")
                if len(ev_splitted) == 2 {
                    day, title := ev_splitted[0], strings.TrimSpace(ev_splitted[1])
                    if strings.Contains(day, tmp_date){
                        if title != "TBA" {
                            events = add_event_info(events, "Pratersauna",
                                date.Weekday().String(), []string{title})
                        }
                    }
                }
           }
        }
    })

	coll.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error on '%s': %s", r.Request.URL, err.Error())
	})

	coll.Visit("https://pratersauna.tv/programm/")
    return events
}

func get_loft() ([]EV_Day){
    events := []EV_Day{}

	coll := colly.NewCollector()
	coll.OnRequest(func(req *colly.Request) {
		// fmt.Println(fmt.Printf("Visiting %s", req.URL))
	})

    coll.OnHTML("div.box-wrap", func(h *colly.HTMLElement) {
        selection := h.DOM
        day := strings.TrimSpace(selection.Find("div.datum").Text())
        time := strings.TrimSpace(selection.Find("span.open").Text())
        title := strings.TrimSpace(selection.Find("div.content-middle").Text())
        location := strings.TrimSpace(selection.Find("div.content-right").Text())
        for _, date := range weekendDates {
            tmp_date := date.Format("02.1.2006")
            if strings.Contains(day, tmp_date){
                full_title := fmt.Sprintf("%s: %s (%s)", time, title, location)
                events = add_event_info(events, "theLoft", date.Weekday().String(),
                    []string{full_title})
            }
        }
    })

	coll.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error on '%s': %s", r.Request.URL, err.Error())
	})

	coll.Visit("https://www.theloft.at/programm/")
    return events
}

func get_black() ([]EV_Day){
    events := []EV_Day{}

	coll := colly.NewCollector()
	coll.OnRequest(func(req *colly.Request) {
		// fmt.Println(fmt.Printf("Visiting %s", req.URL))
	})

    coll.OnHTML("p:not([class])", func(h *colly.HTMLElement) {
        info := h.DOM.Text()
        for _, date := range weekendDates {
            tmp_date := date.Format("02.01.06")
            if strings.Contains(info, tmp_date){
                splitted := strings.Split(info, " // ")
                if len(splitted) != 3 {
                    fmt.Println("[e] could not parse blackmarket info")
                    return
                }
                title, ev_time := splitted[1], splitted[2]
                ev_time = strings.ReplaceAll(ev_time, " UHR", "")
                split_time := strings.Split(ev_time, "-")
                if len(split_time) != 2{
                    fmt.Println("[e] could not parse blackmarket time")
                    return
                }
                start, _ := strconv.Atoi(split_time[0])
                end, _ := strconv.Atoi(split_time[1])
                full_title := fmt.Sprintf("%02d:00-%02d:00: %s", start, end, title)
                events = add_event_info(events, "Black Market", date.Weekday().String(),
                    []string{full_title})
                }
        }
    })

	coll.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error on '%s': %s", r.Request.URL, err.Error())
	})

	coll.Visit("http://www.blackmarket.at/?page_id=49")
    return events
}


func get_freytag(club string) ([]EV_Day){
    events := []EV_Day{}

	coll := colly.NewCollector()
	coll.OnRequest(func(req *colly.Request) {
		// fmt.Println(fmt.Printf("Visiting %s", req.URL))
	})

    coll.OnHTML("div.listKalender_Event__14qVM", func(h *colly.HTMLElement) {
        selection := h.DOM
        day := strings.TrimSpace(selection.Find("span.listKalender_EventDate__hz06c").Text())
        time := strings.TrimSpace(selection.Find("div.listKalender_EventTime__3Xw8c").Text())
        title := strings.TrimSpace(selection.Find("h2").Text())
        sub_title := strings.TrimSpace(selection.Find("h3").Text())
        if sub_title != ""{
            sub_title = " " + sub_title
        }
        location := strings.TrimSpace(selection.Find("span.listKalender_EventLocation__2vPrT").Text())
        for _, date := range weekendDates {
            tmp_date := date.Format("02.01.2006")
            if day == tmp_date{
                full_title := fmt.Sprintf("%s: %s%s", time, title, sub_title)
                events = add_event_info(events, location, date.Weekday().String(),
                    []string{full_title})
            }
        }
    })

	coll.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error on '%s': %s", r.Request.URL, err.Error())
	})

	coll.Visit(fmt.Sprintf("https://frey-tag.at/locations/%s",club))
    return events
}

func get_all_events() (events){
    cur_events := events{}

    eventChan := make(chan []EV_Day)

    functions := []func() []EV_Day{
        get_fluc,
        get_fish,
        get_flex,
        get_exil,
        get_werk,
        get_loft,
        get_black,
    }

    // run funcs in goroutine without argument
    for _, fn := range functions {
        go func(f func() []EV_Day) {
            eventChan <- f()
        }(fn)
    }


    // run freytag separate since it has args -> also as gorotines
    frey_clubs := []string{"club-praterstrasse", "ponyhof", "club-u", "kramladen", "o-der-klub"}
    for _, club_name := range frey_clubs {
        go func(club string) {
            eventChan <- get_freytag(club)
        }(club_name)
    }

    for i := 0; i < len(functions) + len(frey_clubs); i++ {
        cur_events.Events = append(cur_events.Events, <-eventChan...)
    }

    return cur_events

}

func main() {
    weekendDates = getWeekendDates()
    cur_events := get_all_events()

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
