package rest

import (
	"backend/src/model"
	"backend/src/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type commentController struct {
	commentRepository *repository.CommentRepository
}

func (cc *commentController) Create(c *gin.Context) {
	var comment model.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	if err := cc.commentRepository.CreateComment(&comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create comment"})
		return
	}

	c.JSON(http.StatusCreated, comment)
}

func (cc *commentController) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
		return
	}

	comment, err := cc.commentRepository.GetCommentByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "comment not found"})
		return
	}

	c.JSON(http.StatusOK, comment)
}

func (cc *commentController) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
		return
	}

	var comment model.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	comment.CommentID = id
	if err := cc.commentRepository.UpdateComment(&comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update comment"})
		return
	}

	c.JSON(http.StatusOK, comment)
}

func (cc *commentController) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
		return
	}

	if err := cc.commentRepository.DeleteComment(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete comment"})
		return
	}

	c.Status(http.StatusNoContent)
}

func AddCommentRoutes(r *gin.Engine, commentRepository *repository.CommentRepository) {
	cc := commentController{commentRepository: commentRepository}

	comments := r.Group("/comments")
	{
		comments.POST("", cc.Create)
		comments.GET("/:id", cc.Get)
		comments.PUT("/:id", cc.Update)
		comments.DELETE("/:id", cc.Delete)
	}
}
