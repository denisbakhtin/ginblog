package controllers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/denisbakhtin/ginblog/config"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// user is oauth user info struct
type user struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Profile       string `json:"profile"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Gender        string `json:"gender"`
}

// OauthCallback handles authentication of a user and initiates a session.
func OauthCallback(c *gin.Context) {
	// Handle the exchange code to initiate a transport.
	session := sessions.Default(c)
	returnURL := session.Get("oauth2-return").(string)
	retrievedState := session.Get("state").(string)
	queryState := c.Request.URL.Query().Get("state")
	if retrievedState != queryState {
		slog.Error("Invalid session state", "Retrieved", retrievedState, "Param", queryState)
		session.AddFlash(fmt.Sprintf("Invalid session state: retrieved: %s; Param: %s", retrievedState, queryState))
		session.Save()
		c.Redirect(http.StatusFound, returnURL)
		return
	}
	code := c.Request.URL.Query().Get("code")
	tok, err := googleConfig().Exchange(context.TODO(), code)
	if err != nil {
		slog.Error(err.Error())
		session.AddFlash("Authentication error, please try again")
		session.Save()
		c.Redirect(http.StatusFound, returnURL)
		return
	}

	client := googleConfig().Client(context.TODO(), tok)
	userinfo, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		slog.Error(err.Error())
		session.AddFlash("Authentication error, please try again")
		session.Save()
		c.Redirect(http.StatusFound, returnURL)
		return
	}
	defer userinfo.Body.Close()
	data, _ := io.ReadAll(userinfo.Body)
	u := user{}
	if err = json.Unmarshal(data, &u); err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500", gin.H{"Error": err})
		slog.Error(err.Error())
		return
	}
	session.Set("oauth-email", u.Email)
	session.Set("oauth-username", u.Name)
	session.Save()
	c.Redirect(http.StatusFound, returnURL)
}

// OauthGoogleLogin handles the google oauth login procedure.
func OauthGoogleLogin(c *gin.Context) {
	state := randToken(32)
	session := sessions.Default(c)
	session.Set("state", state)
	session.Set("oauth2-return", c.Request.Referer()+"#comments")
	session.Save()
	link := googleConfig().AuthCodeURL(state)
	c.Redirect(http.StatusSeeOther, link)
}

// randToken generates a random @l length token.
func randToken(l int) string {
	b := make([]byte, l)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func googleConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     config.GetConfig().Oauth.GoogleClientID,
		ClientSecret: config.GetConfig().Oauth.GoogleSecret,
		RedirectURL:  fmt.Sprintf("%s/oauthcallback", config.GetConfig().Domain),
		Scopes: []string{
			"email", // You have to select your own scope from here -> https://developers.google.com/identity/protocols/googlescopes#google_sign-in
			"profile",
		},
		Endpoint: google.Endpoint,
	}
}
