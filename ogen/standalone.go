// +build !appengine

package ogen

import (
	"net/http"
)

func client(r *http.Request) *http.Client {
	return &http.Client{}
}
