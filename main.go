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


func main() {
    cur_events := events{}
    fmt.Println("Get Fluc info.")
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
