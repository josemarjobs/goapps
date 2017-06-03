Start nsqlookd

```$ nsqlookd```

Sart nsqd and tell it which instance of nsqlookd to use

```$ nsqd --lookupd-tcp-address=localhost:4160```

Start mongod for data services

```$ mongod --dbpath ./db```

Watch the messages being published to the ```votes``` topic.

```$ nsq_tail --topic="votes" --lookupd-http-address=localhost:4161```