package mid

import (
	"context"
	"fmt"
	"time"

	"github.com/lazyspell/enterprise-backend/foundation/logger"
	"github.com/lazyspell/enterprise-backend/foundation/web"
)

func Logger(ctx context.Context, log *logger.Logger, path string, rawQuery string, method string, remoteAddr string, handler Handler) error {
	v := web.GetValues(ctx)

	if rawQuery != "" {
		path = fmt.Sprintf("%s?%s", path, rawQuery)
	}

	log.Info(ctx, "request started", "method", method, "path", path, "remoteaddr", remoteAddr)

	err := handler(ctx)
	log.Info(ctx, "request compelted", "method", method, "path", path, "remoteaddr", remoteAddr,
		"statuscode", v.StatusCode, "since", time.Since(v.Now).String())

	return err

}
