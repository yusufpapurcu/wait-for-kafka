package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

// This script will try and connect to the topic leader within the timeout, else return an os.Exit(1).
// We connect to the topic leader to ensure a leadership election succeeds and the topic is ready to be consumed/produced.
func main() {
	kafkaEndpoint := os.Getenv("KAFKA_ENDPOINT")
	if kafkaEndpoint == "" {
		log.Fatal("missing KAFKA_ENDPOINT environment variable")
	}

	var topics []string
	for _, tn := range os.Environ() {
		key, value := strings.Split(tn, "=")[0], strings.Split(tn, "=")[1]
		if strings.HasPrefix(key, "KAFKA_TOPIC") {
			topics = append(topics, value)
		}
	}

	fmt.Println(topics)

	timeout := time.After(2 * time.Minute)

	for _, topic := range topics {
		connected := false
		for {
			select {
			case <-timeout:
				log.Fatal("timeout trying to connect to kafka")
			default:
				_, err := kafka.DialLeader(context.Background(), "tcp", kafkaEndpoint, topic, 0)
				if err != nil {
					log.Println("failed to dial leader:", err)
					time.Sleep(10 * time.Second)
				} else {
					log.Printf("connected successfully to %s leader\n", topic)
					connected = true
					break
				}
			}

			if connected {
				break
			}
		}
	}

	fmt.Println("connected successfully to kafka")
}
