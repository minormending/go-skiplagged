package models

import "time"

// Flight represents a flight
type Flight struct {
	Price        int           `json:"price"`
	Airline      string        `json:"airline"`
	FlightNumber int           `json:"flight_number"`
	Duration     time.Duration `json:"duration"`
	Departure    struct {
		Time    time.Time `json:"time"`
		Airport string    `json:"airport"`
	} `json:"departure"`
	Arrival struct {
		Time    time.Time `json:"time"`
		Airport string    `json:"airport"`
	} `json:"arrival"`
}
