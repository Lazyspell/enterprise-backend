package checkapi

import (
	"github.com/lazyspell/enterprise-backend/foundation/web"
)

func Routes(app *web.App) {

	app.HandleFunc("/liveness", Liveness)
	app.HandleFunc("/readiness", Readiness)

}
