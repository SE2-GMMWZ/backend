package main

import (
	"backend/src/repository"
	"backend/src/rest"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"gorm.io/driver/postgres"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {

	dbString := os.Getenv("GOOSE_DBSTRING")
	db, err := sql.Open("postgres", dbString)
	if err != nil {
		log.Fatalf("Failed to open DB connection: %v", err)
	}
	defer db.Close()

	if err := runMigrations(db, "migrations"); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	gormDb, err := gorm.Open(postgres.Open(dbString), &gorm.Config{})
	userRepository := repository.NewUserRepository(gormDb)

	r := gin.Default()
	r.Static("/static", "static")

	rest.AddAuthRoutes(r, userRepository)
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
