package main

import (
	"html/template"
	"net/http"
)

func main() {
	http.HandleFunc("/", indexHandler)
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("monitor"))))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("monitor/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := t.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
