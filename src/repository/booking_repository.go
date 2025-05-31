package repository

import (
	"backend/src/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
)

type BookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) *BookingRepository {
	return &BookingRepository{db}
}

func (r *BookingRepository) CreateBooking(booking *model.Booking) error {
	return r.db.Create(booking).Error
}

func (r *BookingRepository) GetBookingByID(id uuid.UUID) (*model.Booking, error) {
	var booking model.Booking
	err := r.db.First(&booking, "booking_id = ?", id).Error
	return &booking, err
}

func (r *BookingRepository) UpdateBooking(booking *model.Booking) error {
	return r.db.Save(booking).Error
}

func (r *BookingRepository) DeleteBooking(id uuid.UUID) error {
	return r.db.Delete(&model.Booking{}, "booking_id = ?", id).Error
}



func (r *BookingRepository) ListBookings(limit, offset int, query string) ([]model.Booking, error) {
	var bookings []model.Booking
	db := r.db.Model(&model.Booking{})

	// Basic query parsing: support "sailor_id=...", "dock_id=...", "payment_status=...", "payment_method=..."
	if query != "" {
		parts := strings.Split(query, "&")
		for _, part := range parts {
			kv := strings.SplitN(part, "=", 2)
			if len(kv) != 2 {
				continue
			}
			key := kv[0]
			value := kv[1]
			switch key {
			case "sailor_id":
				db = db.Where("sailor_id = ?", value)
			case "dock_id":
				db = db.Where("dock_id = ?", value)
			case "payment_status":
				db = db.Where("payment_status ILIKE ?", value)
			case "payment_method":
				db = db.Where("payment_method ILIKE ?", value)
			}
		}
	}

	err := db.Order("booking_id").Limit(limit).Offset(offset).Find(&bookings).Error
	return bookings, err
}
