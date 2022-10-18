package celeritas

import (
	"net/http"
	"strconv"
)
import "github.com/justinas/nosurf"

func (c *Celeritas) SessionLoad(next http.Handler) http.Handler {
	return c.Session.LoadAndSave(next)
}

func (c *Celeritas) NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	secure, _ := strconv.ParseBool(c.config.cookie.secure)

	//exemption:
	csrfHandler.ExemptGlob("/api/*")

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   secure, // Only accept https ?
		SameSite: http.SameSiteStrictMode,
		Domain:   c.config.cookie.domain,
	})
	return csrfHandler
}
