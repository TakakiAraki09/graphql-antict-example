package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"

	"github.com/TakakiAraki/graphql-antict-example/graph"
	"github.com/TakakiAraki/graphql-antict-example/graph/services"
	"github.com/TakakiAraki/graphql-antict-example/internal"

	_ "github.com/mattn/go-sqlite3"
)

const (
	defaultPort = "8080"
	dbFile      = "./mygraphql.db"
)

func main() {
  fmt.Println("initializing enviroment")
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
  fmt.Println("initialized enviroment")

  fmt.Println("connecting database")
	db, err := sql.Open("sqlite3", fmt.Sprintf("%s?_foreign_keys=on", dbFile))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
  fmt.Println("connected database")

  fmt.Println("register service")
	service := services.New(db)
	srv := handler.NewDefaultServer(internal.NewExecutableSchema(internal.Config{Resolvers: &graph.Resolver{
    Srv: service,
  }}))
  fmt.Println("registed service")

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

}
