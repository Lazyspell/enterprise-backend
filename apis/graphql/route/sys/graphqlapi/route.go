package graphqlapi

import (
	"github.com/lazyspell/enterprise-backend/foundation/web"
)

func Routes(app *web.App) {
	app.HandleFunc("/query", Graphql)
}
