package main

import (
	"context"
	"errors"
	"expvar"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/ardanlabs/conf/v3"
	"github.com/lazyspell/enterprise-backend/apis/graphql/mux"
	// "github.com/lazyspell/enterprise-backend/apis/services/api/debug"
	"github.com/lazyspell/enterprise-backend/foundation/logger"
	"github.com/lazyspell/enterprise-backend/foundation/web"
)

const build = "develop"

func main() {
	var log *logger.Logger
	events := logger.Events{
		Error: func(ctx context.Context, r logger.Record) {
			log.Info(ctx, "******* SEND ALERT *******")
		},
	}

	traceIDFn := func(ctx context.Context) string {
		return web.GetTraceID(ctx)
	}

	log = logger.NewWithEvents(os.Stdout, logger.LevelInfo, "GRAPHQL", traceIDFn, events)
	//----------------------------------------------------------------------------------------
	ctx := context.Background()

	if err := run(ctx, log); err != nil {
		log.Error(ctx, "startup", "msg", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, log *logger.Logger) error {
	//----------------------------------------------------------------------------------------
	// GOMAXPROCS

	log.Info(ctx, "startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	cfg := struct {
		conf.Version
		Web struct {
			ReadTimeout        time.Duration `conf:"default:5s"`
			WriteTimeout       time.Duration `conf:"default:10s"`
			IdleTimeout        time.Duration `conf:"default:120s"`
			ShutdownTimeout    time.Duration `conf:"default:20s"`
			APIHost            string        `conf:"default:0.0.0.0:8080"`
			DebugHost          string        `conf:"default:0.0.0.0:8081"`
			CORSAllowedOrigins []string      `conf:"default:*,mask"`
		}
	}{
		Version: conf.Version{
			Build: build,
			Desc:  "Graphql",
		},
	}

	const prefix = "GRAPHQL"
	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}
		return fmt.Errorf("parsing config: %w", err)
	}

	//----------------------------------------------------------------------------------------
	// App Starting
	log.Info(ctx, "starting service", "version", cfg.Build)
	defer log.Info(ctx, "shutdown complete")

	out, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("generating config for output: %w", err)
	}
	log.Info(ctx, "startup", "config", out)

	expvar.NewString("build").Set(cfg.Build)

	//----------------------------------------------------------------------------------------

	// go func() {
	// 	log.Info(ctx, "startup", "status", "debug v1 router started", "host", cfg.Web.DebugHost)
	// 	if err := http.ListenAndServe(cfg.Web.DebugHost, debug.Mux()); err != nil {
	// 		log.Error(ctx, "shutdown", "status", "debug v1 router closed", "host", cfg.Web.DebugHost, "msg", err)
	// 	}
	// }()
	//----------------------------------------------------------------------------------------

	log.Info(ctx, "startup", "status", "initializing V1 API support")
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	api := http.Server{
		Addr:         cfg.Web.APIHost,
		Handler:      mux.GraphqlAPI(log, shutdown),
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
		ErrorLog:     logger.NewStdLogger(log, logger.LevelError),
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Info(ctx, "startup", "status", "api router started", "host", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	//----------------------------------------------------------------------------------------
	// Shutdown
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case sig := <-shutdown:
		log.Info(ctx, "shutdown", "status", "shutdown started", "signal", sig)
		defer log.Info(ctx, "shutdown", "status", "shutdown complete", "signal", sig)
		ctx, cancel := context.WithTimeout(ctx, cfg.Web.ShutdownTimeout)
		defer cancel()

		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}

	}
	return nil
}

// port := os.Getenv("PORT")
// if port == "" {
// 	port = defaultPort
// }
//
// srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))
//
// srv.AddTransport(transport.Options{})
// srv.AddTransport(transport.GET{})
// srv.AddTransport(transport.POST{})
//
// srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))
//
// srv.Use(extension.Introspection{})
// srv.Use(extension.AutomaticPersistedQuery{
// 	Cache: lru.New[string](100),
// })
//
// http.Handle("/", playground.Handler("GraphQL playground", "/query"))
// http.Handle("/query", srv)
//
// log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
// log.Fatal(http.ListenAndServe(":"+port, nil))
