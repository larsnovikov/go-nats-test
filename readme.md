Subscriber
`go run nats.go --last -c test-cluster -qgroup testgroup -id myID1 post_detail`

`go run nats.go --last -c test-cluster -qgroup testgroup -id myID2 post_detail`

`go run nats.go --last -c test-cluster -qgroup testgroup -id myID3 post_detail`

SERVER
`cd $GOPATH/src/github.com/nats-io/nats-streaming-server`

`go run nats-streaming-server.go`