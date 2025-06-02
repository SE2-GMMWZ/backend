package repository

import (
	"backend/src/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
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

func (r *BookingRepository) ListBookings(limit, page int, sailorID, dockID, dockOwnerID, paymentStatus, paymentMethod string) ([]model.Booking, int, int, error) {
	var bookings []model.Booking
	db := r.db.Model(&model.Booking{})

	// Apply filters based on provided parameters
	if sailorID != "" {
		db = db.Where("sailor_id = ?", sailorID)
	}
	if dockID != "" {
		db = db.Where("dock_id = ?", dockID)
	}
	if dockOwnerID != "" {
		// Join with docking_spots to filter by owner_id
		db = db.Joins("JOIN docking_spots ON bookings.dock_id = docking_spots.dock_id").
			Where("docking_spots.owner_id = ?", dockOwnerID)
	}
	if paymentStatus != "" {
		db = db.Where("payment_status ILIKE ?", paymentStatus)
	}
	if paymentMethod != "" {
		db = db.Where("payment_method ILIKE ?", paymentMethod)
	}

	offset := (page - 1) * limit
	var totalRecords int64
	db.Count(&totalRecords)

	err := db.Order("booking_id").Limit(limit).Offset(offset).Find(&bookings).Error
	totalPages := int((totalRecords + int64(limit) - 1) / int64(limit))
	return bookings, page, totalPages, err
}
