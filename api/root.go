package api

import (
	"github.com/gorilla/mux"
)

func New() (*mux.Router, error) {

	r := mux.NewRouter()
	hostelrouter := r.PathPrefix("/hostels").Subrouter()

	hostelHandler(hostelrouter)

	return r, nil
}
