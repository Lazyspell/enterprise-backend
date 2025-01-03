package graphqlapi

import (
	"context"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"

	// "github.com/lazyspell/enterprise-backend/foundation/web"
	// "github.com/99designs/gqlgen/graphql/playground"
	"github.com/lazyspell/enterprise-backend/apis/graphql/graph"
	"github.com/lazyspell/enterprise-backend/foundation/web"
	"github.com/vektah/gqlparser/v2/ast"
)

func Graphql(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	return web.GraphqlResponse(ctx, w, r, srv, 200)
}
