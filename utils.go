package webutil

import (
	"encoding/base64"
	"math/rand"
	"net/http"
	"time"
)

const (
	alphabet     = "-_0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	alphabetSize = len(alphabet)
)

var randSource *rand.Rand

func init() {
	randSource = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func RandByteSliceWithSize(size int) []byte {
	b := make([]byte, size)
	for i := 0; i != size; i++ {
		b[i] = alphabet[randSource.Intn(alphabetSize)]
	}
	return b
}

// The hash space size of this function is close to a hex string of size 32.
// According to http://en.wikipedia.org/wiki/Birthday_problem#Probability_table
// this function has a collision probability of 10^-12 in 2.6*10^13 elements.
func RandByteSlice() []byte {
	return RandByteSliceWithSize(8)
}

func Base64EncodeWithoutEq(data []byte) string {
	str := base64.URLEncoding.EncodeToString(data)
	equalNum := 0
	for i := len(str) - 1; i >= 0; i-- {
		if str[i] == '=' {
			equalNum++
		} else {
			break
		}
	}
	return str[0 : len(str)-equalNum]
}

func Base64DecodeWithoutEq(str string) ([]byte, error) {
	extras := len(str) - (len(str)/4)*4
	if extras == 0 {
		return base64.URLEncoding.DecodeString(str)
	}

	input := []byte(str)
	for i := 0; i < 4-extras; i++ {
		input = append(input, '=')
	}
	return base64.URLEncoding.DecodeString(string(input))
}

// ExpiresHeader adds an "Expires" in the http response.
// For example, we can use it to set an Expires of 360 days for static content:
// http.Handle("/static/", expiresHeader(360*24*time.Hour, http.StripPrefix("/static/", http.FileServer(http.Dir("prod/static")))))
func ExpiresHeader(d time.Duration, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Expires", time.Now().Add(d).Format(time.RFC1123))
		h.ServeHTTP(w, r)
	})
}
