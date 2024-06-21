package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"project1/package/initializer"
	"project1/package/middleware"
	"project1/package/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	VerifiedEmail bool   `json:"verified_email"`
}

// ================ get url and other details =================
func OauthSetup() *oauth2.Config {
	conf := &oauth2.Config{
		ClientID:     "1045963824475-g4k95hd785nt0ehes50ukaa8ufnql5j4.apps.googleusercontent.com",
		ClientSecret: "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",
		RedirectURL:  "http://localhost:8080/auth/google/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
	return conf
}

// =================== check the authentication =============
// Googlelogin initiates the Google login process.
// @Summary Initiate Google login
// @Description Initiates the Google login process by redirecting to Google's OAuth authorization endpoint.
// @Tags auth
// @Produce html
// @Success 302 {string} string "Redirects to Google login page"
// @Router /auth/login [get]
func Googlelogin(c *gin.Context) {
	googleConfig := OauthSetup()
	url := googleConfig.AuthCodeURL("state")
	c.Redirect(http.StatusTemporaryRedirect, url)

}
func HandleGoogleCallback(c *gin.Context) {
	code := c.Request.URL.Query().Get("code")
	googleConfig := OauthSetup()
	token, err := googleConfig.Exchange(context.Background(), code)
	if err != nil {
		fmt.Println("Code not received properly. Please try again.")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	client := googleConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		log.Fatalf("GET /userinfo error %v", resp)
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	var googleUser GoogleUser
	err = json.NewDecoder(resp.Body).Decode(&googleUser)
	if err != nil {
		log.Fatalf("Couldn't decode response from Google: %v", err)
		return
	}
	user := models.Users{
		Name:  googleUser.Name,
		Email: googleUser.Email,
	}

	if googleUser.VerifiedEmail {
		if err := initializer.DB.First(&user, "email=?", user.Email).Error; err != nil {
			if err := initializer.DB.Create(&user).Error; err != nil {
				log.Fatal("Failed to create a new User")
				return
				
			} else {
				initializer.DB.First(&user, "email=?", user.Email)
			}
		}
		token := middleware.JwtTokenStart(c, user.ID, user.Email, user.Name)
		c.SetCookie("jwtTokenUser", token, 365, "/", "localhost", false, true)
		c.Redirect(http.StatusFound, "/")
	}
	fmt.Println("user ------------", user)
}
