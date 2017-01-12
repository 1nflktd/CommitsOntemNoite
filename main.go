package main

import (
	"html/template"
	"net/http"
	"strings"
	"time"
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var templates = template.Must(template.ParseFiles("view.html"))

type Commits struct {
	Items []Commit
}

type Commit struct {
	ID bson.ObjectId `bson:"_id,omitempty"`
	Message string `json:"message"`
	AuthorName string `json:"author_name"`
	CommitterDate time.Time `json:"committer_date"`
}

type Page struct {
	Body []Commit
}

type DataStore struct {
    session *mgo.Session
}

func (ds * DataStore) init() {
	var err error
	ds.session, err = mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}
}

func (ds * DataStore) close() {
	ds.session.Close()
}

func (ds * DataStore) copy() *DataStore {
	return &DataStore{ds.session.Copy()}
}

func (ds * DataStore) addCommit(commit Commit) error {
	c := ds.session.DB("commits_ontem").C("commits")
	err := c.Insert(commit)
	return err
}

func getActualDate(format string) string {
	date := time.Now().UTC().AddDate(0, 0, -1) // get day before
	local := date
	location, err := time.LoadLocation("America/Sao_Paulo")
	if err == nil {
		local = local.In(location)
	}
	return date.Format(format)
}

/// Curl call, like
/// curl "https://api.github.com/search/commits?q=shit+committer-date:2017-01-10" \
/// -H 'User-Agent: CommitsNoiteOntem' \
/// -H 'Accept: application/vnd.github.cloak-preview'

func getCommitsAPI() (*Commits, error) {
	searchWords := []string{/*"merda", "coco", "cagada", "droga", "desgraça", "bosta", "pqp", "caralho", */"shit"}
	url := "https://api.github.com/search/commits?q="
	url += strings.Join(searchWords, "+")
	url += "+committer-date:" + getActualDate("2006-01-02") // Y-m-d

	client := &http.Client{Timeout: 30 * time.Second}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Accept", "application/vnd.github.cloak-preview")
	request.Header.Set("User-Agent", "CommitsNoiteOntem")
		
	res, err := client.Do(request)
	if err != nil {
		fmt.Printf("Error na chamada: %v\n", err)
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var data Commits
	err = decoder.Decode(&data)
	if err != nil {
        return nil, err
	}
	
	return &data, nil
}

func loadPage(ds *DataStore) (*Page, error) {
	data, err := getCommitsAPI()
	if err != nil {
		return nil, err
	}
	for _, commit := range data.Items {
		ds.addCommit(commit)
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
		fmt.Printf("Erro loadPage: %v\n", err)
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
	ds := &DataStore{}
	ds.init()
	defer ds.close()

	http.HandleFunc("/", makeHandler(viewHandler, ds))

	http.ListenAndServe(":8080", nil)
}
