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
	"github.com/goodsign/monday"
)

type event struct {
	Title string `json:"title"`
	Time  string `json:"time"`
	URL   string `json:"url"`
}

type events struct {
	Events []EV_Day `json:"host_events"`
}

type EV_Day struct {
	Host   string  `json:"host"`
	Day    string  `json:"day"`
	Events []event `json:"events"`
}

var weekendDates []time.Time

func add_event_info(events []EV_Day, host string, day string, event_info event) []EV_Day {
	found := false
	for i, ev := range events {
		if ev.Day == day {
			events[i].Events = append(events[i].Events, event_info)
			found = true
			break
		}
	}
	if !found {
		eve := EV_Day{
			Host:   host,
			Day:    day,
			Events: []event{event_info},
		}
		events = append(events, eve)
	}
	return events
}

func get_fluc() []EV_Day {
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
		sel_wanne := selection.Find("ul.info").Find("li.wanne")
		info := strings.TrimSpace(sel_wanne.Text())
		ev_link := sel_wanne.Find("a[href]")
		link, exists := ev_link.Attr("href")
		url := ""
		if exists {
			url = link
		}

		switch day {
		case "Freitag":
			ev_day = "Friday"
			is_weekend = true
		case "Samstag":
			ev_day = "Saturday"
			is_weekend = true
		default:
		}

		if is_weekend && info != "" {
			info_splitted := strings.SplitN(info, " ", 2)

			event_info := event{
				Title: info_splitted[1],
				Time:  strings.TrimRight(info_splitted[0], ":"),
				URL:   url,
			}
			events = add_event_info(events, "Fluc Wanne", ev_day, event_info)
		}
	})
	coll.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error on '%s': %s", r.Request.URL, err.Error())
	})
	coll.Visit(fmt.Sprintf("https://fluc.at/programm/2024_Flucwoche%02d.html", week))
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

func get_fish() []EV_Day {
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
					coll.OnHTML("div.et_pb_text", func(h *colly.HTMLElement) {
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
				event_info := event{
					Title: title,
					Time:  ev_time,
					URL:   link,
				}
				events = add_event_info(events, "Grelle Forelle", date.Weekday().String(), event_info)
			}
		}
	})
	coll.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error on '%s': %s", r.Request.URL, err.Error())
	})
	coll.Visit("https://www.grelleforelle.com/programm/")
	return events
}

func get_flex() []EV_Day {
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
							sel_title := sel_art.Find("h3.tribe-events-calendar-month__calendar-event-title")

							ev_link := sel_title.Find("a[href]")
							link, exists := ev_link.Attr("href")
							url := ""
							if exists {
								url = link
							}
							title := strings.TrimSpace(sel_title.Text())
							event_info := event{
								Title: title,
								Time:  time,
								URL:   url,
							}
							events = add_event_info(events, "Flex", date.Weekday().String(), event_info)
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

func get_exil() []EV_Day {
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
			if strings.Contains(day, tmp_date) {
				event_info := event{
					Title: title,
					Time:  time,
					URL:   "",
				}
				events = add_event_info(events, "Exil", date.Weekday().String(), event_info)
			}
		}
	})

	coll.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error on '%s': %s", r.Request.URL, err.Error())
	})

	coll.Visit("https://exil1.ticket.io/")
	return events
}

func fix_date(input string) string {
	month_mapping := make(map[string]string)

	month_mapping["Januar"] = "Jänner"
	month_mapping["Februar"] = "Feber"
	// should not have a special short name
	// month_mapping["März"] = ""
	// unknown
	// month_mapping["April"] = 4
	// month_mapping["Mai"] = 5
	// month_mapping["Juni"] = 6
	// month_mapping["Juli"] = 7
	// month_mapping["August"] = 8
	// month_mapping["September"] = 9
	// month_mapping["Oktober"] = 10
	// month_mapping["November"] = 11
	// month_mapping["Dezember"] = 12

	for en_month, de_month := range month_mapping {
		input = strings.Replace(input, en_month, de_month, 1)
	}

	return input
}

func get_werk() []EV_Day {
	events := []EV_Day{}

	coll := colly.NewCollector()
	coll.OnRequest(func(req *colly.Request) {
		// fmt.Println(fmt.Printf("Visiting %s", req.URL))
	})

	coll.OnHTML("div.events--preview-item", func(h *colly.HTMLElement) {
		ev_day := ""
		selection := h.DOM
		title := strings.TrimSpace(selection.Find(".preview-item--headline").Text())
		ev_link := selection.Find(".preview-item--link")
		link, exists := ev_link.Attr("href")
		url := ""
		if exists {
			url = link
		}
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
				tmp_date := monday.Format(date, "02. January", monday.LocaleDeDE)
				tmp_date = fix_date(tmp_date)
				if ev_day == date.Weekday().String() && strings.Contains(day, tmp_date) {
					event_info := event{
						Title: title,
						Time:  time,
						URL:   url,
					}
					events = add_event_info(events, "dasWerk", date.Weekday().String(), event_info)
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

func get_loft() []EV_Day {
	events := []EV_Day{}

	coll := colly.NewCollector()
	coll.OnRequest(func(req *colly.Request) {
		// fmt.Println(fmt.Printf("Visiting %s", req.URL))
	})

	coll.OnHTML("a[href]", func(h *colly.HTMLElement) {
		selection := h.DOM
		day := strings.TrimSpace(selection.Find("div.datum").Text())
		time := strings.TrimSpace(selection.Find("span.open").Text())
		title := strings.TrimSpace(selection.Find("div.content-middle").Text())
		location := strings.TrimSpace(selection.Find("div.content-right").Text())

		link, exists := selection.Attr("href")
		url := ""
		if exists {
			url = link
		}

		for _, date := range weekendDates {
			tmp_date := date.Format("02.1.2006")
			if strings.Contains(day, tmp_date) {
				event_info := event{
					Title: fmt.Sprintf("%s (%s)", title, location),
					Time:  time,
					URL:   url,
				}
				events = add_event_info(events, "theLoft", date.Weekday().String(), event_info)
			}
		}
	})

	coll.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error on '%s': %s", r.Request.URL, err.Error())
	})

	coll.Visit("https://www.theloft.at/programm/")
	return events
}

func get_black() []EV_Day {
	events := []EV_Day{}

	coll := colly.NewCollector()
	coll.OnRequest(func(req *colly.Request) {
		// fmt.Println(fmt.Printf("Visiting %s", req.URL))
	})

	coll.OnHTML("p:not([class])", func(h *colly.HTMLElement) {
		info := h.DOM.Text()
		for _, date := range weekendDates {
			tmp_date := date.Format("02.01.06")
			if strings.Contains(info, tmp_date) {
				splitted := strings.Split(info, " // ")
				if len(splitted) != 3 {
					fmt.Println("[e] could not parse blackmarket info")
					return
				}
				title, ev_time := splitted[1], splitted[2]
				ev_time = strings.ReplaceAll(ev_time, " UHR", "")
				split_time := strings.Split(ev_time, "-")
				if len(split_time) != 2 {
					fmt.Println("[e] could not parse blackmarket time")
					return
				}
				start, _ := strconv.Atoi(split_time[0])
				end, _ := strconv.Atoi(split_time[1])
				full_time := fmt.Sprintf("%02d:00-%02d:00", start, end)
				event_info := event{
					Title: title,
					Time:  full_time,
					URL:   "http://www.blackmarket.at/?page_id=49",
				}
				events = add_event_info(events, "Black Market", date.Weekday().String(), event_info)
			}
		}
	})

	coll.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error on '%s': %s", r.Request.URL, err.Error())
	})

	coll.Visit("http://www.blackmarket.at/?page_id=49")
	return events
}

func get_rhiz() []EV_Day {
	events := []EV_Day{}

	coll := colly.NewCollector()
	coll.OnRequest(func(req *colly.Request) {
		// fmt.Println(fmt.Printf("Visiting %s", req.URL))
	})

	coll.OnHTML("div.grid-item", func(h *colly.HTMLElement) {
		selection := h.DOM
		day := strings.TrimSpace(selection.Find("div.event-date").Text())
		for _, date := range weekendDates {
			tmp_date := date.Format("020106")
			if strings.Contains(day, tmp_date) {
				splitted := strings.Split(day, " ")
				if len(splitted) != 3 {
					fmt.Println("[e] could not parse rhiz info")
					return
				}
				time := splitted[2]
				sel_title := selection.Find("h3")
				title := strings.TrimSpace(sel_title.Text())

				ev_link := selection.Find("a[href]")
				link, exists := ev_link.Attr("href")
				url := ""
				if exists {
					url = link
				}
				event_info := event{
					Title: title,
					Time:  time,
					URL:   url,
				}
				events = add_event_info(events, "rhiz", date.Weekday().String(), event_info)
			}
		}
	})

	coll.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error on '%s': %s", r.Request.URL, err.Error())
	})

	coll.Visit("https://rhiz.wien/programm/")
	return events
}

func get_sass() []EV_Day {
	events := []EV_Day{}

	coll := colly.NewCollector()
	coll.OnRequest(func(req *colly.Request) {
		// fmt.Println(fmt.Printf("Visiting %s", req.URL))
	})

	coll.OnHTML("div.event", func(h *colly.HTMLElement) {
		selection := h.DOM
		day := strings.TrimSpace(selection.Find("span.start_date").Text())
		for _, date := range weekendDates {
			tmp_date := monday.Format(date, "2. Jan", monday.LocaleDeDE)
			if strings.Contains(day, tmp_date) {

				start := strings.TrimSpace(selection.Find("span.start_time").Text())
				end := strings.TrimSpace(selection.Find("span.end_time").Text())
				title := strings.TrimSpace(selection.Find("div.title").Text())
				sub_title := strings.TrimSpace(selection.Find("div.subline").Text())
				if sub_title != "" {
					sub_title = " " + sub_title
				}

				full_title := fmt.Sprintf("%s%s", title, sub_title)
				full_time := fmt.Sprintf("%s-%s", start, end)

				ev_link := selection.Find("a[href]")
				link, exists := ev_link.Attr("href")
				url := ""
				if exists {
					url = fmt.Sprintf("https://sassvienna.com%s", link)
				}

				event_info := event{
					Title: full_title,
					Time:  full_time,
					URL:   url,
				}
				events = add_event_info(events, "SASS", date.Weekday().String(), event_info)
			}
		}
	})

	coll.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error on '%s': %s", r.Request.URL, err.Error())
	})

	coll.Visit("https://sassvienna.com/programm")
	return events
}

func get_b72() []EV_Day {
	events := []EV_Day{}

	coll := colly.NewCollector()
	coll.OnRequest(func(req *colly.Request) {
		// fmt.Println(fmt.Printf("Visiting %s", req.URL))
	})

	coll.OnHTML("div.coming-up", func(h *colly.HTMLElement) {
		ev_time := ""
		selection := h.DOM
		day := strings.TrimSpace(selection.Find("h4").Text())
		for _, date := range weekendDates {
			tmp_date := date.Format("02.01")
			if tmp_date == day {
				title := selection.Find("h6")

				ev_link := title.Find("a[href]")
				link, exists := ev_link.Attr("href")
				url := ""
				if exists {
					url = fmt.Sprintf("https://www.b72.at%s", link)

					coll.OnHTML("div.show-detail", func(h *colly.HTMLElement) {
						link_sel := h.DOM
						link_sel.Find("b:not([class])").Each(func(_ int, s *goquery.Selection) {
							cur_text := s.Text()
							splitted := strings.Split(cur_text, " ")
							if len(splitted) != 2 {
								fmt.Println("[e] could not parse b72 time")
								return
							}
							ev_time = splitted[1]
						})
					})
					coll.Visit(h.Request.AbsoluteURL(link))
				}
				title_text := strings.TrimSpace(title.Text())

				event_info := event{
					Title: title_text,
					Time:  ev_time,
					URL:   url,
				}
				events = add_event_info(events, "B72", date.Weekday().String(), event_info)
			}
		}
	})

	coll.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error on '%s': %s", r.Request.URL, err.Error())
	})

	coll.Visit("https://www.b72.at/program")
	return events
}

func get_freytag(club string) []EV_Day {
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
		if sub_title != "" {
			sub_title = " " + sub_title
		}

		ev_link := selection.Find("a[href]")
		link, exists := ev_link.Attr("href")
		url := ""
		if exists {
			url = fmt.Sprintf("https://frey-tag.at%s", link)
		}
		location := strings.TrimSpace(selection.Find("span.listKalender_EventLocation__2vPrT").Text())
		for _, date := range weekendDates {
			tmp_date := date.Format("02.01.2006")
			if day == tmp_date {
				event_info := event{
					Title: fmt.Sprintf("%s%s", title, sub_title),
					Time:  time,
					URL:   url,
				}
				events = add_event_info(events, location, date.Weekday().String(), event_info)
			}
		}
	})

	coll.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error on '%s': %s", r.Request.URL, err.Error())
	})

	coll.Visit(fmt.Sprintf("https://frey-tag.at/locations/%s", club))
	return events
}

func get_all_events() events {
	cur_events := events{}

	eventChan := make(chan []EV_Day)

	functions := []func() []EV_Day{
		get_fluc,
		get_fish,
		get_flex,
		// get_exil, // FIXME: layout changed
		get_werk,
		get_loft,
		get_black,
		get_rhiz,
		get_sass,
		get_b72,
	}

	// run funcs in goroutine without argument
	for _, fn := range functions {
		go func(f func() []EV_Day) {
			eventChan <- f()
		}(fn)
	}

	// run freytag separate since it has args -> also as gorotines
	frey_clubs := []string{"club-praterstrasse", "ponyhof", "club-u", "kramladen", "o-der-klub", "pratersauna"}
	for _, club_name := range frey_clubs {
		go func(club string) {
			eventChan <- get_freytag(club)
		}(club_name)
	}

	for i := 0; i < len(functions)+len(frey_clubs); i++ {
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
