package main

import (
	"html/template"
	"net/http"
	"log"
	"os"
)

var templates = template.Must(template.ParseFiles("templates/view.tmpl", "templates/header.tmpl", "templates/footer.tmpl"))

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
	templates.ExecuteTemplate(w, "header.tmpl", p)
	err := templates.ExecuteTemplate(w, templ + ".tmpl", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	templates.ExecuteTemplate(w, "footer.tmpl", p)
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

    address := os.Getenv("DATABASE_URL")

    if address == "" {
    	log.Fatal("$DATABASE_URL must be set")
    }

	ds := &DataStore{}
	ds.init(address)
	defer ds.close()

	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	http.HandleFunc("/", makeHandler(viewHandler, ds))

	http.ListenAndServe(":" + port, nil)
}
