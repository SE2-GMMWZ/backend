package main

import (
	"backend/src/repository"
	"backend/src/rest"
	"database/sql"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"strings"
	"time"
)

func main() {

	dbString := os.Getenv("GOOSE_DBSTRING")
	migrationDir := os.Getenv("GOOSE_MIGRATION_DIR")
	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")

	if dbString == "" || migrationDir == "" {
		log.Fatalf("Missing required environment variables: GOOSE_DBSTRING or GOOSE_MIGRATION_DIR")
	}

	if allowedOrigins == "" {
		log.Fatalf("Missing required environment variable: ALLOWED_ORIGINS")
	}

	db, err := sql.Open("postgres", dbString)
	if err != nil {
		log.Fatalf("Failed to open DB connection: %v", err)
	}
	defer db.Close()

	if err := runMigrations(db, migrationDir); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	gormDb, err := gorm.Open(postgres.Open(dbString), &gorm.Config{})
	userRepository := repository.NewUserRepository(gormDb)
	bookingRepository := repository.NewBookingRepository(gormDb)
	commentRepository := repository.NewCommentRepository(gormDb)
	dockingSpotRepository := repository.NewDockingSpotRepository(gormDb)
	guideRepository := repository.NewGuideRepository(gormDb)
	notificationRepository := repository.NewNotificationRepository(gormDb)
	portRepository := repository.NewPortRepository(gormDb)
	reviewRepository := repository.NewReviewRepository(gormDb)

	r := gin.Default()

	// Split ALLOWED_ORIGINS into a slice
	origins := strings.Split(allowedOrigins, ",")

	r.Use(cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Static("/static", "static")
	rest.AddAuthRoutes(r, userRepository)
	rest.AddBookingRoutes(r, bookingRepository)
	rest.AddCommentRoutes(r, commentRepository)
	rest.AddDockingSpotRoutes(r, dockingSpotRepository)
	rest.AddGuideRoutes(r, guideRepository)
	rest.AddNotificationRoutes(r, notificationRepository)
	rest.AddPortRoutes(r, portRepository)
	rest.AddReviewRoutes(r, reviewRepository)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}

func runMigrations(db *sql.DB, migrationDir string) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set Goose dialect: %w", err)
	}

	if err := goose.Up(db, migrationDir); err != nil {
		return fmt.Errorf("goose up failed: %w", err)
	}

	return nil
}
