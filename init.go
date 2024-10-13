package main

import (
	"os"
	"strconv"

	expressgo "github.com/lamgiahung112/express-go"
)

var App *expressgo.Express

func InitApp() {
	var port int
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		port = 3000
	}
	cfg := &expressgo.ExpressConfig{
		Host: "localhost",
		Port: uint16(port),
	}

	App = expressgo.NewExpress(cfg)
}

func StartApp() {
	App.StartServer()
}
