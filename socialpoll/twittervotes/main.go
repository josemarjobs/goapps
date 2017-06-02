package main

import "log"
import "gopkg.in/mgo.v2"

var db *mgo.Session

const (
	MongoDBHosts = "localhost:27017"
	DB           = "ballots"
	COLLECTION   = "polls"
)

func main() {
	log.Printf("Credentials: %+v\n", ts)
}

type poll struct {
	Options []string
}

func loadOptions() ([]string, error) {
	var options []string
	iter := db.DB(DB).C(COLLECTION).Find(nil).Iter()
	var p poll
	for iter.Next(&p) {
		options = append(options, p.Options...)
	}
	iter.Close()
	return options, iter.Err()
}

func dialdb() error {
	var err error
	log.Println("dialing mongodb: " + MongoDBHosts)
	db, err = mgo.Dial(MongoDBHosts)
	return err
}

func closedb() {
	db.Close()
	log.Println("closed database connection")
}
