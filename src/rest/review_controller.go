package rest

import (
	"backend/src/model"
	"backend/src/repository"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type reviewController struct {
	reviewRepository *repository.ReviewRepository
}

// Create godoc
// @Summary      Create review
// @Description  Creates a new review
// @Tags         Reviews
// @Accept       json
// @Produce      json
// @Param        body  body      model.Review  true  "Review to create"
// @Success      201   {object}  model.Review
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /reviews [post]
func (rc *reviewController) Create(c *gin.Context) {
	var review model.Review
	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	if err := rc.reviewRepository.CreateReview(&review); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create review"})
		return
	}

	c.JSON(http.StatusCreated, review)
}

// Get godoc
// @Summary      Get review
// @Description  Returns a single review by its UUID
// @Tags         Reviews
// @Produce      json
// @Param        id   path      string        true  "Review UUID"
// @Success      200  {object}  model.Review
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /reviews/{id} [get]
func (rc *reviewController) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
		return
	}

	review, err := rc.reviewRepository.GetReviewByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "review not found"})
		return
	}

	c.JSON(http.StatusOK, review)
}

// Update godoc
// @Summary      Update review
// @Description  Updates a review by its UUID
// @Tags         Reviews
// @Accept       json
// @Produce      json
// @Param        id    path      string        true  "Review UUID"
// @Param        body  body      model.Review  true  "Review data to update"
// @Success      200   {object}  model.Review
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /reviews/{id} [put]
func (rc *reviewController) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
		return
	}

	var review model.Review
	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	review.ReviewID = id
	if err := rc.reviewRepository.UpdateReview(&review); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update review"})
		return
	}

	c.JSON(http.StatusOK, review)
}

// Delete godoc
// @Summary      Delete review
// @Description  Deletes a review by its UUID
// @Tags         Reviews
// @Produce      json
// @Param        id   path      string  true  "Review UUID"
// @Success      204  "No Content"
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /reviews/{id} [delete]
func (rc *reviewController) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
		return
	}

	if err := rc.reviewRepository.DeleteReview(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete review"})
		return
	}

	c.Status(http.StatusNoContent)
}

// List godoc
// @Summary      List reviews
// @Description  Returns paginated list of reviews; supports query parameters forwarded to repository
// @Tags         Reviews
// @Produce      json
// @Param        reviewer_id     query     string  false "Filter by Reviewer UUID"
// @Param        rating          query     number  false "Minimum rating"
// @Param        comment         query     string  false "Text search in comment"
// @Param        docking_spot_id query     string  false "Filter by Docking Spot UUID"
// @Param        limit           query     int     false "Items per page"
// @Param        page            query     int     false "Page number"
// @Success      200    {object}  []model.Review
// @Router       /reviews/list [get]
func (rc *reviewController) List(c *gin.Context) {
	reviewerID := c.Query("reviewer_id")
	ratingStr := c.Query("rating")
	comment := c.Query("comment")
	dockingSpotID := c.Query("docking_spot_id")
	limitStr := c.DefaultQuery("limit", "10")

	rating := 0.0
	if ratingStr != "" {
		rating, _ = strconv.ParseFloat(ratingStr, 64)
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	var ratingPtr *float64
	if ratingStr != "" {
		ratingPtr = &rating
	}

	reviews, currentPage, totalPages, err := rc.reviewRepository.ListReviews(limit, page, reviewerID, ratingPtr, comment, dockingSpotID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list reviews"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"reviews":     reviews,
		"current_page": currentPage,
		"total_pages":  totalPages,
	})
}

func AddReviewRoutes(r *gin.Engine, reviewRepository *repository.ReviewRepository) {
	rc := reviewController{reviewRepository: reviewRepository}

	reviews := r.Group("/reviews")
	{
		reviews.POST("", rc.Create)
		reviews.GET("/:id", rc.Get)
		reviews.GET("/list", rc.List)
		reviews.PUT("/:id", rc.Update)
		reviews.DELETE("/:id", rc.Delete)
	}
}
