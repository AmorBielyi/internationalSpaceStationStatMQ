package main

import (
	"fmt"
	"github.com/amorbielyi/internationalSpaceStationStatMQ/models"
	"github.com/amorbielyi/internationalSpaceStationStatMQ/utils"
	"github.com/streadway/amqp"
	"log"
	"os"
	"strings"
)

func main() {
	mode := ""
	if len(os.Args) > 1 {
		mode = os.Args[1]
	} else {
		_, _ = fmt.Fprintf(os.Stderr, "Error: argument not provided")
		os.Exit(1)
	}
	if mode == "visibility" {
		fmt.Println("Mode is set to tracking ISS Visibility")
	}
	if mode == "location" {
		fmt.Println("Mode is set to tracking ISS Geolocation and Velocity")
	}

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	utils.FailOnErr(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	utils.FailOnErr(err, "Failed to open a channel")
	defer ch.Close()

	// Declare RMQ Exchange
	err = ch.ExchangeDeclare(
		"iss_geo",
		"fanout",
		true,
		false,
		false,
		false,
		nil)

	utils.FailOnErr(err, "Failed to declare a exchange")

	// Declare RMQ Queue
	q, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil)

	utils.FailOnErr(err, "Failed to declare a queue")

	// QueueBind RMQ
	err = ch.QueueBind(
		q.Name,
		"",
		"iss_geo",
		false,
		nil,
	)

	utils.FailOnErr(err, "Failed to bind a queue")

	// Consume messages
	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil)

	utils.FailOnErr(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {

			issResponse := new(models.IssResponse)
			utils.UnmarshalResponse(strings.NewReader(string(d.Body)), issResponse)

			if mode == "visibility" {
				if issResponse.Visibility == "eclipsed" {
					log.Printf("ISS visibility now is eclipsed")
				} else if issResponse.Visibility == "daylight" {
					log.Printf("ISS visibilily now is daylight")
				}
			} else if mode == "location" {
				log.Printf("ISS geolocation is: longitude: %v and latitude: %v;\n\t\t      "+
					"orbital speed is about %v km/h",
					issResponse.Longitude, issResponse.Latitude, issResponse.Velocity)
			}
		}
	}()

	log.Printf(" Waiting for new ISS data message....")
	<-forever
}
