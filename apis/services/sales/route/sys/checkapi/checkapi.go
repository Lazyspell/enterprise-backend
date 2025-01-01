package checkapi

import (
	"context"
	"net/http"

	"github.com/lazyspell/enterprise-backend/foundation/web"
)

func Liveness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	lineness := struct {
		Status string
	}{
		Status: "ready and alive",
	}

	return web.Respond(ctx, w, lineness, 200)

}

func Readiness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	readiness := struct {
		Status string
	}{
		Status: "status is ready",
	}

	return web.Respond(ctx, w, readiness, 200)

}
