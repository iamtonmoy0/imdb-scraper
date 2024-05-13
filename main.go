package main

import (
	"flag"
	"fmt"

	"github.com/gocolly/colly"
)

type star struct {
	Name      string
	Photo     string
	JobTitle  string
	BirthDate string
	Bio       string
	TopMovies string
}
type movie struct {
	Title string
	Year  string
}

func main() {
	month := flag.Int("month", 1, "Month to fetch birthdays for")
	day := flag.Int("day", 1, "day to fetch birthday for")
	flag.Parse()
	crawl(*month, *day)
}

func crawl(month int, day int) {
c:=colly.NewCollector(
	colly.AllowedDomains("imdb.com","www.imdb.com")
)
infoCollector:=c.Clone()
c.OnHTML(".mode-detail",func(e *colly.HTMLElement){
	profileUrl:=e.ChildAttr("div.listener-item-image>a","href")
	profileUrl=e.Request.AbsoluteURL(profileUrl)
	infoCollector.Visit(profileUrl)
})

	startUrl := fmt.Sprintf("https://www.imdb.com/search/name/?birth monthday=%d-%d", month, day)
	c.Visit(startUrl)
}
