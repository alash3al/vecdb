package http

import (
	"github.com/alash3al/vecdb/internals/embeddings"
	"github.com/alash3al/vecdb/internals/store"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func EmbeddingsTextWrite(validate *validator.Validate, embedder embeddings.Driver, db store.Driver) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		if embedder == nil {
			return ctx.SendStatus(http.StatusNotImplemented)
		}

		var input struct {
			Bucket  string `json:"bucket" validate:"required"`
			Key     string `json:"key" validate:"required"`
			Content string `json:"content" validate:"required"`
		}

		if err := ctx.BodyParser(&input); err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
				"hint":  `seems, you didn't specify a valid json body {"bucket": "<String>", "key": "<String>", "content": "<String>"}`,
			})
		}

		if err := validate.Struct(&input); err != nil {
			return ctx.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
				"error": err.Error(),
				"hint":  `the required JSON body fields are {"bucket": "<String>", "key": "<String>", "content": "<String>"}`,
			})
		}

		vec, err := embedder.TextEmbedding(ctx.Context(), input.Content)
		if err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
				"hint":  `I was unable to perform TextEmbedding(), kindly check the embedder configurations`,
			})
		}

		if err := db.Put(input.Bucket, input.Key, vec); err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
				"hint":  `I was unable to write to the data store`,
			})
		}

		return ctx.SendStatus(http.StatusOK)
	}
}

func EmbeddingsTextSearch(validate *validator.Validate, embedder embeddings.Driver, db store.Driver) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		if embedder == nil {
			return ctx.SendStatus(http.StatusNotImplemented)
		}

		var input struct {
			Bucket              string   `json:"bucket" validate:"required"`
			Content             string   `json:"content" validate:"required"`
			MinCosineSimilarity *float64 `json:"min_cosine_similarity" validate:"required"`
			MaxResultCount      *int64   `json:"max_result_count" validate:"required"`
		}

		if err := ctx.BodyParser(&input); err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
				"hint":  `seems, you didn't specify a valid json body {"bucket": "<String>", "content": "<String>", "min_cosine_similarity": <float>, "max_result_count": <int>}`,
			})
		}

		if err := validate.Struct(&input); err != nil {
			return ctx.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
				"error": err.Error(),
				"hint":  `the required JSON body fields are {"bucket": "<String>", "content": "<String>", "min_cosine_similarity": <float>, "max_result_count": <int>}`,
			})
		}

		vec, err := embedder.TextEmbedding(ctx.Context(), input.Content)
		if err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
				"hint":  `I was unable to perform TextEmbedding(), kindly check the embedder configurations`,
			})
		}

		result, err := db.Query(store.VectorQueryInput{
			Bucket:              input.Bucket,
			Vector:              vec,
			MinCosineSimilarity: *input.MinCosineSimilarity,
			MaxResultCount:      *input.MaxResultCount,
		})

		if err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
				"hint":  "I was unable to find the data in the data store",
			})
		}

		if result.Items == nil {
			result.Items = []store.VectorQueryResultItem{}
		}

		return ctx.Status(http.StatusOK).JSON(fiber.Map{
			"result": result,
		})
	}
}
