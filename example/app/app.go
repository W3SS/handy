package app

import (
	"fmt"
	. "github.com/go4r/handy"
	"net/http"
)

type IndexController struct{ Annotations }

func (controller *IndexController) Annotates() {

	controller.Annotation(
		Name("index.getuser"),
		Any("/getuser/(number:userId)"),
		controller.getUsers,
	).Annotation(
		Name("index.index"),
		Any("/"),
		controller.indexAction,
	).Annotation(
		Name("index.setuser"),
		Any("/setuser/(:name)"),
		controller.setUser,
	)

}

func (controller *IndexController) setUser(w http.ResponseWriter, r *http.Request) {

	session := Session(r)
	session.Values["user"] = StringParameter(r, "name")
	Forward(r, "index.getuser", 13)
}

func (controller *IndexController) getUsers(w http.ResponseWriter, r *http.Request) string {
	userId := NumberParameter(r, "userId")

	session := Session(r)

	if _, ok := session.Values["userId"]; !ok {
		session.Values["userId"] = 0.0
	}

	userId += session.Values["userId"].(float64)

	session.Values["userId"] = userId

	return fmt.Sprintf("Hello There User %f,%s", userId, session.Values["user"])
}

func (t *IndexController) indexAction(r *http.Request) {
	Forward(r, "index.getuser", 5)
}
