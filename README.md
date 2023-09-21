# \[Even]t \[par]ser in \[go]

Little script to learn go and scrape data from club websites in vienna to show the upcoming events.
Target is to output a json file to feed to a bot/web page to display them nicely.

Currently parsed websites:
- [Fluc Wanne](https://fluc.at/programm/2023_Flucwoche%d.html)
- [Grelle Forelle](https://www.grelleforelle.com/programm/)
- [Flex](https://flex.at/events/monat/)
- [Exil](https://exil1.ticket.io/)
- [dasWerk](https://www.daswerk.org/programm/)
- [Pratersauna](https://pratersauna.tv/programm/)
- [theLoft](https://www.theloft.at/programm/)
- [BlackMarket](http://www.blackmarket.at/?page_id=49)

Parsed via [frey-tag.at](https://frey-tag.at/locations/) (an amazing website containing all necessary
information in a beautiful (and parseable) manner) for different reasons:
- [Kramladen](https://frey-tag.at/locations/kramladen) -> website down / no website
- [Praterstrasse/PRST](https://frey-tag.at/locations/club-praterstrasse) -> club=IG
- [ponyhof](https://frey-tag.at/locations/ponyhof) -> website=IG
- [club-u](https://frey-tag.at/locations/club-u) -> broken site (no events)
- [O-Klub](https://frey-tag.at/locations/o-der-klub) -> 'website=facebook' -> [Ticket Shop](https://shop.eventjet.at/o-vienna) does not contain all events

To come:
- [B72](https://www.b72.at/program)
- [Rhiz](https://rhiz.wien/programm/)
- [SASS](https://sassvienna.com/programm)
- [USUS](https://amwasser.wien/events)
