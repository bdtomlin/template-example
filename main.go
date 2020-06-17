package main

import (
	"github.com/gofiber/fiber"
	"github.com/gofiber/template/pug"
)

func main() {
	engine := pug.New("./views", ".pug")
	engine.Reload(true)

	app := fiber.New(&fiber.Settings{
		Views: engine,
	})

	app.Get("/", func(c *fiber.Ctx) {
		c.Render("index", nil, "layouts/main")
	})

	app.Listen(3000)
}
