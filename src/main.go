package main

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var (
	tokens []string
	jwtKey = []byte("my_secret_2026")
)

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func main() {
	router := gin.Default()

	router.GET("/login", gin.BasicAuth(gin.Accounts{
		"admin": "secret",
	}), loginHandler)

	router.GET("/resource", resourceHandler)

	router.Run(":8080")
}

func loginHandler(c *gin.Context) {
	token, _ := generateJWT()
	tokens = append(tokens, token)

	c.JSON(200, gin.H{
		"token": token,
	})
}

func resourceHandler(c *gin.Context) {
	jwtToken := c.Request.Header.Get("Authorization")
	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(jwtToken, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Signature invalid",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad Request",
		})
	}

	if !tkn.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": http.StatusText(http.StatusUnauthorized),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    "resource data",
	})
}

func generateJWT() (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: "username",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func randomHex(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
