package service

import (
	"github.com/gorilla/mux"
	"log"
)

func NewRouter() *mux.Router{
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes{
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
		log.Printf("Route : %s, %s %s", route.Name, route.Method, route.Pattern)
	}

	return router
}
