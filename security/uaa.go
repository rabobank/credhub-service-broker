package security

import (
	"net/http"
)

func UAA(*http.Request) (bool, *string) { return true, nil }
