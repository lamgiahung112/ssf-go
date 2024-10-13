package controllers

import (
	"net/http"
	views "ssf/views/helpers"

	expressgo "github.com/lamgiahung112/express-go"
)

type HomeController struct{}

func NewHomeController() *HomeController {
	return &HomeController{}
}

func (controller *HomeController) HandleHomePage(w http.ResponseWriter, r *http.Request, context *expressgo.RequestContext) {
	context.Set("title_for_layout", "Hi mom")

	views.Render("index.html", context, w, r)
}
