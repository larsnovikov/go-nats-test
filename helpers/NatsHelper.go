package Helper

import (
	"../structs"
	"encoding/json"
	"fmt"
	"github.com/nats-io/go-nats-streaming"
	"github.com/nats-io/go-nats-streaming/pb"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"
)

var NatsHandler NatsClient

func (n NatsClient) prepareRequest(route *structs.Route) []byte {
	jsonOut, _ := json.Marshal(route.Request)

	return jsonOut
}

func operateRequest(msg *stan.Msg) {
	request := structs.Request{}

	err := json.Unmarshal([]byte(msg.Data), &request)
	ErrorHandler.Handle(err, func(err error) {}, func(err error) {})

	result := Routing[request.PathPattern].Handler(request)
	fmt.Print(result)
}

func (n NatsClient) Subscribe() {
	sc, err := stan.Connect(n.ClusterID, n.ClientID, stan.NatsURL(n.URL),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Fatalf("Connection lost, reason: %v", reason)
		}))
	if err != nil {
		log.Fatalf("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, n.URL)
	}
	log.Printf("Connected to %s clusterID: [%s] clientID: [%s]\n", n.URL, n.ClusterID, n.ClientID)

	i := 0

	mcb := func(msg *stan.Msg) {
		i++
		operateRequest(msg)
	}

	startOpt := stan.StartAt(pb.StartPosition_NewOnly)

	if n.StartSeq != 0 {
		startOpt = stan.StartAtSequence(n.StartSeq)
	} else if n.DeliverLast {
		startOpt = stan.StartWithLastReceived()
	} else if n.DeliverAll {
		log.Print("subscribing with DeliverAllAvailable")
		startOpt = stan.DeliverAllAvailable()
	} else if n.StartDelta != "" {
		ago, err := time.ParseDuration(n.StartDelta)
		if err != nil {
			sc.Close()
			log.Fatal(err)
		}
		startOpt = stan.StartAtTimeDelta(ago)
	}

	sub, err := sc.QueueSubscribe(n.Topic, n.Qgroup, mcb, startOpt, stan.DurableName(n.Durable))
	if err != nil {
		sc.Close()
		log.Fatal(err)
	}

	log.Printf("Listening on [%s], clientID=[%s], qgroup=[%s] durable=[%s]\n", n.Topic, n.ClientID, n.Qgroup, n.Durable)

	if n.ShowTime {
		log.SetFlags(log.LstdFlags)
	}

	// Wait for a SIGINT (perhaps triggered by user with CTRL-C)
	// Run cleanup when signal is received
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for range signalChan {
			fmt.Printf("\nReceived an interrupt, unsubscribing and closing connection...\n\n")
			// Do not unsubscribe a durable on exit, except if asked to.
			if n.Durable == "" || n.Unsubscribe {
				sub.Unsubscribe()
			}
			sc.Close()
			cleanupDone <- true
		}
	}()
	<-cleanupDone
}

func (n NatsClient) Publish(route structs.Route) {
	var clusterID string
	var clientID string
	var async bool
	var URL string
	var topic string

	topic = route.Topic
	URL = stan.DefaultNatsURL
	clusterID = "test-cluster"
	clientID = "stan-pub"
	async = false

	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(URL))
	if err != nil {
		log.Fatalf("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, URL)
	}
	defer sc.Close()

	msg := []byte(NatsHandler.prepareRequest(&route))

	ch := make(chan bool)
	var glock sync.Mutex
	var guid string
	acb := func(lguid string, err error) {
		glock.Lock()
		log.Printf("Received ACK for guid %s\n", lguid)
		defer glock.Unlock()
		if err != nil {
			log.Fatalf("Error in server ack for guid %s: %v\n", lguid, err)
		}
		if lguid != guid {
			log.Fatalf("Expected a matching guid in ack callback, got %s vs %s\n", lguid, guid)
		}
		ch <- true
	}

	if !async {
		err = sc.Publish(topic, msg)
		if err != nil {
			log.Fatalf("Error during publish: %v\n", err)
		}
		log.Printf("Published [%s] : '%s'\n", topic, msg)
	} else {
		glock.Lock()
		guid, err = sc.PublishAsync(topic, msg, acb)
		if err != nil {
			log.Fatalf("Error during async publish: %v\n", err)
		}
		glock.Unlock()
		if guid == "" {
			log.Fatal("Expected non-empty guid to be returned.")
		}
		log.Printf("Published [%s] : '%s' [guid: %s]\n", topic, msg, guid)

		select {
		case <-ch:
			break
		case <-time.After(5 * time.Second):
			log.Fatal("timeout")
		}
	}
}
