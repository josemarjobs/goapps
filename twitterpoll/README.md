start nsqlookd

```$ nsqlookd```

start nsqd and tell it which instance of nsqlookd to use

```$ nsqd --lookupd-tcp-address=localhost:4160```

start mongod for data services

```$ mongod --dbpath ./db```