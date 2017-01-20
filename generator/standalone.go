// +build !appengine

package generator

import (
	"net/http"
)

func client(r *http.Request) *http.Client {
	return &http.Client{}
}
