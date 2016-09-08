// +build appengine

package ogen

import (
	"net/http"

	"appengine"
	"appengine/urlfetch"
)

func client(r *http.Request) *http.Client {
	c := appengine.NewContext(r)
	return urlfetch.Client(c)
}
