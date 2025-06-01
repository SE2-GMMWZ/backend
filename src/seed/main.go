package main

import (
	"backend/src/model"
	"backend/src/repository"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// Initialize faker
	gofakeit.Seed(time.Now().UnixNano())
	// Seed math/rand for random selections
	rand.Seed(time.Now().UnixNano())

	// Database connection (from environment)
	dbString := os.Getenv("GOOSE_DBSTRING")

	db, err := gorm.Open(postgres.Open(dbString), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	bookingRepo := repository.NewBookingRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	dockingSpotRepo := repository.NewDockingSpotRepository(db)
	guideRepo := repository.NewGuideRepository(db)
	notificationRepo := repository.NewNotificationRepository(db)
	reviewRepo := repository.NewReviewRepository(db)

	// Create users
	var users []model.User
	var sailors []model.User
	var regularUsers []model.User
	roles := []model.UserRole{model.RoleEditor, model.RoleAdmin, model.RoleDockOwner, model.RoleSailor}
	for i := 0; i < 200; i++ {
		u := model.User{
			UserID:   uuid.New(),
			Name:     gofakeit.FirstName(),
			Surname:  gofakeit.LastName(),
			Email:    gofakeit.Email(),
			Password: "password",
			Role:     roles[rand.Intn(len(roles))],
		}
		if err := userRepo.CreateUser(&u); err != nil {
			log.Fatalf("failed to create user %d: %v", i, err)
		}
		users = append(users, u)
		// Collect sailors and regular users for bookings/comments
		switch u.Role {
		case model.RoleSailor:
			sailors = append(sailors, u)
		case model.UserRole("user"):
			regularUsers = append(regularUsers, u)
		}
	}

	// Create docking spots
	var docks []model.DockingSpot
	for i := 0; i < 200; i++ {
		owner := users[rand.Intn(len(users))]
		loc := map[string]interface{}{"town": gofakeit.City(), "latitude": gofakeit.Latitude(), "longitude": gofakeit.Longitude()}
		locJSON, _ := json.Marshal(loc)
		ds := model.DockingSpot{
			DockID:          uuid.New(),
			Name:            gofakeit.Company(),
			Location:        locJSON,
			Description:     ptrString(gofakeit.Sentence(10)),
			OwnerID:         owner.UserID,
			Services:        ptrString(gofakeit.Word()),
			ServicesPricing: ptrFloat64(gofakeit.Float64Range(10, 100)),
			PricePerNight:   gofakeit.Float64Range(50, 500),
			PricePerPerson:  ptrFloat64(gofakeit.Float64Range(20, 200)),
			Availability:    gofakeit.RandomString([]string{"available", "unavailable"}),
		}
		if err := dockingSpotRepo.CreateDockingSpot(&ds); err != nil {
			log.Fatalf("failed to create docking spot %d: %v", i, err)
		}
		docks = append(docks, ds)
	}

	// Create guides
	var guides []model.Guide
	for i := 0; i < 200; i++ {
		author := users[rand.Intn(len(users))]
		loc := map[string]interface{}{"town": gofakeit.City(), "latitude": gofakeit.Latitude(), "longitude": gofakeit.Longitude()}
		locJSON, _ := json.Marshal(loc)
		g := model.Guide{
			GuideID:         uuid.New(),
			Title:           gofakeit.Sentence(3),
			Content:         gofakeit.Paragraph(1, 2, 10, " "),
			AuthorID:        author.UserID,
			PublicationDate: gofakeit.Date(),
			Location:        locJSON,
			IsApproved:      gofakeit.Bool(),
		}
		if err := guideRepo.CreateGuide(&g); err != nil {
			log.Fatalf("failed to create guide %d: %v", i, err)
		}
		guides = append(guides, g)
	}

	// Create bookings
	for i := 0; i < 200; i++ {
		if len(sailors) == 0 {
			log.Fatalf("no sailors available for bookings")
		}
		sailor := sailors[rand.Intn(len(sailors))]
		dock := docks[rand.Intn(len(docks))]
		start := gofakeit.Date()
		end := start.AddDate(0, 0, gofakeit.Number(1, 7))
		b := model.Booking{
			BookingID:     uuid.New(),
			SailorID:      sailor.UserID,
			DockID:        dock.DockID,
			StartDate:     start,
			EndDate:       end,
			PaymentMethod: gofakeit.RandomString([]string{"online", "in-person"}),
			PaymentStatus: gofakeit.RandomString([]string{"paid", "unpaid"}),
			People:        gofakeit.Number(1, 5),
		}
		if err := bookingRepo.CreateBooking(&b); err != nil {
			log.Fatalf("failed to create booking %d: %v", i, err)
		}
	}

	// Create comments
	for i := 0; i < 200; i++ {
		if len(regularUsers) == 0 {
			log.Fatalf("no regular users available for comments")
		}
		user := regularUsers[rand.Intn(len(regularUsers))]
		guide := guides[rand.Intn(len(guides))]
		c := model.Comment{
			CommentID: uuid.New(),
			GuideID:   guide.GuideID,
			UserID:    user.UserID,
			Content:   gofakeit.Sentence(10),
			Timestamp: gofakeit.Date(),
		}
		if err := commentRepo.CreateComment(&c); err != nil {
			log.Fatalf("failed to create comment %d: %v", i, err)
		}
	}

	// Create notifications
	for i := 0; i < 200; i++ {
		user := users[rand.Intn(len(users))]
		n := model.Notification{
			NotificationID: uuid.New(),
			UserID:         user.UserID,
			Message:        gofakeit.Sentence(5),
			Timestamp:      time.Now(),
		}
		if err := notificationRepo.CreateNotification(&n); err != nil {
			log.Fatalf("failed to create notification %d: %v", i, err)
		}
	}

	// Create reviews
	for i := 0; i < 200; i++ {
		user := users[rand.Intn(len(users))]
		r := model.Review{
			ReviewID:     uuid.New(),
			ReviewerID:   user.UserID,
			Rating:       gofakeit.Float64Range(1, 5),
			Comment:      ptrString(gofakeit.Sentence(15)),
			DateOfReview: time.Now(),
		}
		if err := reviewRepo.CreateReview(&r); err != nil {
			log.Fatalf("failed to create review %d: %v", i, err)
		}
	}

	fmt.Println("Database seeding completed successfully.")
}

// Helpers for pointer values
func ptrString(s string) *string    { return &s }
func ptrFloat64(f float64) *float64 { return &f }
