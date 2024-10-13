package routes

import (
	"ssf/controllers"

	expressgo "github.com/lamgiahung112/express-go"
)

func initHomeRoutes(app *expressgo.Express) {
	homeController := controllers.NewHomeController()

	route := app.NewRoute("/")

	route.Get(homeController.HandleHomePage)
}
