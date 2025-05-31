package rest

import (
	"backend/src/model"
	"backend/src/repository"
	"backend/src/rest/shared"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type bookingController struct {
	bookingRepository *repository.BookingRepository
}

// Create godoc
// @Summary      Create booking
// @Description  Creates a new booking
// @Tags         Bookings
// @Accept       json
// @Produce      json
// @Param        body  body      model.Booking  true  "Booking to create"
// @Success      201   {object}  model.Booking
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /bookings [post]
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

// Get godoc
// @Summary      Get booking
// @Description  Returns a single booking by its UUID
// @Tags         Bookings
// @Produce      json
// @Param        id   path      string        true  "Booking UUID"
// @Success      200  {object}  model.Booking
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /bookings/{id} [get]
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

// Update godoc
// @Summary      Update booking
// @Description  Updates a booking by its UUID
// @Tags         Bookings
// @Accept       json
// @Produce      json
// @Param        id    path      string        true  "Booking UUID"
// @Param        body  body      model.Booking true  "Booking data to update"
// @Success      200   {object}  model.Booking
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /bookings/{id} [put]
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

// Delete godoc
// @Summary      Delete booking
// @Description  Deletes a booking by its UUID
// @Tags         Bookings
// @Produce      json
// @Param        id   path      string  true  "Booking UUID"
// @Success      204  "No Content"
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /bookings/{id} [delete]
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

// List godoc
// @Summary      List bookings
// @Description  Returns paginated list of bookings; supports query parameters forwarded to repository
// @Tags         Bookings
// @Produce      json
// @Param        sailor_id      query     string false "Filter by Sailor UUID"
// @Param        dock_id        query     string false "Filter by Dock UUID"
// @Param        payment_status query     string false "Filter by payment status"
// @Param        payment_method query     string false "Filter by payment method"
// @Param        limit          query     int    false "Items per page"
// @Success      200    {object}  []model.Booking
// @Router       /bookings/list [get]
func (bc *bookingController) List(c *gin.Context) {
	shared.ListWithQuery(c, bc.bookingRepository.ListBookings)
}

func AddBookingRoutes(r *gin.Engine, bookingRepository *repository.BookingRepository) {
	bc := bookingController{bookingRepository: bookingRepository}

	bookings := r.Group("/bookings")
	{
		bookings.POST("", bc.Create)
		bookings.GET("/:id", bc.Get)
		bookings.GET("/list", bc.List)
		bookings.PUT("/:id", bc.Update)
		bookings.DELETE("/:id", bc.Delete)
	}
}
