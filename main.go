package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly"
)

type star struct {
	Name      string
	Photo     string
	JobTitle  string
	BirthDate string
	Bio       string
	TopMovies []movie
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
	c := colly.NewCollector(
		colly.AllowedDomains("imdb.com", "www.imdb.com"),
	)
	infoCollector := c.Clone()
	c.OnHTML(".mode-detail", func(e *colly.HTMLElement) {
		profileUrl := e.ChildAttr("div.listener-item-image>a", "href")
		profileUrl = e.Request.AbsoluteURL(profileUrl)
		infoCollector.Visit(profileUrl)
	})
	c.OnHTML("a.lister-page-next", func(e *colly.HTMLElement) {
		nextPage := e.Request.AbsoluteURL(e.Attr("href"))
		c.Visit(nextPage)
	})

	infoCollector.OnHTML("#content-2-wide", func(e *colly.HTMLElement) {
		tmpProfile := star{}
		tmpProfile.Name = e.ChildText("h1.header>span.itemprop")
		tmpProfile.Photo = e.ChildAttr("#name-poster", "src")
		tmpProfile.JobTitle = e.ChildText("#name-job-categories>a>span.itemprop")
		tmpProfile.BirthDate = e.ChildAttr("#name-born-info time", "datetime")
		tmpProfile.Bio = strings.TrimSpace(e.ChildText("#name-bio-text>div.name-trivia-bio-text>div.inline"))

		e.ForEach("div.knownfor-title", func(_ int, kf *colly.HTMLElement) {
			m := movie{}
			m.Title = kf.ChildText("div.knownfor-title-role>a.knownfor-ellipsis")
			m.Year = kf.ChildText("div.knownfor-year>span.knownfor-ellipsis")
			// Append the movie to the TopMovies slice
			tmpProfile.TopMovies = append(tmpProfile.TopMovies, m)
		})
		js, err := json.MarshalIndent(tmpProfile, "", " ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(js))
	})
	// on sending request
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("visiting:", r.URL.String())
	})
	// visiting data
	infoCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("visiting profile URL", r.URL.String())
	})

	startUrl := fmt.Sprintf("https://www.imdb.com/search/name/?birth_monthday=%d-%d", month, day)
	c.Visit(startUrl)
}
