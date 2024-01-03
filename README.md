# \[Even\]t \[par\]set in \[go\]

Little script to learn go and scrape data from club websites in vienna to show the upcoming events.
Outputs a JSON file, which is then used by an [telegram bot](telbot), to display the information.

The let the bot trigger the update of the JSON automatically, the binary needs to be put inside its
directory. Then the bot will run twice a day. In addition, the update can be triggered with
`/update` (only by the hoster).

To built the binary inside the directory, run the following command:

```
go build -o telbot
```

### Currently parsed websites:

- [B72](https://www.b72.at/program)
- [BlackMarket](http://www.blackmarket.at/?page_id=49)
- [Exil](https://exil1.ticket.io/)
- [Flex](https://flex.at/events/monat/)
- [Fluc Wanne](https://fluc.at/programm/2024_Flucwoche01.html)
- [Grelle Forelle](https://www.grelleforelle.com/programm/)
- [Rhiz](https://rhiz.wien/programm/)
- [SASS](https://sassvienna.com/programm)
- [dasWerk](https://www.daswerk.org/programm/)
- [theLoft](https://www.theloft.at/programm/)

Parsed via [frey-tag.at](https://frey-tag.at/locations/) (an amazing website containing all necessary
information in a beautiful (and parseable) manner) for different reasons:

- [Kramladen](https://frey-tag.at/locations/kramladen) -> website down / no website
- [O-Klub](https://frey-tag.at/locations/o-der-klub) -> 'website=facebook' -> [Ticket Shop](https://shop.eventjet.at/o-vienna) does not contain all events
- [Pratersauna](https://frey-tag.at/locations/pratersauna) -> website info to unreliable for parsing
- [Praterstrasse/PRST](https://frey-tag.at/locations/club-praterstrasse) -> club=IG
- [club-u](https://frey-tag.at/locations/club-u) -> broken site (no events)
- [ponyhof](https://frey-tag.at/locations/ponyhof) -> website=IG

To come:

- [USUS](https://amwasser.wien/events)
