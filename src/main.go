package main

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
)

var tokens []string

func main() {
	router := gin.Default()

	router.GET("/login", gin.BasicAuth(gin.Accounts{
		"admin": "secret",
	}), loginHandler)

	router.GET("/resource", resourceHandler)

	router.Run(":8080")
}

func loginHandler(c *gin.Context) {
	token, _ := randomHex(20)
	tokens = append(tokens, token)

	c.JSON(200, gin.H{
		"token": token,
	})
}

func resourceHandler(c *gin.Context) {
	bearerToken := c.Request.Header.Get("Authorization")
	reqToken := strings.Split(bearerToken, " ")[1]
	if slices.Contains(tokens, reqToken) {
		c.JSON(200, gin.H{
			"success": true,
			"data":    "resource data",
		})
		return
	}

	c.JSON(http.StatusUnauthorized, gin.H{
		"success": false,
		"message": http.StatusText(http.StatusUnauthorized),
	})
}

func randomHex(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
