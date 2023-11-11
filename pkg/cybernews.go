package pkg

import (
	"fmt"
	"sync"

	colly "github.com/gocolly/colly/v2"
	"github.com/google/uuid"
)

var wg sync.WaitGroup

type News struct {
	Id       string `json:"id"`
	Headline string `json:"headlines"`
	Author   string `json:"articles"`
	FullNews string `json:"status"`
	Url      string `json:"url"`
	Image    string `json:"image"`
	Date     string `json:"date"`
}

type NewsFieldSelectors struct {
	payloadUrl string
	headlines  string
	author     string
	fullNews   string
	url        string
	image      string
	date       string
}

var NewsType = map[string][]NewsFieldSelectors{
	"general": {
		{
			payloadUrl: "https://ciosea.economictimes.indiatimes.com/news/next-gen-technologies",
			headlines:  "article.desc div h3.heading",
			author:     "",
			fullNews:   "article.desc div p.desktop-view",
			url:        "article.desc figure a",
			image:      "article.desc figure a img",
			date:       "",
		},
		{
			payloadUrl: "https://telecom.economictimes.indiatimes.com/news/internet",
			headlines:  "article.desc div h3.heading",
			author:     "",
			fullNews:   "article.desc div p.desktop-view",
			url:        "article.desc figure a",
			image:      "article.desc figure a img",
			date:       "",
		},
	},
}

func scrapeNews(c *colly.Collector, goQuerySelector string, index int8, name string, data *[][]string) {
	var results []string

	if goQuerySelector == "" {
		(*data)[index] = results
		return
	}

	c.OnHTML(goQuerySelector, func(r *colly.HTMLElement) {
		if name == "image" {
			results = append(results, r.Attr("data-src"))
		} else if name == "url" {
			results = append(results, r.Attr("href"))
		} else {
			results = append(results, r.Text)
		}
	})

	c.OnScraped(func(r *colly.Response) {
		(*data)[index] = results
	})

	c.Wait()
}

func formatNews(data [][]string, index int) News {
	news := News{
		Id: uuid.New().String(),
	}

	if index < len(data[0]) {
		news.Headline = data[0][index]
	} else {
		news.Headline = ""
	}

	if index < len(data[1]) {
		news.Author = data[1][index]
	} else {
		news.Author = ""
	}

	if index < len(data[2]) {
		news.FullNews = data[2][index]
	} else {
		news.FullNews = ""
	}

	if index < len(data[3]) {
		news.Url = data[3][index]
	} else {
		news.Url = ""
	}

	if index < len(data[4]) {
		news.Image = data[4][index]
	} else {
		news.Image = ""
	}

	if index < len(data[5]) {
		news.Date = data[5][index]
	} else {
		news.Date = ""
	}

	return news
}

func GetNews(newsType string) ([]News, error) {
	results := make([]News, 0)

	NewsSelectors, exists := NewsType[newsType]

	if !exists {
		return results, fmt.Errorf("News Type %s does not exist", newsType)
	}

	c := colly.NewCollector(colly.Async(true))

	for _, newsSelector := range NewsSelectors {
		wg.Add(1)
		go func(newsData *NewsFieldSelectors) {
			defer wg.Done()

			data := make([][]string, 6)
			collyClone := c.Clone()

			// Get Headlines
			scrapeNews(collyClone, newsData.headlines, 0, "headline", &data)

			// Get Authors
			scrapeNews(collyClone, newsData.author, 1, "author", &data)

			// Get FullNews
			scrapeNews(collyClone, newsData.fullNews, 2, "fullNews", &data)

			// Get Images
			scrapeNews(collyClone, newsData.image, 3, "image", &data)

			// Get Urls
			scrapeNews(collyClone, ".desc figure a", 4, "url", &data)

			// Get Date
			scrapeNews(collyClone, "", 5, "date", &data)

			collyClone.OnRequest(func(r *colly.Request) {
				fmt.Println("Visiting", r.URL)
			})

			collyClone.OnError(func(r *colly.Response, err error) {
				fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nand error:", err.Error())
			})

			collyClone.Visit(newsData.payloadUrl)
			collyClone.Wait()

			size := max(len(data[0]), len(data[1]), len(data[2]), len(data[3]), len(data[4]), len(data[5]))

			for i := 0; i < size; i++ {
				results = append(results, formatNews(data, i))
			}
		}(&newsSelector)
	}
	wg.Wait()

	return results, nil
}
