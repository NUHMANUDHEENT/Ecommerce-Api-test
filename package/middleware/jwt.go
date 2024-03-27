package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var BlacklistedTokens = make(map[string]bool)

// var UserData models.Users

type Claims struct {
	Email  string `json:"username"`
	Role   string `json:"roles"`
	UserID uint
	jwt.StandardClaims
}

func JwtTokenStart(c *gin.Context, userId uint, email string, role string) string {
	tokenString, err := createToken(userId, email, role)
	if err != nil {
		c.JSON(200, gin.H{
			"Error": "Failed to create Token",
		})
	}
	return tokenString
}

func createToken(userId uint, email string, role string) (string, error) {
	claims := Claims{
		Email:  email,
		Role:   role,
		UserID: uint(userId),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 60).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRETKEY")))
	if err != nil {
		fmt.Println("----", err, tokenString)
		return "", err
	}
	return tokenString, nil
}

func AuthMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("jwtToken")
		if err != nil {
			c.JSON(401, gin.H{
				"message": "Can't find cookie",
				"error":   err,
			})
			c.Abort()
			return
		}
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRETKEY")), nil
		})
		if err != nil || !token.Valid {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		if claims.Role != requiredRole {
			c.JSON(403, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}
		c.Set("userid", claims.UserID)
		c.Next()
	}
}
