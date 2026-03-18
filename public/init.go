package public

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func MountRoutes(router *fiber.App) {
	router.Get("/", func(c *fiber.Ctx) error {
		log.Println("OK")
		return c.SendString("OK")
	})

	userGroup := router.Group("/customers")
	{
		userGroup.Post("/", parseBody(createUserRequest{}), createUser)
		userGroup.Get("/", parseQuery(fetchUserByUserIdQuery{}), fetchUserByUserId)
		userGroup.Get("/:id", parseParam(fetchUserByIdParam{}), fetchUserById)
	}

}
