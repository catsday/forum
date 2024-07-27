package internal

import (
	"net/http"
)

func Router() *http.ServeMux {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", Home)
	mux.HandleFunc("/forum/view", PostView)
	mux.HandleFunc("/forum/create", PostCreate)

	return mux
}
