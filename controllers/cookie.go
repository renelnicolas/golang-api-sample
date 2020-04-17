package controllers

import "net/http"

// CookieController :
type CookieController struct {
}

// https://gist.github.com/liubin/5705616
// https://www.socketloop.com/tutorials/golang-drop-cookie-to-visitor-s-browser-and-http-setcookie-example
// data, err := base64.StdEncoding.DecodeString(cookie.Value)

// Create :
func (h CookieController) Create(w http.ResponseWriter, r *http.Request) {
	c := http.Cookie{
		Name:   "ithinkidroppedacookie",
		Value:  "thedroppedcookiehasgoldinit",
		Path:   "/",
		MaxAge: 3600,
	}
	http.SetCookie(w, &c)

	w.Write([]byte("new cookie created!\n"))
}

// Read :
func (h CookieController) Read(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("ithinkidroppedacookie")
	if err != nil {
		w.Write([]byte("error in reading cookie : " + err.Error() + "\n"))
	} else {
		value := c.Value
		w.Write([]byte("cookie has : " + value + "\n"))
	}
}

// Delete :
// see https://golang.org/pkg/net/http/#Cookie
// Setting MaxAge<0 means delete cookie now.
func (h CookieController) Delete(w http.ResponseWriter, r *http.Request) {
	c := http.Cookie{
		Name:   "ithinkidroppedacookie",
		MaxAge: -1,
	}
	http.SetCookie(w, &c)

	w.Write([]byte("old cookie deleted!\n"))
}
