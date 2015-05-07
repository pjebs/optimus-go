// +build !appengine

package optimus

import (
	"net/http"
)

func client() *http.Client {
	return &http.Client{}
}
