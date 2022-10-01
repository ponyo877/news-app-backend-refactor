package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-redis/redis/v9"
	"github.com/olivere/elastic/v7"
	"github.com/studio-b12/gowebdav"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ponyo877/news-app-backend-refactor/api/handler"
	"github.com/ponyo877/news-app-backend-refactor/config"
	"github.com/ponyo877/news-app-backend-refactor/repository"
	"github.com/ponyo877/news-app-backend-refactor/usecase/article"
	"github.com/ponyo877/news-app-backend-refactor/usecase/comment"
	"github.com/ponyo877/news-app-backend-refactor/usecase/fileio"
	"github.com/ponyo877/news-app-backend-refactor/usecase/site"
	"github.com/ponyo877/news-app-backend-refactor/usecase/stock"
	"github.com/ponyo877/news-app-backend-refactor/usecase/user"
)

func main() {
	appConfig, err := config.LoadAppConfig()
	if err != nil {
		log.Panicf("LoadAppConfigに失敗しました: %v", err)
	}
	mysqlConfig, err := config.LoadMysqlConfig()
	if err != nil {
		log.Panicf("LoadMysqlConfigに失敗しました: %v", err)
	}
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true", mysqlConfig.DBUser, mysqlConfig.DBPassword, mysqlConfig.DBHost, mysqlConfig.DBDatabase)
	gormDB, err := gorm.Open(mysql.Open(dataSourceName), &gorm.Config{})
	if err != nil {
		log.Panicf("gormDBクライアント作成に失敗しました: %v", err)
	}
	db, err := gormDB.DB()
	if err != nil {
		log.Panicf("SQLクライアント作成に失敗しました: %v", err)
	}
	defer db.Close()

	redisConfig, err := config.LoadRedisConfig()
	if err != nil {
		log.Panicf("LoadRedisConfigに失敗しました: %v", err)
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisConfig.KVSHost + ":6379",
		Password: redisConfig.KVSPassword,
		DB:       redisConfig.KVSDatabase,
	})
	defer rdb.Close()

	elasticSearchConfig, err := config.LoadElasticSearchConfig()
	if err != nil {
		log.Panicf("LoadElasticSearchConfigに失敗しました: %v", err)
	}
	es, err := elastic.NewClient(
		elastic.SetURL("http://"+elasticSearchConfig.SESHost+":9200"),
		elastic.SetBasicAuth(elasticSearchConfig.SEUser, elasticSearchConfig.SEPassword),
		elastic.SetSniff(false),
	)
	if err != nil {
		log.Panicf("ElasticSearchクライアント作成に失敗しました: %v", err)
	}
	webdavConfig, err := config.LoadWebDAVConfig()
	if err != nil {
		log.Panicf("LoadWebDAVConfigに失敗しました: %v", err)
	}

	webdav := gowebdav.NewClient(
		"http://"+webdavConfig.WDSHost+":8080",
		webdavConfig.WDUser,
		webdavConfig.WDPassword,
	)
	// ReadDir("/homes/scott/tmp")
	// if err := webdav.Mkdir("static", 0644); err != nil {
	// 	log.Panicf("staticディレクトリの作成に失敗しました: %v", err)
	// }

	articlRepository := repository.NewArticleRepository(gormDB, rdb, es)
	articleService := article.NewService(articlRepository)
	siteRepository := repository.NewSiteMySQL(gormDB)
	userRepository := repository.NewUserMySQL(gormDB)
	commentRepository := repository.NewCommentMySQL(gormDB)
	imageRepository := repository.NewImageWebDAV(webdav)
	siteService := site.NewService(siteRepository)
	stockService := stock.NewService(articleService, siteService)
	fileioService := fileio.NewService(imageRepository, appConfig.APRoot)
	userService := user.NewService(userRepository, fileioService)
	commentService := comment.NewService(commentRepository, userService)

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	handler.MakeArticleHandlers(e, articleService)
	handler.MakeStockHandlers(e, stockService)
	handler.MakeSiteHandlers(e, siteService)
	handler.MakeUserHandlers(e, userService)
	handler.MakeImageHandlers(e, fileioService)
	handler.MakeCommentHandlers(e, commentService)

	// Start server
	e.Logger.Fatal(e.Start(":80"))
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
