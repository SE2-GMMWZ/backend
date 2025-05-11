package rest

import (
	"backend/src/model"
	"backend/src/repository"
	"backend/src/rest/shared"
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
// @Param        page   query     int false "Page number"
// @Param        limit  query     int false "Items per page"
// @Success      200    {object}  map[string]interface{}
// @Router       /reviews/list [get]
func (rc *reviewController) List(c *gin.Context) {
	shared.ListEntities(c, rc.reviewRepository.ListReviews)
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
