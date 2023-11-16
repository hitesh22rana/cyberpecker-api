package main

import (
	"net/http"

	cybernews "github.com/hitesh22rana/cyberpecker-api/pkg"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var port = ":8000"

func getNews(c echo.Context) error {
	newsType := c.QueryParam("type")
	news, err := cybernews.GetNews(newsType)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, news)
}

func main() {
	e := echo.New()

	api := e.Group("/api/v2", middleware.RemoveTrailingSlash())
	api.GET("/news", getNews)

	e.Logger.Fatal(e.Start(port))
}
