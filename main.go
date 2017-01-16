package main

import (
	"html/template"
	"net/http"
	"log"
	"os"
)

var templates = template.Must(template.ParseFiles("templates/view.html"))

type Page struct {
	Body []Commit
}

func loadPage(ds *DataStore) (*Page, error) {
	data, err := getCommits(ds)
	if err != nil {
		return nil, err
	}
	return &Page{Body: data.Items}, nil
}

func renderTemplate(w http.ResponseWriter, templ string, p *Page) {
	err := templates.ExecuteTemplate(w, templ + ".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request, ds *DataStore) {
	p, err := loadPage(ds)
	if err != nil {
		log.Printf("error loadPage\nerror: %v\n", err)
		return
	}
	renderTemplate(w, "view", p)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, *DataStore), dsMaster *DataStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ds := dsMaster.copy()
		defer ds.close()
		fn(w, r, ds)
	}
}

func main() {
	port := os.Getenv("PORT")

    if port == "" {
        log.Fatal("$PORT must be set")
    }

	ds := &DataStore{}
	ds.init()
	defer ds.close()

	http.HandleFunc("/", makeHandler(viewHandler, ds))

	http.ListenAndServe(":" + port, nil)
}
