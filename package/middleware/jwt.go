package middleware

import (
	"fmt"
	"net/http"
	"project1/package/initializer"
	"project1/package/models"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var SecretKey = []byte("qwertyuiop")
var BlacklistedTokens = make(map[string]bool)
var UserEmail string
var UserData models.Users

type Claims struct {
	Email string `json:"username"`
	Role  string `json:"roles"`
	jwt.StandardClaims
}

func JwtTokenStart(c *gin.Context, email string, role string) {
	tokenString, err := createToken(email, role)
	if err != nil {
		c.JSON(401, gin.H{
			"Error": "Failed to create Token",
		})
	}
	c.Set("token", tokenString)
	c.JSON(201, gin.H{
		"Token": tokenString,
	})
	fmt.Println("---------------===  ", tokenString, "  ===-----------------")
}

func createToken(email string, role string) (string, error) {
	claims := Claims{
		Email: email,
		Role:  role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 4).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func AuthMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		UserData = models.Users{}
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token not provided"})
			c.Abort()
			return
		}
		if BlacklistedTokens[tokenString] {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token revoked"})
			c.Abort()
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return SecretKey, nil
		})
		fmt.Println("token_-----  ", token)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		if claims.Role == "user" {
			fmt.Println("email  -----  ", claims.Email)
			if err := initializer.DB.First(&UserData, "email=?", claims.Email).Error; err != nil {
				c.JSON(400, gin.H{
					"error": "failed fetch user details",
				})
				c.Abort()
				return
			}
		}
		if claims.Role != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}
		c.Set("claims", claims)
		claims = &Claims{}
		c.Next()
	}
}
