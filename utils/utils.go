package utils

import (
	"encoding/json"
	"fmt"
	"github.com/amorbielyi/internationalSpaceStationStatMQ/models"
	"io"
	"log"
)

func FailOnErr(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func UnmarshalResponse(response io.Reader, issResponse *models.IssResponse) {
	err := json.NewDecoder(response).Decode(issResponse)
	if err != nil {
		panic(fmt.Sprintf("cannot unmarshal response into struct.\n %v", err))
	}
}
