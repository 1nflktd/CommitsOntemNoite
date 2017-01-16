package main

import (
	"strings"
	"encoding/json"
	"time"
	"net/http"
	"log"
	"gopkg.in/mgo.v2/bson"
)

func (ds * DataStore) addCommit(commit Commit) error {
	c := ds.session.DB("commits_ontem").C("commits")
	err := c.Insert(commit)
	return err
}

func (ds * DataStore) getCommits() (*Commits, error) {
	c := ds.session.DB("commits_ontem").C("commits")
	commits := Commits{}
	err := c.Find(bson.M{
		"committer_date" : bson.M{
			"$gte" : getOnlyDate(0, 0, -1),
			"$lte" : getOnlyDate(0, 0, 0),
		},
	}).Limit(100).All(&commits.Items)
	return &commits, err
}

func (ds * DataStore) truncateCommits() error {
	_, err := ds.session.DB("commits_ontem").C("commits").RemoveAll(bson.M{})
	return err
}

/// Curl call, like
/// curl "https://api.github.com/search/commits?q=shit+committer-date:2017-01-10" \
/// -H 'User-Agent: CommitsNoiteOntem' \
/// -H 'Accept: application/vnd.github.cloak-preview'

func getCommitsAPI() (*Commits, error) {
	searchWords := []string{/*"merda", "coco", "cagada", "droga", "desgraça", "bosta", "pqp", "caralho", */"shit"}
	url := "https://api.github.com/search/commits?q="
	url += strings.Join(searchWords, "+")
	url += "+committer-date:" + getOnlyDateFormat("2006-01-02", 0, 0, -1) // Y-m-d

	log.Printf("call getCommitsAPI\nurl: %v\n", url)

	client := &http.Client{Timeout: 30 * time.Second}

	request, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Printf("error new request\nerror: %v\n", err)
		return nil, err
	}

	request.Header.Set("Accept", "application/vnd.github.cloak-preview")
	request.Header.Set("User-Agent", "CommitsNoiteOntem")
		
	res, err := client.Do(request)
	if err != nil {
		log.Printf("error receiving data\nerror: %v\n", err)
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var data Commits
	err = decoder.Decode(&data)

	if err != nil {
		log.Printf("error decoding data\nerror: %v\n", err)
        return nil, err
	}
	
	return &data, nil
}

func getCommitsDB(ds *DataStore) (*Commits, error) {
	commits, err := ds.getCommits()

	if err != nil {
		log.Printf("error call getCommitsDB\nerror: %v\n", err)
		return nil, err
	}

	if len(commits.Items) == 0 { // database nao contem dados de ontem, entao pesquisar
		ds.truncateCommits()
		log.Printf("truncate commits\n")
		return nil, nil
	}

	return commits, nil
}

func saveCommitsDB(ds *DataStore, commits *Commits) (err error) {
	for _, commit := range commits.Items {
		err = ds.addCommit(commit)
		if err != nil {
			break
		}
	}
	return err
}

func getCommits(ds *DataStore) (*Commits, error) {
	data, err := getCommitsDB(ds)
	if err != nil {
		return nil, err
	}

	if data == nil {
		data, err = getCommitsAPI()
		if err != nil {
			return nil, err
		}
		saveCommitsDB(ds, data)
	}

	return data, nil
}
