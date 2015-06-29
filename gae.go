// +build appengine

package optimus

import (
	"appengine"
	"appengine/urlfetch"
	"net/http"
)

func client(r *http.Request) *http.Client {
	c := appengine.NewContext(r)
	return urlfetch.Client(c)
}
