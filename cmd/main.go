package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	db "github.com/777Lava/ozonTest/internal/database"
	"github.com/777Lava/ozonTest/internal/graph"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("can't find .env file. continue work using system env vars from main")
	}

	var server *handler.Server
	fmt.Println(os.Getenv("MODE"))
	if os.Getenv("MODE") == "inmemory" {
		server = handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))
	} else {
		line := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"),
			os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))

		DB, err := gorm.Open(postgres.Open(line), &gorm.Config{})
		if err != nil {
			panic(err)
		}
		server = handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
			PostRepo:    db.PostRepo{DB: DB},
			CommentRepo: db.CommentRepo{DB: DB},
		}}))
	}

	server.AddTransport(&transport.Websocket{})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", server)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))

}
