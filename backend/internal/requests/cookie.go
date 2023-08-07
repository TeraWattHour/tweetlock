package requests

import (
	"net/http"
	"os"
	"time"
)

func GetCookieValue(cookieHeader string, cookieName string) string {
	cookies := http.Header{}
	cookies.Add("Cookie", cookieHeader)
	request := http.Request{Header: cookies}

	cookie, err := request.Cookie(cookieName)
	if err != nil {
		return ""
	}

	return cookie.Value
}

func GenerateCookie(name string, value string, expires time.Time) string {
	cookieDomain := os.Getenv("COOKIE_DOMAIN")

	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Domain:   cookieDomain,
		Expires:  expires,
		Secure:   cookieDomain != "localhost",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	return cookie.String()
}

func AccessCookie(value string, expires time.Time) string {
	return GenerateCookie("x-access", value, expires)
}

func RefreshCookie(value string, expires time.Time) string {
	return GenerateCookie("x-refresh", value, expires)
}
