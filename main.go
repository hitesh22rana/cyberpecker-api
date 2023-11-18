package main

import (
	"errors"
	"log"
	"net/http"
	"os"

	cybernews "github.com/hitesh22rana/cyberpecker-api/src/cybernews"
	database "github.com/hitesh22rana/cyberpecker-api/src/database"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
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
	newsCategory := c.QueryParam("category")
	if _, err := cybernews.ValidateNewsCategory(newsCategory); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	dbClient, ok := c.Get("database").(*redis.Client)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, database.ErrorConnecting.Error())
	}

	data, err := database.RetrieveNews(dbClient, c.Request().Context(), newsCategory)
	if data != nil {
		return c.JSON(http.StatusOK, data)
	}

	if err != nil && (errors.Is(err, database.ErrorConnecting) || errors.Is(err, database.ErrorJsonParsing)) {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	news, err := cybernews.GetNews(newsCategory)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	database.SaveNews(dbClient, c.Request().Context(), newsCategory, news)
	return c.JSON(http.StatusOK, news)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln("error loading .env file")
	}

	dbAddr := os.Getenv("DATABASE_ADDRESS")
	dbPassword := os.Getenv("DATABASE_PASSWORD")
	dbClient := database.NewRedisClient(&redis.Options{
		Addr:     dbAddr,
		Password: dbPassword,
		DB:       0,
	})
	defer dbClient.Close()

	if err := database.Health(dbClient); err != nil {
		log.Fatalln(err)
	}

	e := echo.New()
	e.Use(dbMiddleware(dbClient))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	api := e.Group("/api/v2")
	api.GET("/news", getNews)

	e.Logger.Fatal(e.Start(port))
}
