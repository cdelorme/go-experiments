package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()

	router.Handle(http.MethodOptions, "/", DefaultOptions)
	router.Handle(http.MethodOptions, "/api*all", CORSOptions)

	http.ListenAndServe(":3000", router)
}

func DefaultOptions(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Allow", "GET,PUT,POST,PATCH,DELETE,OPTIONS")
	w.WriteHeader(http.StatusOK)
}

func CORSOptions(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Allow", "GET,PUT,POST,PATCH,DELETE,OPTIONS")
	w.Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,PATCH,DELETE,OPTIONS")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"))
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Vary", "Origin")
	w.WriteHeader(http.StatusOK)

	// @todo: can parse URI here to load self-documenting API as
	// an optional addendum when responding to an OPTIONS request
}
