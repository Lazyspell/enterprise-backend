// Package mux provides support to bind domain level routes
// to the application mux
package mux

import (
	"os"

	"github.com/lazyspell/enterprise-backend/apis/services/sales/route/sys/checkapi"
	"github.com/lazyspell/enterprise-backend/foundation/web"
)

// WebAPI constructs a http.Handler with all applicatoin routes bound.
func WebAPI(shutdown chan os.Signal) *web.App {
	mux := web.NewApp(shutdown)

	checkapi.Routes(mux)

	return mux
}
