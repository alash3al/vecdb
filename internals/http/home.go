package http

import "github.com/gofiber/fiber/v2"

func Home() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"message": "I'm Live!",
		})
	}
}
