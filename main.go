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
	Events []EvDay `json:"host_events"`
}

type EvDay struct {
	Host   string  `json:"host"`
	Day    string  `json:"day"`
	Events []event `json:"events"`
}

var weekendDates []time.Time

func addEventInfo(events []EvDay, host string, day string, eventInfo event) []EvDay {
	found := false
	for i, ev := range events {
		if ev.Day == day {
			events[i].Events = append(events[i].Events, eventInfo)
			found = true
			break
		}
	}
	if !found {
		eve := EvDay{
			Host:   host,
			Day:    day,
			Events: []event{eventInfo},
		}
		events = append(events, eve)
	}
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

func getClubFlucc(selector string) []EvDay {
	events := []EvDay{}

	coll := colly.NewCollector()
	coll.OnRequest(func(req *colly.Request) {
		// fmt.Println(fmt.Printf("Visiting %s", req.URL))
	})
	coll.OnHTML("section.events-block", func(h *colly.HTMLElement) {
		selection := h.DOM
		evBlock := selection.Find("div.container")
		for _, date := range weekendDates {
			tmpDate := date.Format("02.01.06")

			evBlock.Find("div.day-title").Each(func(_ int, day *goquery.Selection) {
				currentText := strings.TrimSpace(day.Text())
				if strings.Contains(currentText, tmpDate) {
					evList := day.Next()
					evList.Find("li.card").Each(func(_ int, eve_info *goquery.Selection) {
						loc := strings.TrimSpace(eve_info.Find("div.location").Text())
						if strings.Contains(loc, selector) {
							evTime := eve_info.Find("div.time-location-info").Text()
							evTime = strings.TrimSpace(evTime)
							evTime = strings.ReplaceAll(evTime, "\t", "")
							evTime = strings.ReplaceAll(evTime, "\n", "")
							evTime = strings.Split(evTime, "@")[0]
							title := strings.TrimSpace(eve_info.Find("div.title-dimension").Find("h4").First().Text())
							evLink := eve_info.Find("a[href]")
							link, exists := evLink.Attr("href")
							url := ""
							if exists {
								url = fmt.Sprintf("https://flucc.at%s", link)
							}
							eventInfo := event{
								Title: title,
								Time:  evTime,
								URL:   url,
							}
							host := fmt.Sprintf("Flucc %s", selector)
							events = addEventInfo(events, host, date.Weekday().String(), eventInfo)
						}

					})
				}
			})
		}
	})
	coll.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error on '%s': %s", r.Request.URL, err.Error())
	})
	coll.Visit("https://flucc.at")
	return events
}

func getClubFish() []EvDay {
	events := []EvDay{}

	coll := colly.NewCollector()
	coll.OnRequest(func(req *colly.Request) {
		// fmt.Println(fmt.Printf("Visiting %s", req.URL))
	})
	coll.OnHTML("div.project", func(h *colly.HTMLElement) {
		evTime := ""
		selection := h.DOM
		evText := selection.Find("h2")
		fullText := evText.Text()
		for _, date := range weekendDates {
			tmpDate := date.Format("02/01")
			if strings.HasPrefix(fullText, tmpDate) {

				evLink := evText.Find("a[href]")
				link, exists := evLink.Attr("href")
				if exists {
					coll.OnHTML("div.et_pb_text", func(h *colly.HTMLElement) {
						linkSel := h.DOM
						linkSel.Find("p:not([class])").Each(func(_ int, s *goquery.Selection) {
							curText := s.Text()
							if strings.Contains(curText, "DOORS") {
								evTime = curText
							}
						})
					})
					coll.Visit(h.Request.AbsoluteURL(link))
				}

				title := strings.TrimSpace(strings.Trim(fullText, tmpDate))
				evTime = strings.Trim(evTime, "DORS \n")
				eventInfo := event{
					Title: title,
					Time:  evTime,
					URL:   link,
				}
				events = addEventInfo(events, "Grelle Forelle", date.Weekday().String(), eventInfo)
			}
		}
	})
	coll.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error on '%s': %s", r.Request.URL, err.Error())
	})
	coll.Visit("https://www.grelleforelle.com/programm/")
	return events
}

func getClubFlex() []EvDay {
	events := []EvDay{}

	coll := colly.NewCollector()
	coll.OnRequest(func(req *colly.Request) {
		// fmt.Println(fmt.Printf("Visiting %s", req.URL))
	})

	coll.OnHTML("td.tribe-events-calendar-month__day", func(h *colly.HTMLElement) {
		selection := h.DOM
		for _, date := range weekendDates {
			tmpDate := date.Format("2006-01-02")
			selection.Find(fmt.Sprintf("div#tribe-events-calendar-day-%s", tmpDate)).Each(
				func(_ int, sel_day *goquery.Selection) {
					sel_day.Find("article.tribe-events-calendar-month__calendar-event").Each(
						func(_ int, sel_art *goquery.Selection) {
							time := strings.TrimSpace(sel_art.Find("div.tribe-events-calendar-month__calendar-event-datetime").Text())
							time = strings.ReplaceAll(time, "\t", "")
							time = strings.ReplaceAll(time, "\n", "")
							selTitle := sel_art.Find("div.tribe-events-calendar-month__calendar-event-title")

							evLink := selTitle.Find("a[href]")
							link, exists := evLink.Attr("href")
							url := ""
							if exists {
								url = link
							}
							title := strings.TrimSpace(selTitle.Text())
							eventInfo := event{
								Title: title,
								Time:  time,
								URL:   url,
							}
							events = addEventInfo(events, "Flex", date.Weekday().String(), eventInfo)
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

func fixDate(input string) string {
	monthMapping := make(map[string]string)

	monthMapping["Januar"] = "Jänner"
	monthMapping["Februar"] = "Feber"
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

	for enMonth, deMonth := range monthMapping {
		input = strings.Replace(input, enMonth, deMonth, 1)
	}

	return input
}

func getClubWerk() []EvDay {
	events := []EvDay{}

	coll := colly.NewCollector()
	coll.OnRequest(func(req *colly.Request) {
		// fmt.Println(fmt.Printf("Visiting %s", req.URL))
	})

	coll.OnHTML("div.grid", func(h *colly.HTMLElement) {
		elem := h.DOM
		for _, date := range weekendDates {
			tmpDate := date.Format("02.01.")

			elem.Find("div.shrink-0").Each(func(_ int, day *goquery.Selection) {
				dateTime := day.Find("p").First().Text()
				if strings.Contains(dateTime, tmpDate) {
					title := day.Find("p").Eq(1).Text()
					splitted := strings.Split(dateTime, " // ")
					if len(splitted) != 2 {
						fmt.Println("[e] could not parse werk info")
						return
					}
					_, time := splitted[0], splitted[1]
					eventInfo := event{
						Title: title,
						Time:  time,
						URL:   "https://www.daswerk.org/program/",
					}
					events = addEventInfo(events, "dasWerk", date.Weekday().String(), eventInfo)
				}
			})
		}

	})

	coll.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error on '%s': %s", r.Request.URL, err.Error())
	})

	coll.Visit("https://www.daswerk.org/program/")
	return events
}

func getClubLoft() []EvDay {
	events := []EvDay{}

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
			tmpDate := date.Format("02.1.2006")
			if strings.Contains(day, tmpDate) {
				eventInfo := event{
					Title: fmt.Sprintf("%s (%s)", title, location),
					Time:  time,
					URL:   url,
				}
				events = addEventInfo(events, "theLoft", date.Weekday().String(), eventInfo)
			}
		}
	})

	coll.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error on '%s': %s", r.Request.URL, err.Error())
	})

	coll.Visit("https://www.theloft.at/programm/")
	return events
}

func getClubBlack() []EvDay { //lint:ignore U1000
	events := []EvDay{}

	coll := colly.NewCollector()
	coll.OnRequest(func(req *colly.Request) {
		// fmt.Println(fmt.Printf("Visiting %s", req.URL))
	})

	coll.OnHTML("p:not([class])", func(h *colly.HTMLElement) {
		info := h.DOM.Text()
		for _, date := range weekendDates {
			tmpDate := date.Format("02.01.06")
			if strings.Contains(info, tmpDate) {
				splitted := strings.Split(info, " // ")
				if len(splitted) != 3 {
					fmt.Println("[e] could not parse blackmarket info")
					return
				}
				title, evTime := splitted[1], splitted[2]
				evTime = strings.ReplaceAll(evTime, " UHR", "")
				splittedTime := strings.Split(evTime, "-")
				if len(splittedTime) != 2 {
					fmt.Println("[e] could not parse blackmarket time")
					return
				}
				start, _ := strconv.Atoi(splittedTime[0])
				end, _ := strconv.Atoi(splittedTime[1])
				fullTime := fmt.Sprintf("%02d:00-%02d:00", start, end)
				eventInfo := event{
					Title: title,
					Time:  fullTime,
					URL:   "http://www.blackmarket.at/?page_id=49",
				}
				events = addEventInfo(events, "Black Market", date.Weekday().String(), eventInfo)
			}
		}
	})

	coll.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error on '%s': %s", r.Request.URL, err.Error())
	})

	coll.Visit("http://www.blackmarket.at/?page_id=49")
	return events
}

func getClubSass() []EvDay {
	events := []EvDay{}

	coll := colly.NewCollector()
	coll.OnRequest(func(req *colly.Request) {
		// fmt.Println(fmt.Printf("Visiting %s", req.URL))
	})

	coll.OnHTML("div.event", func(h *colly.HTMLElement) {
		selection := h.DOM
		day := strings.TrimSpace(selection.Find("span.start_date").Text())
		for _, date := range weekendDates {
			tmpDate := monday.Format(date, "2. Jan", monday.LocaleDeDE)
			if strings.Contains(day, tmpDate) {

				start := strings.TrimSpace(selection.Find("span.start_time").Text())
				end := strings.TrimSpace(selection.Find("span.end_time").Text())
				title := strings.TrimSpace(selection.Find("div.title").Text())
				subTitle := strings.TrimSpace(selection.Find("div.subline").Text())
				if subTitle != "" {
					subTitle = " " + subTitle
				}

				fullTitle := fmt.Sprintf("%s%s", title, subTitle)
				fullTime := fmt.Sprintf("%s-%s", start, end)

				evLink := selection.Find("a[href]")
				link, exists := evLink.Attr("href")
				url := ""
				if exists {
					url = fmt.Sprintf("https://sassvienna.com%s", link)
				}

				eventInfo := event{
					Title: fullTitle,
					Time:  fullTime,
					URL:   url,
				}
				events = addEventInfo(events, "SASS", date.Weekday().String(), eventInfo)
			}
		}
	})

	coll.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error on '%s': %s", r.Request.URL, err.Error())
	})

	coll.Visit("https://sassvienna.com/programm")
	return events
}

func getClubB72() []EvDay {
	events := []EvDay{}

	coll := colly.NewCollector()
	coll.OnRequest(func(req *colly.Request) {
		// fmt.Println(fmt.Printf("Visiting %s", req.URL))
	})

	coll.OnHTML("div.coming-up", func(h *colly.HTMLElement) {
		evTime := ""
		selection := h.DOM
		day := strings.TrimSpace(selection.Find("h4").Text())
		for _, date := range weekendDates {
			tmpDate := date.Format("02.01")
			if tmpDate == day {
				title := selection.Find("h6")

				evLink := title.Find("a[href]")
				link, exists := evLink.Attr("href")
				url := ""
				if exists {
					url = fmt.Sprintf("https://www.b72.at%s", link)

					coll.OnHTML("div.show-detail", func(h *colly.HTMLElement) {
						linkSel := h.DOM
						linkSel.Find("b:not([class])").Each(func(_ int, s *goquery.Selection) {
							curText := s.Text()
							splitted := strings.Split(curText, " ")
							if len(splitted) != 2 {
								fmt.Println("[e] could not parse b72 time")
								return
							}
							evTime = splitted[1]
						})
					})
					coll.Visit(h.Request.AbsoluteURL(link))
				}
				titleText := strings.TrimSpace(title.Text())

				eventInfo := event{
					Title: titleText,
					Time:  evTime,
					URL:   url,
				}
				events = addEventInfo(events, "B72", date.Weekday().String(), eventInfo)
			}
		}
	})

	coll.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error on '%s': %s", r.Request.URL, err.Error())
	})

	coll.Visit("https://www.b72.at/program")
	return events
}

func getClubFreytag(club string) []EvDay {
	events := []EvDay{}

	coll := colly.NewCollector()
	coll.OnRequest(func(req *colly.Request) {
		// fmt.Println(fmt.Printf("Visiting %s", req.URL))
	})

	coll.OnHTML("div.listKalender_Event__14qVM", func(h *colly.HTMLElement) {
		selection := h.DOM
		day := strings.TrimSpace(selection.Find("span.listKalender_EventDate__hz06c").Text())
		time := strings.TrimSpace(selection.Find("div.listKalender_EventTime__3Xw8c").Text())
		title := strings.TrimSpace(selection.Find("h2").Text())
		subTitle := strings.TrimSpace(selection.Find("h3").Text())
		if subTitle != "" {
			subTitle = " " + subTitle
		}

		evLink := selection.Find("a[href]")
		link, exists := evLink.Attr("href")
		url := ""
		if exists {
			url = fmt.Sprintf("https://frey-tag.at%s", link)
		}
		location := strings.TrimSpace(selection.Find("span.listKalender_EventLocation__2vPrT").Text())
		for _, date := range weekendDates {
			tmpDate := date.Format("02.01.2006")
			if day == tmpDate {
				eventInfo := event{
					Title: fmt.Sprintf("%s%s", title, subTitle),
					Time:  time,
					URL:   url,
				}
				events = addEventInfo(events, location, date.Weekday().String(), eventInfo)
			}
		}
	})

	coll.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error on '%s': %s", r.Request.URL, err.Error())
	})

	coll.Visit(fmt.Sprintf("https://frey-tag.at/locations/%s", club))
	return events
}

func getAllEvents() events {
	currentEvents := events{}

	eventChan := make(chan []EvDay)

	functions := []func() []EvDay{
		getClubFish,
		getClubFlex,
		getClubWerk,
		getClubLoft,
		// get_black,
		getClubSass,
		getClubB72,
	}

	// run funcs in goroutine without argument
	for _, fn := range functions {
		go func(f func() []EvDay) {
			eventChan <- f()
		}(fn)
	}

	// run freytag separate since it has args -> also as gorotines
	freyClubs := []string{"club-praterstrasse", "ponyhof", "club-u", "kramladen", "o-der-klub", "pratersauna", "jolly-roger", "club-exil"}
	for _, clubName := range freyClubs {
		go func(club string) {
			eventChan <- getClubFreytag(club)
		}(clubName)
	}

	// run freytag separate since it has args -> also as gorotines
	fluccLocations := []string{"Wanne", "Deck"}
	for _, location := range fluccLocations {
		go func(loc string) {
			eventChan <- getClubFlucc(loc)
		}(location)
	}

	for i := 0; i < len(functions)+len(freyClubs)+len(fluccLocations); i++ {
		currentEvents.Events = append(currentEvents.Events, <-eventChan...)
	}

	return currentEvents

}

func main() {
	weekendDates = getWeekendDates()
	curretnEvents := getAllEvents()

	content, err := json.MarshalIndent(curretnEvents, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	err = os.WriteFile("events.json", content, 0644)
	if err != nil {
		fmt.Println("Error writing file")
	}

	fmt.Println(string(content))
}
