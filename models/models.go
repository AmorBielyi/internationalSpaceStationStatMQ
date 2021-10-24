package models

import "fmt"

type IssResponse struct {
	Visibility string  `json:"visibility"`
	Longitude  float64 `json:"longitude"`
	Latitude   float64 `json:"latitude"`
	Velocity   float64 `json:"velocity"`
}

func (issResponse IssResponse) String() string {
	return fmt.Sprintf("{\"latitude\":%v,\"longitude\":%v,\"velocity\":%v,\"visibility\":\"%s\"}",
		issResponse.Latitude, issResponse.Longitude, issResponse.Velocity, issResponse.Visibility)
}
