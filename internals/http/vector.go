package http

import (
	"github.com/alash3al/vecdb/internals/store"
	"github.com/alash3al/vecdb/internals/vector"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func VectorWrite(validate *validator.Validate, db store.Driver) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var input struct {
			Bucket string     `json:"bucket" validate:"required"`
			Key    string     `json:"key" validate:"required"`
			Vector vector.Vec `json:"vector" validate:"required"`
		}

		if err := ctx.BodyParser(&input); err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
				"hint":  `seems, you didn't specify a valid json body {"bucket": "<String>", "key": "<String>", "vector": "List<float>"}`,
			})
		}

		if err := validate.Struct(&input); err != nil {
			return ctx.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
				"error": err.Error(),
				"hint":  `the required JSON body fields are {"bucket": "<String>", "key": "<String>", "vector": "List<float>"}`,
			})
		}

		if err := db.Put(input.Bucket, input.Key, input.Vector); err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
				"hint":  "I was unable to write the data into the data store",
			})
		}

		return ctx.SendStatus(http.StatusOK)
	}
}

func VectorSearch(validate *validator.Validate, db store.Driver) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var input struct {
			Bucket              string     `json:"bucket" validate:"required"`
			Vector              vector.Vec `json:"vector" validate:"required"`
			MinCosineSimilarity *float64   `json:"min_cosine_similarity" validate:"required"`
			MaxResultCount      *int64     `json:"max_result_count" validate:"required"`
		}

		if err := ctx.BodyParser(&input); err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
				"hint":  `seems, you didn't specify a valid json body {"bucket": "<String>", "vector": "List<float>", "min_cosine_similarity": <float>, "max_result_count": <int>}`,
			})
		}

		if err := validate.Struct(&input); err != nil {
			return ctx.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
				"error": err.Error(),
				"hint":  `the required JSON body fields are {"bucket": "<String>", "vector": "List<float>", "min_cosine_similarity": <float>, "max_result_count": <int>}`,
			})
		}

		result, err := db.Query(store.VectorQueryInput{
			Bucket:              input.Bucket,
			Vector:              input.Vector,
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
