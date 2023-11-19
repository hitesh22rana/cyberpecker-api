package cybernews

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"sync"

	colly "github.com/gocolly/colly/v2"
)

var (
	wg sync.WaitGroup
)

func ValidateNewsCategory(category string) ([]NewsFields, error) {
	NewsSelectors, exists := newsCategory[category]
	if !exists {
		return nil, fmt.Errorf("news category: (%s) does not exist", category)
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

			scrapeNews(collyClone, newsData.headline, 0, "headline", &data)
			scrapeNews(collyClone, newsData.news, 1, "news", &data)
			scrapeNews(collyClone, newsData.link, 2, "link", &data)
			scrapeNews(collyClone, newsData.image, 3, "image", &data)

			collyClone.OnError(func(r *colly.Response, err error) {
				log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nand error:", err.Error())
			})

			collyClone.Visit(newsData.url)
			collyClone.Wait()

			size := max(len(data[0]), len(data[1]), len(data[2]), len(data[3]))

			for i := 0; i < size; i++ {
				news := formatNews(data, i, newsData.source, newsData.url)
				if advertisement := news.checkAdvertisement(newsData.domain); advertisement {
					continue
				}
				results = append(results, news)
			}
		}(newsSelector)
	}
	c.Wait()
	wg.Wait()

	return results, nil
}
