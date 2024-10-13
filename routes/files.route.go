package routes

import (
	"ssf/controllers"

	expressgo "github.com/lamgiahung112/express-go"
)

func InitFilesRoute(app *expressgo.Express) {
	filesController := controllers.NewFilesController()

	route := app.NewRoute("/files")
	route.Route("/store").Post(filesController.HandleSubmitShareFile)
	route.Route("/check-password").Post(filesController.CheckFilePassword)
	route.Route("/download/{slug}").Get(filesController.Download)
	route.Route("/{slug}").Get(filesController.HandleShowSharedFile)
}
