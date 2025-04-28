package rest

import (
	"backend/src/model"
	"backend/src/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type bookingController struct {
	bookingRepository *repository.BookingRepository
}

func (bc *bookingController) Create(c *gin.Context) {
	var booking model.Booking
	if err := c.ShouldBindJSON(&booking); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	if err := bc.bookingRepository.CreateBooking(&booking); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create booking"})
		return
	}

	c.JSON(http.StatusCreated, booking)
}

func (bc *bookingController) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
		return
	}

	booking, err := bc.bookingRepository.GetBookingByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "booking not found"})
		return
	}

	c.JSON(http.StatusOK, booking)
}

func (bc *bookingController) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
		return
	}

	var booking model.Booking
	if err := c.ShouldBindJSON(&booking); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	booking.BookingID = id
	if err := bc.bookingRepository.UpdateBooking(&booking); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update booking"})
		return
	}

	c.JSON(http.StatusOK, booking)
}

func (bc *bookingController) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
		return
	}

	if err := bc.bookingRepository.DeleteBooking(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete booking"})
		return
	}

	c.Status(http.StatusNoContent)
}

func AddBookingRoutes(r *gin.Engine, bookingRepository *repository.BookingRepository) {
	bc := bookingController{bookingRepository: bookingRepository}

	bookings := r.Group("/bookings")
	{
		bookings.POST("", bc.Create)
		bookings.GET("/:id", bc.Get)
		bookings.PUT("/:id", bc.Update)
		bookings.DELETE("/:id", bc.Delete)
	}
}
