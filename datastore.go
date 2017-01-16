package main

import (
	"gopkg.in/mgo.v2"
	"log"
)

type DataStore struct {
    session *mgo.Session
}

func (ds * DataStore) init(address string) {
	var err error
	ds.session, err = mgo.Dial(address)
	if err != nil {
		log.Fatal("Error opening database\nerror: %v\n", err)
	}
}

func (ds * DataStore) close() {
	ds.session.Close()
}

func (ds * DataStore) copy() *DataStore {
	return &DataStore{ds.session.Copy()}
}
