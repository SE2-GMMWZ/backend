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

func (r *BookingRepository) ListBookings(limit, offset int) ([]model.Booking, error) {
	var bookings []model.Booking
	err := r.db.Order("booking_id").Limit(limit).Offset(offset).Find(&bookings).Error
	return bookings, err
}
