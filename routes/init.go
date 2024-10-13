package routes

import expressgo "github.com/lamgiahung112/express-go"

func InitRoutes(app *expressgo.Express) {
	initHomeRoutes(app)
	InitFilesRoute(app)
}
