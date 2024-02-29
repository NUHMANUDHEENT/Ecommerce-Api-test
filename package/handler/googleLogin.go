package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)
//================ get url and other details =================
func OauthSetup() *oauth2.Config {
	conf := &oauth2.Config{
		ClientID:     "1045963824475-g4k95hd785nt0ehes50ukaa8ufnql5j4.apps.googleusercontent.com",
		ClientSecret: "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",
		RedirectURL:  "http://localhost:8081/products",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
	return conf
}
// =================== check the authentication =============
func Googlelogin(c *gin.Context) {
	var googleConfig *oauth2.Config
	googleConfig = OauthSetup()

	code := c.Query("code")
	fmt.Println("---------------", code)
	token, err := googleConfig.Exchange(c, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": token.AccessToken})
}

