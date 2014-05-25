package main

import (
	"github.com/go4r/handy"
	"github.com/go4r/handy/example/app"
	"net/http"
)

func main() {

	//http.DefaultServeMux
	var (
		mainController = &app.IndexController{}
	)


	println("Starting Server")

	handy.Server.Map(
	mainController,
	)

	handy.Server.Get("/public/(*:filepath)",
	func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/"+handy.StringParameter(r, "filepath"))
	})

	http.ListenAndServe(":8080", handy.Server)
}
