package cybernews

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"sync"

	colly "github.com/gocolly/colly/v2"
	"github.com/google/uuid"
)

var (
	wg sync.WaitGroup
)

type News struct {
	Id       string `json:"id"`
	Headline string `json:"headlines"`
	FullNews string `json:"news"`
	Url      string `json:"url"`
	Image    string `json:"image"`
}

type NewsFields struct {
	payloadUrl string
	headlines  string
	fullNews   string
	url        string
	image      string
}

var NewsCategory = map[string][]NewsFields{
	"general": {
		{
			payloadUrl: "https://ciosea.economictimes.indiatimes.com/news/next-gen-technologies",
			headlines:  "article.desc div h3.heading",
			fullNews:   "article.desc div p.desktop-view",
			url:        "article.desc figure a",
			image:      "article.desc figure a img",
		},
		{
			payloadUrl: "https://telecom.economictimes.indiatimes.com/news/internet",
			headlines:  "article.desc div h3.heading",
			fullNews:   "article.desc div p.desktop-view",
			url:        "article.desc figure a",
			image:      "article.desc figure a img",
		},
		{
			payloadUrl: "https://ciosea.economictimes.indiatimes.com/news/consumer-tech",
			headlines:  "article.desc div h3.heading",
			fullNews:   "article.desc div p.desktop-view",
			image:      ".desc figure a img",
			url:        ".desc figure a",
		},
	},
	"dataBreach": {
		{
			payloadUrl: "https://thehackernews.com/search/label/data%20breach",
			headlines:  "h2.home-title",
			fullNews:   ".home-desc",
			url:        "a.story-link",
			image:      ".img-ratio img",
		},
		{
			payloadUrl: "https://ciso.economictimes.indiatimes.com/news/data-breaches",
			headlines:  "article.desc div h3.heading",
			fullNews:   "article.desc div p.desktop-view",
			url:        ".desc figure a img",
			image:      ".desc figure a",
		},
	},
	"cyberAttack": {
		{
			payloadUrl: "https://thehackernews.com/search/label/Cyber%20Attack",
			headlines:  "h2.home-title",
			fullNews:   ".home-desc",
			image:      ".img-ratio img",
			url:        "a.story-link",
		},
		{
			payloadUrl: "https://ciso.economictimes.indiatimes.com/news/cybercrime-fraud",
			headlines:  "article.desc div h3.heading",
			fullNews:   "article.desc div p.desktop-view",
			image:      ".desc figure a img",
			url:        ".desc figure a",
		},
	},
	"vulnerability": {
		{
			payloadUrl: "https://thehackernews.com/search/label/Vulnerability",
			headlines:  "h2.home-title",
			fullNews:   ".home-desc",
			image:      ".img-ratio img",
			url:        "a.story-link",
		},
		{
			payloadUrl: "https://ciso.economictimes.indiatimes.com/news/vulnerabilities-exploits",
			headlines:  "article.desc div h3.heading",
			fullNews:   "article.desc div p.desktop-view",
			image:      ".desc figure a img",
			url:        ".desc figure a",
		},
	},
	"malware": {
		{
			payloadUrl: "https://thehackernews.com/search/label/Malware",
			headlines:  "h2.home-title",
			fullNews:   ".home-desc",
			image:      ".img-ratio img",
			url:        "a.story-link",
		},
	},
	"security": {
		{
			payloadUrl: "https://ciosea.economictimes.indiatimes.com/news/security",
			headlines:  "article.desc div h3.heading",
			fullNews:   "article.desc div p.desktop-view",
			image:      ".desc figure a img",
			url:        ".desc figure a",
		},
		{
			payloadUrl: "https://telecom.economictimes.indiatimes.com/tag/hacking",
			headlines:  "article.desc div h3.heading",
			fullNews:   "article.desc div p.desktop-view",
			image:      ".desc figure a img",
			url:        ".desc figure a",
		},
	},
	"cloud": {
		{
			payloadUrl: "https://ciosea.economictimes.indiatimes.com/news/cloud-computing",
			headlines:  "article.desc div h3.heading",
			fullNews:   "article.desc div p.desktop-view",
			image:      ".desc figure a img",
			url:        ".desc figure a",
		},
	},
	"bigData": {
		{
			payloadUrl: "https://ciosea.economictimes.indiatimes.com/news/big-data",
			headlines:  "article.desc div h3.heading",
			fullNews:   "article.desc div p.desktop-view",
			image:      ".desc figure a img",
			url:        ".desc figure a",
		},
		{
			payloadUrl: "https://ciosea.economictimes.indiatimes.com/news/data-center",
			headlines:  "article.desc div h3.heading",
			fullNews:   "article.desc div p.desktop-view",
			image:      ".desc figure a img",
			url:        ".desc figure a",
		},
	},
	"research": {
		{
			payloadUrl: "https://ciosea.economictimes.indiatimes.com/tag/research",
			headlines:  "article.desc div h3.heading",
			fullNews:   "article.desc div p.desktop-view",
			image:      ".desc figure a img",
			url:        ".desc figure a",
		},
	},
	"socialMedia": {
		{
			payloadUrl: "https://telecom.economictimes.indiatimes.com/search/social",
			headlines:  "article.desc div h3.heading",
			fullNews:   "article.desc div p.desktop-view",
			image:      ".desc figure a img",
			url:        ".desc figure a",
		},
	},
	"corporate": {
		{
			payloadUrl: "https://ciosea.economictimes.indiatimes.com/news/corporate",
			headlines:  "article.desc div h3.heading",
			fullNews:   "article.desc div p.desktop-view",
			image:      ".desc figure a img",
			url:        ".desc figure a",
		},
		{
			payloadUrl: "https://telecom.economictimes.indiatimes.com/news/industry",
			headlines:  "article.desc div h3.heading",
			fullNews:   "article.desc div p.desktop-view",
			image:      ".desc figure a img",
			url:        ".desc figure a",
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

	c.OnScraped(func(_ *colly.Response) {
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
		news.FullNews = data[1][index]
	} else {
		news.FullNews = ""
	}

	if index < len(data[2]) {
		news.Url = data[2][index]
	} else {
		news.Url = ""
	}

	if index < len(data[3]) {
		news.Image = data[3][index]
	} else {
		news.Image = ""
	}

	return news
}

func ValidateNewsCategory(newsCategory string) ([]NewsFields, error) {
	NewsSelectors, exists := NewsCategory[newsCategory]
	if !exists {
		return nil, fmt.Errorf("news category: (%s) does not exist", newsCategory)
	}

	return NewsSelectors, nil
}

func GetNews(newsCategory string) ([]News, error) {
	results := make([]News, 0)

	NewsSelectors, err := ValidateNewsCategory(newsCategory)
	if err != nil {
		return nil, err
	}

	c := colly.NewCollector(colly.Async(true))
	c.WithTransport(&http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	})

	for _, newsSelector := range NewsSelectors {
		wg.Add(1)
		go func(newsData NewsFields) {
			defer wg.Done()

			data := make([][]string, 4)
			collyClone := c.Clone()

			// Get Headlines
			scrapeNews(collyClone, newsData.headlines, 0, "headline", &data)

			// Get FullNews
			scrapeNews(collyClone, newsData.fullNews, 1, "fullNews", &data)

			// Get Images
			scrapeNews(collyClone, newsData.image, 2, "image", &data)

			// Get Urls
			scrapeNews(collyClone, newsData.url, 3, "url", &data)

			collyClone.OnError(func(r *colly.Response, err error) {
				log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nand error:", err.Error())
			})

			collyClone.Visit(newsData.payloadUrl)
			collyClone.Wait()

			size := max(len(data[0]), len(data[1]), len(data[2]), len(data[3]))

			for i := 0; i < size; i++ {
				results = append(results, formatNews(data, i))
			}
		}(newsSelector)
	}
	c.Wait()
	wg.Wait()

	return results, nil
}
