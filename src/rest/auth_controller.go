package rest

import (
	"backend/src/auth"
	"backend/src/model"
	"backend/src/repository"
	"backend/src/rest/shared"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
)

type authController struct {
	userRepository *repository.UserRepository
}

func setAuthCookie(c *gin.Context, token string, maxAge int) {
	expr := fmt.Sprintf(
		"token=%s; Max-Age=%d; Path=/; Secure; HttpOnly; SameSite=None",
		token,
		maxAge,
	)
	c.Writer.Header().Add("Set-Cookie", expr)
}

// GetUserByID godoc
// @Summary      Get user
// @Description  Returns a single user by its UUID
// @Tags         Users
// @Produce      json
// @Param        id   path      string  true  "User UUID"
// @Success      200  {object}  model.User
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /users/{id} [get]
func (ac *authController) GetUserByID(c *gin.Context) {
	userID := c.Param("id")
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format."})
		return
	}

	user, err := ac.userRepository.GetUserByID(userUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// UpdateUser godoc
// @Summary      Update user
// @Description  Updates a user by its UUID
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id    path      string  true  "User UUID"
// @Param        body  body      object  true  "User data to update"
// @Success      200   {object}  model.User
// @Failure      400   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /users/{id} [put]
func (ac *authController) UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format."})
		return
	}

	var updateData struct {
		Email   string         `json:"email"`
		Name    string         `json:"name"`
		Surname string         `json:"surname"`
		Phone   string         `json:"phone"`
		Role    model.UserRole `json:"role"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input."})
		return
	}

	user, err := ac.userRepository.GetUserByID(userUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found."})
		return
	}

	// Update the user details
	user.Email = updateData.Email
	user.Name = updateData.Name
	user.Surname = updateData.Surname
	user.PhoneNumber = &updateData.Phone
	user.Role = updateData.Role

	err = ac.userRepository.UpdateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully.", "user": user})
}

// DeleteUser godoc
// @Summary      Delete user
// @Description  Deletes a user by its UUID
// @Tags         Users
// @Produce      json
// @Param        id   path      string  true  "User UUID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /users/{id} [delete]
func (ac *authController) DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format."})
		return
	}

	err = ac.userRepository.DeleteUser(userUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully."})
}

// Signup godoc
// @Summary      Sign-up
// @Description  Creates a new user and sets JWT cookie
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      object  true  "User sign-up credentials"
// @Success      200   {object}  model.User
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /signup [post]
func (ac *authController) Signup(c *gin.Context) {
	var credentials struct {
		Email    string         `json:"email" binding:"required"`
		Password string         `json:"password" binding:"required"`
		Name     string         `json:"name" binding:"required"`
		Surname  string         `json:"surname" binding:"required"`
		Phone    string         `json:"phone" binding:"required"`
		Role     model.UserRole `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Password encryption failed."})
		return
	}

	err = ac.userRepository.CreateUser(&model.User{
		Email:       credentials.Email,
		Password:    string(hashedPassword),
		PhoneNumber: &credentials.Phone,
		Surname:     credentials.Surname,
		Name:        credentials.Name,
		Role:        credentials.Role,
	})

	if err != nil && strings.Contains(err.Error(), "already exists") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists."})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error."})
		return
	}

	user, err := ac.userRepository.GetUserByEmail(credentials.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error."})
		return
	}

	token, err := auth.GenerateJWT(user.UserID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token."})
		return
	}

	setAuthCookie(c, token, 3600*24*30)

	c.JSON(http.StatusOK, gin.H{"message": "Login successful.", "user": gin.H{
		"id":      user.UserID,
		"email":   user.Email,
		"name":    user.Name,
		"surname": user.Surname,
		"phone":   user.PhoneNumber,
		"role":    user.Role,
	}})
}

// Login godoc
// @Summary      Login
// @Description  Authenticates a user and returns JWT cookie
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      object  true  "User login credentials"
// @Success      200   {object}  model.User
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /login [post]
func (ac *authController) Login(c *gin.Context) {
	var credentials struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	user, err := ac.userRepository.GetUserByEmail(credentials.Email)
	if err != nil && strings.Contains(err.Error(), "not found") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password."})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error."})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password."})
		return
	}

	token, err := auth.GenerateJWT(user.UserID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token."})
		return
	}

	setAuthCookie(c, token, 3600*24*30)

	c.JSON(http.StatusOK, gin.H{"message": "Login successful.", "user": gin.H{
		"id":      user.UserID,
		"email":   user.Email,
		"name":    user.Name,
		"surname": user.Surname,
		"phone":   user.PhoneNumber,
		"role":    user.Role,
	}})
}

// List godoc
// @Summary      List users
// @Description  Returns paginated list of users; supports query parameters forwarded to repository
// @Tags         Users
// @Produce      json
// @Param        page   query     int false "Page number"
// @Param        limit  query     int false "Items per page"
// @Success      200    {object}  model.User
// @Router       /users/list [get]
func (ac *authController) List(c *gin.Context) {
	shared.ListWithQuery(c, ac.userRepository.ListUsers)
}

// Logout godoc
// @Summary      Logout
// @Description  Clears authentication cookie
// @Tags         Auth
// @Success      200  "Successfully logged out"
// @Router       /logout [post]
func (ac *authController) Logout(c *gin.Context) {
	expr := "token=; Max-Age=-1; Path=/; Secure; HttpOnly; SameSite=None"
	c.Writer.Header().Add("Set-Cookie", expr)
	c.Status(http.StatusOK)
}

// UserInfo godoc
// @Summary      Current user info
// @Description  Returns information about the currently authenticated user
// @Tags         Auth
// @Produce      json
// @Success      200  {object}  model.User
// @Failure      500  {object}  map[string]string
// @Security     ApiKeyAuth
// @Router       /user-info [get]
func (ac *authController) UserInfo(c *gin.Context) {
	userId, present := c.Get("userId")
	if !present {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error."})
		return
	}
	user, err := ac.userRepository.GetUserByID(userId.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": gin.H{
		"id":      user.UserID,
		"email":   user.Email,
		"name":    user.Name,
		"surname": user.Surname,
		"phone":   user.PhoneNumber,
		"role":    user.Role,
	}})
}

func AddAuthRoutes(r *gin.Engine, userRepository *repository.UserRepository) {
	ac := authController{userRepository: userRepository}

	r.POST("/signup", auth.RedirectIfAuthenticated(), ac.Signup)
	r.POST("/login", auth.RedirectIfAuthenticated(), ac.Login)
	r.GET("/user-info", auth.AuthMiddleware(), ac.UserInfo)
	r.POST("/logout", ac.Logout)

	r.GET("/users/list", ac.List)
	r.GET("/users/:id", ac.GetUserByID)
	r.PUT("/users/:id", ac.UpdateUser)
	r.DELETE("/users/:id", ac.DeleteUser)
}
