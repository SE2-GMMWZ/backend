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
