package main

import (
	"fmt"
	"time"

	"github.com/gocolly/colly"
)

type events struct {
    Host string `json:"host"`
    Events []evday `json:"host_events"`
}

type evday struct {
    Day string `json:"day"`
    Info string `json:"info"`
}

func get_fluc() (events){
    fluc_events := events{}
    fluc_events.Host = "Fluc"
	tn := time.Now().UTC()
	_, week := tn.ISOWeek()

	coll := colly.NewCollector()
	coll.OnRequest(func(req *colly.Request) {
		// fmt.Println(fmt.Printf("Visiting %s", req.URL))
	})
	coll.OnHTML("li.datum", func(h *colly.HTMLElement) {
        fluc_ev := evday{}
        selection := h.DOM
        days := selection.Find("span.tag")
        info := selection.Find("ul.info").Find("a[href]")
        switch days.Text() {
        case "Freitag":
            fluc_ev.Day = days.Text()
            fluc_ev.Info = info.Text()
            fluc_events.Events = append(fluc_events.Events, fluc_ev)
        case "Samstag":
            fluc_ev.Day = days.Text()
            fluc_ev.Info = info.Text()
            fluc_events.Events = append(fluc_events.Events, fluc_ev)
        default:
        }
	})
	coll.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error on '%s': %s", r.Request.URL, err.Error())
	})
	coll.Visit(fmt.Sprintf("https://fluc.at/programm/2023_Flucwoche%d.html", week-1))
    return fluc_events
}

func main() {
    cur_events := []events{}
    fmt.Println("Get Fluc info.")
    cur_events = append(cur_events, get_fluc())

    fmt.Println(cur_events)
}
