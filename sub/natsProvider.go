package main

import (
	"encoding/json"
	"log"

	"github.com/nats-io/stan.go"
)

const clusterID = "test-cluster"
const clientID = "sub-client"
const subject = "zxc"

func connectToNatsStreaming(orderCh chan<- Order) (stan.Conn, stan.Subscription, error) {
	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL("nats://nats-streaming:4222"))
	if err != nil {
		log.Fatalf("Error connecting to NATS Streaming: %v", err)
	}
	// подписка на Nuts
	sub, err := sc.Subscribe(subject, func(msg *stan.Msg) {
		var order Order
		if err := json.Unmarshal(msg.Data, &order); err != nil {
			log.Println("Error on unmarshal:", err)
			return
		}
		if order.OrderUID != "" {
			log.Println("Recived message:", order.OrderUID)
			orderCh <- order
		}
	}, stan.DeliverAllAvailable())

	return sc, sub, err
}
