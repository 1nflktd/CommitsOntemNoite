package main

import (
	"gopkg.in/mgo.v2"
)

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
