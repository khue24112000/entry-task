package main

import (
	"entry-project/back-end/internal/config"
	"entry-project/back-end/internal/handler"
	"entry-project/back-end/internal/repository"
	"entry-project/back-end/internal/routes"
	"entry-project/back-end/internal/service"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()
	db := config.ConnectDB(cfg)
	rdb := config.ConnectRedis(cfg)

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	sqlDB.SetMaxOpenConns(1000)
	sqlDB.SetMaxIdleConns(50)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	sqlDB.SetConnMaxIdleTime(2 * time.Minute)

	config.MigrateDB(db)

	gin.SetMode(gin.DebugMode)

	r := gin.Default()

	// khoi tao user
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, rdb)
	userHandler := handler.NewUserHandler(userService)
	routes.SetupRoutes(r, userHandler)

	go func() {
		log.Printf("Server is running on port %s", cfg.AppPort)
		r.Run(":" + cfg.AppPort)
	}()

	if os.Getenv("GIN_PPROF_ENABLE") == "true" {
		go func() {
			log.Println("pprof on http://localhost:6060/debug/pprof/")
			// net/http/pprof đã tự đăng ký handler trên DefaultServeMux
			if err := http.ListenAndServe("localhost:6060", nil); err != nil {
				log.Fatal(err)
			}
		}()
	}

	select {}

}
