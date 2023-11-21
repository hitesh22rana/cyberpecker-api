package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	cache "github.com/hitesh22rana/cyberpecker-api/src/cache"
	cybernews "github.com/hitesh22rana/cyberpecker-api/src/cybernews"
	"github.com/joho/godotenv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
)

var (
	port               = ":8000"
	lruCacheContextKey = "lruCache"
	databaseContextkey = "database"
)

func addMiddleware[T any](client *T, contextKey string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(contextKey, client)
			return next(c)
		}
	}
}

func getNews(c echo.Context) error {
	newsCategory := c.QueryParam("category")
	if _, err := cybernews.ValidateNewsCategory(newsCategory); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	lruCache, ok := c.Get(lruCacheContextKey).(*cache.LRUCache)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, cache.ErrorConnecting.Error())
	}

	cachedData := lruCache.GetNews(newsCategory)
	if cachedData != nil {
		return c.JSON(http.StatusOK, cachedData)
	}

	dbClient, ok := c.Get(databaseContextkey).(*cache.RedisClient)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, cache.ErrorConnecting.Error())
	}

	data, err := dbClient.GetNews(c.Request().Context(), newsCategory)
	if data != nil {
		lruCache.SetNews(newsCategory, data)
		return c.JSON(http.StatusOK, data)
	}

	if err != nil && (errors.Is(err, cache.ErrorConnecting) || errors.Is(err, cache.ErrorJsonParsing)) {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	news, err := cybernews.GetNews(newsCategory)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	dbClient.SetNews(c.Request().Context(), newsCategory, news)
	return c.JSON(http.StatusOK, news)
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	errEnvDev := godotenv.Load(".env")
	errEnvProd := godotenv.Load(".env.prod")
	if errEnvDev != nil && errEnvProd != nil {
		log.Fatalln("error loading environment files")
	}
}

func main() {
	dbAddr := os.Getenv("DATABASE_ADDRESS")
	dbPassword := os.Getenv("DATABASE_PASSWORD")
	dbClient := cache.NewRedisClient(&redis.Options{
		Addr:     dbAddr,
		Password: dbPassword,
		DB:       0,
	})
	defer dbClient.Close()

	if err := dbClient.Health(); err != nil {
		log.Fatalln(err)
	}

	lruCache := cache.NewLRUCache(10, 120)

	e := echo.New()
	e.Use(addMiddleware(dbClient, databaseContextkey))
	e.Use(addMiddleware(lruCache, lruCacheContextKey))
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${remote_ip} -> ${method} ${uri} ${status} ${latency_human}\n",
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	api := e.Group("/api/v2")
	api.GET("/news", getNews)

	go func() {
		if err := e.Start(port); err != http.ErrServerClosed {
			e.Logger.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err := e.Shutdown(context.Background()); err != nil {
		e.Logger.Fatal(err)
	}
}
