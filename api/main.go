package main

import (
	"fmt"
	"os"

	"github.com/ponyo877/news-app-backend-refactor/pkg/annoyindex"

	"github.com/labstack/gommon/log"
	"github.com/nlpodyssey/cybertron/pkg/tasks"
	"github.com/nlpodyssey/cybertron/pkg/tasks/textencoding"

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
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", mysqlConfig.DBUser, mysqlConfig.DBPassword, mysqlConfig.DBHost, mysqlConfig.DBPort, mysqlConfig.DBDatabase)
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
		Addr:     redisConfig.KVSHost + ":" + redisConfig.KVSPort,
		Password: redisConfig.KVSPassword,
		DB:       redisConfig.KVSDatabase,
	})
	defer rdb.Close()

	webdavConfig, err := config.LoadWebDAVConfig()
	if err != nil {
		log.Panicf("LoadWebDAVConfigに失敗しました: %v", err)
	}
	nginxEndpoint := fmt.Sprintf("http://%s:%s", webdavConfig.WDSHost, webdavConfig.WDPort)

	mlmodelConfig, err := config.LoadMLModelConfig()
	if err != nil {
		log.Panicf("LoadMLModelConfigに失敗しました: %v", err)
	}
	mlmodel, err := tasks.Load[textencoding.Interface](
		&tasks.Config{
			ModelsDir: mlmodelConfig.MLModelDir,
			ModelName: mlmodelConfig.MLModelName,
		},
	)
	defer tasks.Finalize(mlmodel)
	annindex := annoyindex.NewAnnoyIndexAngular(256)
	if _, err := os.Stat(mlmodelConfig.MLIndexPath); err == nil {
		if ok := annindex.Load(mlmodelConfig.MLIndexPath); ok {
			log.Info("AnnoyIndexの既存モデルの読み込みに成功しました")
		}
	}

	articlRepository := repository.NewArticleRepository(gormDB, rdb, mlmodel, annindex, mlmodelConfig.MLIndexPath)
	articleService := article.NewService(articlRepository)
	siteRepository := repository.NewSiteMySQL(gormDB)
	userRepository := repository.NewUserMySQL(gormDB)
	commentRepository := repository.NewCommentMySQL(gormDB)
	imageRepository := repository.NewImageNginx(nginxEndpoint)
	siteService := site.NewService(siteRepository)
	stockService := stock.NewService(articleService, siteService)
	fileioService := fileio.NewService(imageRepository, appConfig.APRoot)
	userService := user.NewService(userRepository, fileioService)
	commentService := comment.NewService(commentRepository, userService)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	handler.MakeArticleHandlers(e, articleService)
	handler.MakeStockHandlers(e, stockService)
	handler.MakeSiteHandlers(e, siteService)
	handler.MakeUserHandlers(e, userService)
	handler.MakeImageHandlers(e, fileioService)
	handler.MakeCommentHandlers(e, commentService)

	e.Logger.Fatal(e.Start(":" + appConfig.APPort))
}
