// File: config/config.go

package config

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

// Store is the session store used by the application.
var Store *sessions.CookieStore

func init() {
	Store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
	Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false, // Set to true in production when using HTTPS
	}

	log.Println("Session key: ", os.Getenv("SESSION_KEY"))

}
