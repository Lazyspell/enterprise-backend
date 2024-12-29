// Package mux provides support to bind domain level routes
// to the application mux
package mux

import (
	"net/http"

	"github.com/lazyspell/enterprise-backend/apis/services/sales/route/sys/checkapi"
)

// WebAPI constructs a http.Handler with all applicatoin routes bound.
func WebAPI() http.Handler {
	mux := http.NewServeMux()

	checkapi.Routes(mux)

	return mux
}
