package rest

import (
	"backend/src/auth"
	"backend/src/model"
	"backend/src/repository"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type authManager struct {
	userRepository *repository.UserRepository
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(result)
}

func (am *authManager) Verify(c *gin.Context) {
	var credentials struct {
		Email string `form:"email" binding:"required"`
		Code  string `form:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload."})
		return
	}

	user, err := am.userRepository.GetUserByEmail(credentials.Email)
	if err != nil && strings.Contains(err.Error(), "not found") {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found."})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error."})
		return
	}

	user.Verified = true
	err = am.userRepository.UpdateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error."})
		return
	}
	c.Status(http.StatusOK)
	return
}

func (am *authManager) Signup(c *gin.Context) {
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
	credentials.Password = string(hashedPassword)

	code := generateRandomString(9)

	err = am.userRepository.CreateUser(&model.User{
		Email:            credentials.Email,
		Password:         credentials.Password,
		VerificationCode: code,
		Verified:         false,
	})

	if err != nil && strings.Contains(err.Error(), "already exists") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists."})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error."})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully.", "code": code})
}

func (am *authManager) Login(c *gin.Context) {
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
		c.JSON(500, gin.H{"error": "Internal server error."})
		return
	}

	if !user.Verified {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email not verified. Please verify your email."})
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

func (am *authManager) Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "localhost", false, true)
	c.Status(http.StatusOK)
}

func AddAuthRoutes(r *gin.Engine, userRepository *repository.UserRepository) {
	am := authManager{userRepository: userRepository}

	r.POST("/signup", auth.RedirectIfAuthenticated(), am.Signup)
	r.POST("/login", auth.RedirectIfAuthenticated(), am.Login)
	r.POST("/verify", auth.RedirectIfAuthenticated(), am.Verify)
	r.POST("/logout", am.Logout)
}
