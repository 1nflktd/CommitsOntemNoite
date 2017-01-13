package main

import (
	"time"
	"gopkg.in/mgo.v2/bson"
)

type Commits struct {
	Items []Commit
}

type Commit struct {
	ID bson.ObjectId `bson:"_id,omitempty"`
	Message string `bson:"message" json:"message"`
	AuthorName string `bson:"author_name" json:"author_name"`
	CommitterDate time.Time `bson:"committer_date" json:"committer_date"`
}
