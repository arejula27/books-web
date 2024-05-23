package web

import (
	"books/web/pages"
	"net/http"
)

func HelloWebHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}

	name := r.FormValue("name")
	component := pages.HelloPost(name)
	component.Render(r.Context(), w)
}
