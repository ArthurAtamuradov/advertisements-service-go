// cmd/your_app_name/main.go
package main

import (
	"database/sql"
	"fmt"
	"log"
	http "net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/go-sql-driver/mysql"
	migrate "github.com/rubenv/sql-migrate"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	ihttp "advertisements-service/internal/delivery/http"

	"advertisements-service/internal/repositories"

	"advertisements-service/internal/usecases"
)

var db *sql.DB

func main() {
	loadConfig()

	connectDB()
	defer db.Close()

	err := runMigrations(db)
	if err != nil {
		log.Fatal(err)
	}

	adRepo := repositories.NewAdvertisementRepository(db)
	adService := usecases.NewAdvertisementService(adRepo)
	adHandler := ihttp.NewAdvertisementHandler(adService)

	router := setupRouter(adHandler)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", viper.GetInt("server.port")),
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error: %v\n", err)
		}
	}()

	waitForShutdown(server)
}


func loadConfig() {
	viper.SetConfigFile("config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}
}

func connectDB() {
	var err error
	db, err = sql.Open("mysql", viper.GetString("database.url"))
	if err != nil {
		log.Fatalf("Error connecting to the database: %s", err)
	}


	if err := createDatabase(db, viper.GetString("database.name")); err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Error pinging the database: %s", err)
	}
}

func setupRouter(handler *ihttp.AdvertisementHandler) *gin.Engine {
	router := gin.Default()

	router.GET("/advertisements", handler.GetAllAdvertisements)
	router.GET("/advertisements/:id", handler.GetAdvertisement)
	router.POST("/advertisements", handler.CreateAdvertisement)
	router.PUT("/advertisements/:id", handler.UpdateAdvertisement)
	router.DELETE("/advertisements/:id", handler.DeleteAdvertisement)

	return router
}

func runMigrations(db *sql.DB) error {
	migrations := &migrate.FileMigrationSource{
		Dir: "migrations",
	}
	fmt.Printf("%s", migrations)

	n, err := migrate.Exec(db, "mysql", migrations, migrate.Up)

	if err != nil {
		return err
	}

	fmt.Printf("Applied %d migrations!\n", n)
	return nil
}

func waitForShutdown(server *http.Server) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	log.Println("Shutting down gracefully...")

	if err := server.Shutdown(nil); err != nil {
		log.Fatalf("Error during server shutdown: %s", err)
	}
}

func createDatabase(db *sql.DB, dbName string) error {
	_, err := db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName))
	return err
}