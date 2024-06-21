package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var BlacklistedTokens = make(map[string]bool)


type Claims struct {
	Email  string `json:"username"`
	Role   string `json:"roles"`
	UserID uint
	jwt.StandardClaims
}

func JwtTokenStart(c *gin.Context, userId uint, email string, role string) string {
	tokenString, err := createToken(userId, email, role)
	if err != nil {
		c.JSON(500, gin.H{
			"status":"Fail",
			"Error": "Failed to create Token",
			"code": 500,
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
		tokenString, err := c.Cookie("jwtToken" + requiredRole)
		if err != nil {
			c.JSON(401, gin.H{
				"status": "Unauthorized",
				"message": "Can't find cookie",
				"code":    401,
			})
			c.Abort()
			return
		}
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRETKEY")), nil
		})
		if err != nil || !token.Valid {
			c.JSON(401, gin.H{
				"status":    "Unauthorized",
				"message":    "Invalid or expired JWT Token.",
				"code":      401,
			})
			c.Abort()
			return
		}
		if claims.Role != requiredRole {
			c.JSON(403, gin.H{
				"status":   "Forbidden",
				"error": "Insufficient permissions",
				"code":     403,
			})
			c.Abort()
			return
		}
		c.Set("userid", claims.UserID)
		c.Next()
	}
}
