package app

import (
	"fmt"
	. "github.com/go4r/handy"
	"net/http"
	"github.com/go4r/handy/lib"
)

type IndexController struct{ Annotations }

func (controller *IndexController) Annotate() {

	controller.
	Annotation(
		Before("metrics"),
		controller,
	).
	Annotation(
		Before("authenticated"),
		Name("index.getuser"),
		Any("/getuser/(number:userId)"),
		controller.getUsers,
	).
	Annotation(
		Name("index.index"),
		Any("/"),
		controller.indexAction,
	).
	Annotation(
		Name("index.setuser"),
		Any("/setuser/(:name)"),
		controller.setUser,
	)

}

func (controller *IndexController) setUser(w http.ResponseWriter, r *http.Request) {

	session := Session(r)
	session.Values["user"] = StringParameter(r, "name")
	session.Values["userId"] = float64(10)
	Forward(r, "index.getuser", 13)
}

func (controller *IndexController) getUsers(w http.ResponseWriter, r *http.Request) string {
	userId := NumberParameter(r, "userId")

	session := Session(r)

	userId += lib.Float64Session(r, "userId")

	session.Values["userId"] = userId

	return fmt.Sprintf("Hello There User %f,%s", userId, session.Values["user"])
}

func (t *IndexController) indexAction(r *http.Request) {
	Forward(r, "index.getuser", 5)
}
