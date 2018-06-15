package controllers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/denisbakhtin/ginblog/system"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

//user is oauth user info struct
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

//OauthCallback handles authentication of a user and initiates a session.
func OauthCallback(c *gin.Context) {
	// Handle the exchange code to initiate a transport.
	session := sessions.Default(c)
	returnURL := session.Get("oauth2-return").(string)
	retrievedState := session.Get("state")
	queryState := c.Request.URL.Query().Get("state")
	if retrievedState != queryState {
		logrus.Errorf("Invalid session state: retrieved: %s; Param: %s", retrievedState, queryState)
		session.AddFlash(fmt.Sprintf("Invalid session state: retrieved: %s; Param: %s", retrievedState, queryState))
		session.Save()
		c.Redirect(http.StatusFound, returnURL)
		return
	}
	code := c.Request.URL.Query().Get("code")
	tok, err := googleConfig().Exchange(oauth2.NoContext, code)
	if err != nil {
		logrus.Error(err)
		session.AddFlash("Authentication error, please try again")
		session.Save()
		c.Redirect(http.StatusFound, returnURL)
		return
	}

	client := googleConfig().Client(oauth2.NoContext, tok)
	userinfo, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		logrus.Error(err)
		session.AddFlash("Authentication error, please try again")
		session.Save()
		c.Redirect(http.StatusFound, returnURL)
		return
	}
	defer userinfo.Body.Close()
	data, _ := ioutil.ReadAll(userinfo.Body)
	u := user{}
	if err = json.Unmarshal(data, &u); err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500", gin.H{"Error": err})
		logrus.Error(err)
		return
	}
	session.Set("oauth-email", u.Email)
	session.Set("oauth-username", u.Name)
	session.Save()
	c.Redirect(http.StatusFound, returnURL)
}

//OauthGoogleLogin handles the google oauth login procedure.
func OauthGoogleLogin(c *gin.Context) {
	state := randToken(32)
	session := sessions.Default(c)
	session.Set("state", state)
	session.Set("oauth2-return", c.Request.Referer()+"#comments")
	session.Save()
	link := googleConfig().AuthCodeURL(state)
	c.Redirect(http.StatusSeeOther, link)
}

//randToken generates a random @l length token.
func randToken(l int) string {
	b := make([]byte, l)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func googleConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     system.GetConfig().Oauth.GoogleClientID,
		ClientSecret: system.GetConfig().Oauth.GoogleSecret,
		RedirectURL:  "http://localhost:8080/oauthcallback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email", // You have to select your own scope from here -> https://developers.google.com/identity/protocols/googlescopes#google_sign-in
		},
		Endpoint: google.Endpoint,
	}
}
