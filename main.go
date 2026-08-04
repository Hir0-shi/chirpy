package main

import (
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("."))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "index.html")
			return
		}
		fileServer.ServeHTTP(w, r)
	})

	mux.Handle("/assets", http.StripPrefix("/assets/", fileServer))

	_ = http.ListenAndServe(":8080", mux)

}
