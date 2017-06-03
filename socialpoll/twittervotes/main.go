package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	nsq "github.com/bitly/go-nsq"
	"gopkg.in/mgo.v2"
)

var db *mgo.Session

const (
	NSQ_URL      = "localhost:4150"
	MongoDBHosts = "localhost:27017"
	DB           = "ballots"
	COLLECTION   = "polls"
)

func main() {
	var stoplock sync.Mutex
	stop := false
	stopChan := make(chan struct{}, 1)
	signalChan := make(chan os.Signal, 1)
	go func() {
		<-signalChan
		stoplock.Lock()
		stop = true
		stoplock.Unlock()
		log.Println("Stopping...")
		stopChan <- struct{}{}
		closeConn()
	}()
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	if err := dialdb(); err != nil {
		log.Fatal("failed to dial MongoDB:", err)
	}
	defer closedb()

	votes := make(chan string)
	publisherStoppedChan := publichVotes(votes)
	twitterStoppedChan := startTwitterStream(stopChan, votes)
	go func() {
		for {
			time.Sleep(1 * time.Minute)
			closeConn()
			stoplock.Lock()
			if stop {
				stoplock.Unlock()
				return
			}
			stoplock.Unlock()
		}
	}()
	<-twitterStoppedChan
	close(votes)
	<-publisherStoppedChan
}

type poll struct {
	Options []string
}

func publichVotes(votes <-chan string) <-chan struct{} {
	stopchan := make(chan struct{}, 1)
	pub, err := nsq.NewProducer(NSQ_URL, nsq.NewConfig())
	if err != nil {
		log.Panic("failed to create producer:", err)
		return nil
	}
	go func() {
		for vote := range votes {
			pub.Publish("votes", []byte(vote))
		}
		log.Println("Publisher: Stopping")
		pub.Stop()
		log.Println("Publisher: Stopped")
		stopchan <- struct{}{}
	}()
	return stopchan
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
