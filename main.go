package main

import (
	"errors"
	"net/http"

	cybernews "github.com/hitesh22rana/cyberpecker-api/pkg/cybernews"
	database "github.com/hitesh22rana/cyberpecker-api/pkg/database"
	"github.com/redis/go-redis/v9"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	port               = ":8000"
	databaseContextkey = "database"
)

func dbMiddleware(databaseClient *redis.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(databaseContextkey, databaseClient)
			return next(c)
		}
	}
}

func getNews(c echo.Context) error {
	newsType := c.QueryParam("type")
	if _, err := cybernews.ValidateNewsType(newsType); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	dbClient, ok := c.Get("database").(*redis.Client)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, database.ErrorConnecting.Error())
	}

	data, err := database.RetrieveNews(dbClient, c.Request().Context(), newsType)
	if err == nil {
		return c.JSON(http.StatusOK, data)
	}

	if err != nil && errors.Is(err, database.ErrorJsonParsing) {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	news, err := cybernews.GetNews(newsType)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	database.SaveNews(dbClient, c.Request().Context(), newsType, news)
	return c.JSON(http.StatusOK, news)
}

func main() {
	databaseClient := database.NewRedisClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer databaseClient.Close()

	e := echo.New()
	e.Use(dbMiddleware(databaseClient))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	api := e.Group("/api/v2")
	api.GET("/news", getNews)

	e.Logger.Fatal(e.Start(port))
}
