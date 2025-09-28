package main

import (
	"flag"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"go-api-server-sample/cmd/api-server/internal/container"
	"go-api-server-sample/cmd/api-server/internal/controller"
	"go-api-server-sample/cmd/api-server/internal/middleware"
	"go-api-server-sample/internal/infrastructure/database"
)

func main() {
	migrate := flag.Bool("migrate", false, "Run database migration")
	migrateReset := flag.Bool("migrate-reset", false, "Reset database (development only)")
	flag.Parse()

	db, err := database.Connect()
	if err != nil {
		log.Fatal("データベース接続に失敗しました:", err)
	}

	if *migrateReset {
		log.Println("データベースをリセットしています...")
		if err := database.Reset(db); err != nil {
			log.Fatal("データベースリセットに失敗しました:", err)
		}
		log.Println("データベースリセット完了")
		return
	}

	if *migrate {
		log.Println("マイグレーションを実行しています...")
		if err := database.Migrate(db); err != nil {
			log.Fatal("マイグレーション実行に失敗しました:", err)
		}
		log.Println("マイグレーション完了")
		return
	}

	if err := database.Migrate(db); err != nil {
		log.Fatal("マイグレーション実行に失敗しました:", err)
	}

	dependencyContainer := container.NewContainer(db)

	router := setupRouter(dependencyContainer)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("サーバーをポート %s で起動します", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("サーバー起動に失敗しました:", err)
	}
}

func setupRouter(deps *container.Container) *gin.Engine {
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		ginMode = gin.ReleaseMode
	}
	gin.SetMode(ginMode)

	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())

	healthController := controller.NewHealthController(deps.HealthCheckUseCase)
	r.GET("/health", healthController.Check)

	v1 := r.Group("/api/v1")
	v1.Use(middleware.ErrorHandler())

	contentController := controller.NewContentController(
		deps.CreateContentUseCase,
		deps.GetContentUseCase,
		deps.ListContentsUseCase,
		deps.UpdateContentUseCase,
		deps.DeleteContentUseCase,
	)

	contents := v1.Group("/contents")
	{
		contents.POST("", contentController.Create)
		contents.GET("", contentController.List)
		contents.GET("/:id", contentController.GetByID)
		contents.PUT("/:id", contentController.Update)
		contents.DELETE("/:id", contentController.Delete)
	}

	return r
}
