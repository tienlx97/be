package main

import (
	"be/auth"
	"be/graph"
	"be/logger"
	"be/repository"
	"context"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"

	_ "github.com/joho/godotenv/autoload"
)

const defaultPort = "8080"

func main() {

	ctx := context.Background()

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	r := chi.NewRouter()

	f, err := repository.FirebaseApp(ctx)
	if err != nil {
		panic(err)
	}

	r.Use(logger.Middleware())

	a := auth.New()
	r.Use(a.NotLoginMiddleware())
	r.Use(a.FirebaseLoginMiddleware(f))

	fc, err := f.Firestore(ctx)
	if err != nil {
		panic(err)
	}

	app := graph.NewApplication()

	r.Handle("/", playground.Handler("GraphQL Playground", "/query"))
	r.Handle("/query", srv)

	httpsrv := &http.Server{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      r,
		Addr:         ":" + port,
	}

	err = httpsrv.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
