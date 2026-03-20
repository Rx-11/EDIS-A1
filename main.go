package main

import (
	"log"
	"strings"

	"github.com/Rx-11/EDIS-A1/config"
	"github.com/Rx-11/EDIS-A1/db"
	"github.com/Rx-11/EDIS-A1/public"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New(fiber.Config{
		AppName: "BookStore-Backend",
	})

	app.Use(logger.New(logger.Config{Next: func(c *fiber.Ctx) bool {
		return strings.HasPrefix(c.Path(), "/")
	}}))

	config.Init()
	log.Println("Loaded configs.")
	db.Init(config.GetConfig().DbConfig, db.MySQL, db.LogInfo)
	db.Migrate()

	public.MountRoutes(app)

	log.Println("Server started at http://localhost:80")
	log.Fatal(app.Listen("0.0.0.0:80"))

}
