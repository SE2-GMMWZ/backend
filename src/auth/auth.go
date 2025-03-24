package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var jwtSecret = []byte("supersecretkey")

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateJWT(email string) (string, error) {
	claims := &Claims{
		Username: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := ""
		if cookie, err := c.Cookie("token"); err == nil {
			authHeader = cookie
		}

		if authHeader == "" {
			c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("/login?dest=%s", url.QueryEscape(c.Request.URL.Path)))
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("/login?dest=%s", url.QueryEscape(c.Request.URL.Path)))
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			c.Set("email", claims.Username)
		} else {
			c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("/login?dest=%s", url.QueryEscape(c.Request.URL.Path)))
			c.Abort()
			return
		}
		c.Next()
	}
}

func RedirectIfAuthenticated() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := ""
		if cookie, err := c.Cookie("token"); err == nil {
			authHeader = cookie
		}
		token, err := jwt.ParseWithClaims(authHeader, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err == nil && token.Valid {
			referrer := c.Request.Referer()
			if referrer == "" {
				referrer = "/"
			}
			c.Redirect(http.StatusTemporaryRedirect, referrer)
			c.Abort()
			return
		}

		c.Next()
	}
}
