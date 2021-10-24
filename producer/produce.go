package main

import (
	"fmt"
	"github.com/amorbielyi/internationalSpaceStationStatMQ/models"
	"github.com/amorbielyi/internationalSpaceStationStatMQ/utils"

	"github.com/streadway/amqp"
	"log"
	"net/http"
	"strconv"
	"time"
)

func obtainISSRealTimeData() *models.IssResponse {

	sateliteNORADId := "25544"
	path := "https://api.wheretheiss.at/v1/satellites/"

	req, err := http.NewRequest("GET", path+sateliteNORADId, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	issResponse := new(models.IssResponse)

	if xRateLimitRemaining, ok := resp.Header["X-Rate-Limit-Remaining"]; ok {
		if xRateLimitInterval, ok := resp.Header["X-Rate-Limit-Interval"]; ok {
			limit, err := strconv.Atoi(xRateLimitRemaining[0])
			if err != nil {
				panic(err)
			}
			if limit == 0 {
				fmt.Printf("Failed to obtain ISS data: your limit is 0, waiting %s and continue... ",
					xRateLimitInterval[0])
				time.Sleep(6 * time.Minute)
			}
		}
	}

	utils.UnmarshalResponse(resp.Body, issResponse)
	return issResponse
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	utils.FailOnErr(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	utils.FailOnErr(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"iss_geo",
		"fanout",
		true,
		false,
		false,
		false,
		nil)

	utils.FailOnErr(err, "Failed to declare an exchange")

	for {
		issgeo := obtainISSRealTimeData()
		log.Println("New ISS data received for publication: ", issgeo)
		body := issgeo.String()

		err = ch.Publish(
			"iss_geo",
			"",
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})
		utils.FailOnErr(err, "Failed to publish a message")
		time.Sleep(2 * time.Second)

	}
}
