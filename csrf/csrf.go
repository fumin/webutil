package csrf

import (
	"bytes"
	"crypto/rand"
	"html/template"
	"net/http"

	"github.com/fumin/webutil"
)

const csrfCookieName = "_csrf"

func Check(r *http.Request) bool {
	token := r.FormValue(csrfCookieName)
	if token == "" {
		token = r.Header.Get("X-CSRFToken")
	}

	if token == "" || token != tokenFromCookie(r) {
		return false
	}
	return true
}

const hiddenInput = `<input name="{{.Name}}" value="{{.Value}}" type="hidden">`

var hiddenInputTmpl = template.Must(template.New("hiddenInput").Parse(hiddenInput))

func FormInput(w http.ResponseWriter, r *http.Request) template.HTML {
	token := setToken(w, r)
	var b bytes.Buffer
	data := struct {
		Name  string
		Value string
	}{
		Name:  csrfCookieName,
		Value: token,
	}
	err := hiddenInputTmpl.Execute(&b, data)
	if err != nil {
		panic(err)
	}
	return template.HTML(b.String())
}

func setToken(w http.ResponseWriter, r *http.Request) string {
	token := tokenFromCookie(r)
	if token != "" {
		return token
	}
	entropy := make([]byte, 16)
	rand.Read(entropy)
	token = webutil.Base64EncodeWithoutEq(entropy)
	http.SetCookie(w, &http.Cookie{Name: csrfCookieName, Value: token})
	return token
}

func tokenFromCookie(r *http.Request) string {
	cookie, err := r.Cookie(csrfCookieName)
	if err != nil {
		return ""
	}
	return cookie.Value
}
