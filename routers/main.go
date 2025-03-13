package routers

import (
	"rokhelper/db"
	"rokhelper/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, DBMG *db.Mongo) {
	titleHandler := handlers.TitleHandler{DB: DBMG}
	KingdomHandler := handlers.KingdomHandler{DB: DBMG}
	api := app.Group("/api")

	title := api.Group("/title")
	title.Get("/:id", titleHandler.GetTitle)
	title.Post("/:id", titleHandler.AddTitle)
	title.Post("/finish/:id", titleHandler.FinishTitle)
	title.Post("/done/:id", titleHandler.DoneTitle)
	title.Post("/editMap/:id", titleHandler.EditMap)
	title.Post("/editConfig/:id", titleHandler.EditConfig)
	title.Get("/config/:id", titleHandler.GetConfig)

	kingdom := api.Group("/kingdom")

	kingdom.Post("/", KingdomHandler.CreateKingdom)
	kingdom.Get("/", KingdomHandler.GetKingdoms)
	kingdom.Get("/:id", KingdomHandler.GetKingdom)

}
