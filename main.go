package main

import (
	"flag"
	"github.com/alash3al/vecdb/internals/config"
	"github.com/alash3al/vecdb/internals/embeddings"
	"github.com/alash3al/vecdb/internals/http"
	"github.com/alash3al/vecdb/internals/store"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"log"
	// embedder drivers
	_ "github.com/alash3al/vecdb/internals/embeddings/drivers/gemini"

	// store drivers
	_ "github.com/alash3al/vecdb/internals/store/drivers/bolt"
)

var (
	flagConfigFilename = flag.String("config", "./config.yml", "the configuration filename")
)

func main() {
	flag.Parse()

	var cfg *config.Config
	var db store.Driver
	var embedder embeddings.Driver
	var err error

	{
		cfg, err = config.NewFromFile(*flagConfigFilename)
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	{
		db, err = store.Open(cfg.Store.Driver, cfg.Store.Args)
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	if cfg.Embedder.Enabled {
		embedder, err = embeddings.Open(cfg.Embedder.Driver, cfg.Embedder.Args)
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	validate := validator.New(validator.WithRequiredStructEnabled())

	app := fiber.New()

	app.Get("/", http.Home())

	app.Post("/v1/vectors/write", http.VectorWrite(validate, db))
	app.Post("/v1/vectors/search", http.VectorSearch(validate, db))

	app.Post("/v1/embeddings/text/write", http.EmbeddingsTextWrite(validate, embedder, db))
	app.Post("/v1/embeddings/text/search", http.EmbeddingsTextSearch(validate, embedder, db))

	log.Fatal(app.Listen(cfg.Server.ListenAddr))
}
