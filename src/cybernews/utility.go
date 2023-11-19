package cybernews

import (
	"strings"

	colly "github.com/gocolly/colly/v2"
	"github.com/google/uuid"
)

func scrapeNews(c *colly.Collector, querySelector string, index int8, field string, data *[][]string) {
	var results []string

	if querySelector == "" {
		(*data)[index] = results
		return
	}

	c.OnHTML(querySelector, func(r *colly.HTMLElement) {
		if field == "image" {
			results = append(results, r.Attr("data-src"))
		} else if field == "link" {
			results = append(results, r.Attr("href"))
		} else {
			results = append(results, r.Text)
		}
	})

	c.OnScraped(func(_ *colly.Response) {
		(*data)[index] = results
	})

	c.Wait()
}

func fixNewsUrl(link string, url string) string {
	if strings.Contains(link, "https://") {
		return link
	}

	return url + link
}

func formatNews(data [][]string, index int, source string, url string) News {
	news := News{
		Id:     uuid.New().String(),
		Source: source,
	}

	if index < len(data[0]) {
		news.Headline = data[0][index]
	}

	if index < len(data[1]) {
		news.News = data[1][index]
	}

	if index < len(data[2]) {
		news.Link = fixNewsUrl(data[2][index], url)
	}

	if index < len(data[3]) {
		news.Image = data[3][index]
	}

	return news
}
