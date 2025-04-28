package rest

import (
	"backend/src/auth"
	"backend/src/model"
	"backend/src/repository"
	"backend/src/rest/shared"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
)

type authController struct {
	userRepository *repository.UserRepository
}

func (am *authController) Signup(c *gin.Context) {
	var credentials struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
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
		Email:    credentials.Email,
		Password: string(hashedPassword),
	})

	if err != nil && strings.Contains(err.Error(), "already exists") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists."})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error."})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully."})
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

	token, err := auth.GenerateJWT(credentials.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token."})
		return
	}

	c.SetCookie("token", token, 3600*24*30 /* one month */, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Login successful."})
}

func (ac *authController) List(c *gin.Context) {
	shared.ListWithQuery(c, ac.userRepository.ListUsers)
}

func (am *authController) Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "localhost", false, true)
	c.Status(http.StatusOK)
}

func AddAuthRoutes(r *gin.Engine, userRepository *repository.UserRepository) {
	am := authController{userRepository: userRepository}

	r.POST("/signup", auth.RedirectIfAuthenticated(), am.Signup)
	r.POST("/login", auth.RedirectIfAuthenticated(), am.Login)
	r.POST("/logout", am.Logout)
}
