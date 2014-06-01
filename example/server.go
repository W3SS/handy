package main

import (
	. "github.com/go4r/handy"
	"github.com/go4r/handy/example/app"
	"net/http"
	"github.com/go4r/handy/lib"
	"time"
	"fmt"
)

func main() {

	//http.DefaultServeMux
	println("Starting Server")

	Server.SetMiddleware("authenticated", func(r *http.Request) bool {
		userId := lib.Float64Session(r, "userId")

		if userId != 0 {
			return true
		}

		RespondWithStatus(r, "You dont have access to this area", http.StatusUnauthorized)

		return false
	})

	Server.SetMiddleware("metrics", func(r *http.Request) {
		ctime := time.Now()
		lib.Defer(r, func() {
			ntime := time.Since(ctime)
			println("Elapsed Time", ntime.String())
		})
	})

	Server.Map(
	&app.IndexController{},
	)

	Server.HandleGet("/public/(*:filepath)", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/"+StringParameter(r, "filepath"))
	}).Before("authenticated").After(func(r *http.Request) {
		fmt.Println("Respond: " + r.URL.Path)
	})

	ListenAndServer(":8080")
}
