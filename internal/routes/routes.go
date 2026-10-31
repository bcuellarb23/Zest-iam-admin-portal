package routes

import (
	"net/http"

	"github.com/bcuellarb23/Zest-iam-admin-portal/internal/auth"
	"github.com/julienschmidt/httprouter"
)

// Function to handle routes
func SetupRouter() *httprouter.Router {
	// Router will return the the requested http handler
	router := httprouter.New()
	router.GET("/", DashboardHandler)
	router.GET("/signin", auth.SigninHandler)
	router.GET("/signout", auth.SignoutHandler)

	return router
}

func DashboardHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//
}
