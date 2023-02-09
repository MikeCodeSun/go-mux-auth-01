package route

import (
	"net/http"

	"github.com/MikeCodeSun/go-mux-auth/controller"
	"github.com/MikeCodeSun/go-mux-auth/util"
	"github.com/gorilla/mux"
)

func UserRoute() *mux.Router{
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", controller.HomePage).Methods("GET")
	r.HandleFunc("/user/register", controller.Register).Methods("POST")
	r.HandleFunc("/user/login", controller.Login).Methods("POST")
	r.HandleFunc("/user/logout", controller.Logout).Methods("GET")
	apiR := r.PathPrefix("/api").Subrouter()
	apiR.HandleFunc("/protected", controller.ProtectedPage).Methods(http.MethodGet)
	apiR.Use(util.AuthJwt)
  return r
}