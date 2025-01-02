// Package mux provides support to bind domain level routes
// to the application mux
package mux

import (
	"os"

	"github.com/lazyspell/enterprise-backend/apis/services/api/mid"
	"github.com/lazyspell/enterprise-backend/apis/services/sales/route/sys/checkapi"
	"github.com/lazyspell/enterprise-backend/foundation/logger"
	"github.com/lazyspell/enterprise-backend/foundation/web"
)

// WebAPI constructs a http.Handler with all applicatoin routes bound.
func WebAPI(log *logger.Logger, shutdown chan os.Signal) *web.App {
	mux := web.NewApp(shutdown, mid.Logger(log))

	checkapi.Routes(mux)

	return mux
}
