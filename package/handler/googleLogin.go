package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)



func OauthSetup() *oauth2.Config {

    conf := &oauth2.Config{
        ClientID:     "726393225583-l7p1u41ve5mpc5ssi2kl09c9r4fhjeik.apps.googleusercontent.com",
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

func Googlelogin(c *gin.Context) {
    var googleConfig *oauth2.Config
    googleConfig = OauthSetup()
    url := googleConfig.AuthCodeURL("state")
    // c.Redirect(http.StatusFound, url)
    fmt.Println("check",url)

    // code := c.Query("code")
    // fmt.Println("=======>", code)
    token, err := googleConfig.Exchange(c, "authorization-code")
    fmt.Println("=====<>",token)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"access_token": token.AccessToken})
}

// func GoogleCallback(c *gin.Context) {
// }