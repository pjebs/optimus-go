// +build appengine

package optimus

import (
	"net/http"
)

func client() *http.Client {

	transport := http.Transport{}

	return &http.Client{Transport: &transport}
}
