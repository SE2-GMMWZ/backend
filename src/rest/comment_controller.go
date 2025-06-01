package rest

import (
	"backend/src/model"
	"backend/src/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strconv"
)

type commentController struct {
	commentRepository *repository.CommentRepository
}

// Create godoc
// @Summary      Create comment
// @Description  Creates a new comment
// @Tags         Comments
// @Accept       json
// @Produce      json
// @Param        body  body      model.Comment  true  "Comment to create"
// @Success      201   {object}  model.Comment
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /comments [post]
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

// Get godoc
// @Summary      Get comment
// @Description  Returns a single comment by its UUID
// @Tags         Comments
// @Produce      json
// @Param        id   path      string        true  "Comment UUID"
// @Success      200  {object}  model.Comment
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /comments/{id} [get]
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

// Update godoc
// @Summary      Update comment
// @Description  Updates a comment by its UUID
// @Tags         Comments
// @Accept       json
// @Produce      json
// @Param        id    path      string        true  "Comment UUID"
// @Param        body  body      model.Comment true  "Comment data to update"
// @Success      200   {object}  model.Comment
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /comments/{id} [put]
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

// Delete godoc
// @Summary      Delete comment
// @Description  Deletes a comment by its UUID
// @Tags         Comments
// @Produce      json
// @Param        id   path      string  true  "Comment UUID"
// @Success      204  "No Content"
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /comments/{id} [delete]
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

// List godoc
// @Summary      List comments
// @Description  Returns paginated list of comments; supports query parameters forwarded to repository
// @Tags         Comments
// @Produce      json
// @Param        guide_id  query     string false "Filter by Guide UUID"
// @Param        user_id   query     string false "Filter by User UUID"
// @Param        content   query     string false "Filter by content (substring match)"
// @Param        limit     query     int    false "Items per page"
// @Param        page      query     int    false "Page number"
// @Success      200    {object}  []model.Comment
// @Router       /comments/list [get]
func (cc *commentController) List(c *gin.Context) {
	guideID := c.Query("guide_id")
	userID := c.Query("user_id")
	content := c.Query("content")
	limitStr := c.DefaultQuery("limit", "10")
	pageStr := c.DefaultQuery("page", "1")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	comments, currentPage, totalPages, err := cc.commentRepository.ListComments(limit, page, guideID, userID, content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list comments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"comments":     comments,
		"current_page": currentPage,
		"total_pages":  totalPages,
	})
}

func AddCommentRoutes(r *gin.Engine, commentRepository *repository.CommentRepository) {
	cc := commentController{commentRepository: commentRepository}

	comments := r.Group("/comments")
	{
		comments.POST("", cc.Create)
		comments.GET("/:id", cc.Get)
		comments.GET("/list", cc.List)
		comments.PUT("/:id", cc.Update)
		comments.DELETE("/:id", cc.Delete)
	}
}
