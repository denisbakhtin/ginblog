package helpers

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

func GetSession(c *gin.Context) *sessions.Session {
	obj, exists := c.Get("Session")
	if !exists {
		panic("Context 'Session' not set. Use session middleware.")
	}
	session, ok := obj.(*sessions.Session)
	if !ok {
		panic("Context 'Session' should be *session.Session")
	}
	return session
}
