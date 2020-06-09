package main

import (
	"log"
	"net/http"

	"github.com/bdtomlin/template-example/views"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		views.Render(w, "index", nil, "layouts/main")
	})

	http.HandleFunc("/nested", func(w http.ResponseWriter, r *http.Request) {
		views.Render(w, "index", nil, "layouts/main", "layouts/nested")
	})

	log.Fatal(http.ListenAndServe(":3333", nil))
}
