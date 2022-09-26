package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-redis/redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ponyo877/news-app-backend-refactor/api/handler"
	"github.com/ponyo877/news-app-backend-refactor/config"
	"github.com/ponyo877/news-app-backend-refactor/repository"
	"github.com/ponyo877/news-app-backend-refactor/usecase/article"
	"github.com/ponyo877/news-app-backend-refactor/usecase/site"
	"github.com/ponyo877/news-app-backend-refactor/usecase/stock"
)

func main() {

	mysqlConfig, err := config.LoadMysqlConfig()
	if err != nil {
		log.Panic(err.Error())
	}
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true", mysqlConfig.DBUser, mysqlConfig.DBPassword, mysqlConfig.DBHost, mysqlConfig.DBDatabase)
	gormDB, err := gorm.Open(mysql.Open(dataSourceName), &gorm.Config{})
	if err != nil {
		log.Panic(err.Error())
	}
	db, err := gormDB.DB()
	defer db.Close()

	redisConfig, err := config.LoadRedisConfig()
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisConfig.KVSHost + ":6379",
		Password: redisConfig.KVSPassword,
		DB:       redisConfig.KVSDatabase,
	})
	defer rdb.Close()

	articlRepository := repository.NewArticleRepository(gormDB, rdb)
	articleService := article.NewService(articlRepository)
	siteRepository := repository.NewSiteMySQL(gormDB)
	siteService := site.NewService(siteRepository)
	stockService := stock.NewService(articleService, siteService)

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	handler.MakeArticleHandlers(e, articleService)
	handler.MakeStockHandlers(e, stockService)

	// Start server
	e.Logger.Fatal(e.Start(":80"))
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
