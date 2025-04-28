package rest

import (
	"backend/src/auth"
	"backend/src/model"
	"backend/src/repository"
	"backend/src/rest/shared"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
)

type authController struct {
	userRepository *repository.UserRepository
}

func (am *authController) Signup(c *gin.Context) {
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

	err = am.userRepository.CreateUser(&model.User{
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

	user, err := am.userRepository.GetUserByEmail(credentials.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error."})
	}

	token, err := auth.GenerateJWT(user.UserID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token."})
		return
	}
	c.SetCookie("token", token, 3600*24*30 /* one month */, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Login successful.", "user": gin.H{
		"id":      user.UserID,
		"email":   user.Email,
		"name":    user.Name,
		"surname": user.Surname,
		"phone":   user.PhoneNumber,
		"role":    user.Role,
	}})
}

func (am *authController) Login(c *gin.Context) {
	var credentials struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	user, err := am.userRepository.GetUserByEmail(credentials.Email)
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

	c.SetCookie("token", token, 3600*24*30 /* one month */, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Login successful.", "user": gin.H{
		"id":      user.UserID,
		"email":   user.Email,
		"name":    user.Name,
		"surname": user.Surname,
		"phone":   user.PhoneNumber,
		"role":    user.Role,
	}})
}

func (ac *authController) List(c *gin.Context) {
	shared.ListWithQuery(c, ac.userRepository.ListUsers)
}

func (am *authController) Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "localhost", false, true)
	c.Status(http.StatusOK)
}

func (am *authController) UserInfo(c *gin.Context) {
	userId, present := c.Get("userId")
	if !present {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error."})
	}
	user, err := am.userRepository.GetUserByID(userId.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error."})
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
	am := authController{userRepository: userRepository}

	r.POST("/signup", auth.RedirectIfAuthenticated(), am.Signup)
	r.POST("/login", auth.RedirectIfAuthenticated(), am.Login)
	r.GET("/user-info", auth.AuthMiddleware(), am.UserInfo)
	r.POST("/logout", am.Logout)
}
