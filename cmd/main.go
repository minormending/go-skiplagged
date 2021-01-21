package main

import (
	"fmt"
	"os"

	"github.com/minormending/go-skiplagged/models"
	"github.com/minormending/go-skiplagged/skiplagged"
)

func main() {
	os.Setenv("HTTP_PROXY", "http://localhost:8888")

	req, err := models.NewRequest("NYC", "AUS", "2021-02-18", "2021-02-22", 1)
	if err != nil {
		panic(err)
	}
	req.WithMaxPrice(200).
		WithLeavingCriteria(8, 19).
		WithReturningCriteria(11, 22)

	cities := []skiplagged.CitySummary{}
	if len(req.TripCity) == 0 {
		summary, err := skiplagged.GetCitySummaryLeavingCity(req)
		if err != nil {
			panic(err)
		}
		for _, city := range summary {
			fmt.Printf("%s is $%d\n", city.Name, city.MinRoundTripPrice)
		}
	} else {
		cities = append(cities, skiplagged.CitySummary{
			Name: req.TripCity,
		})
	}

	for _, city := range cities {
		req.TripCity = city.Name
		summary, err := skiplagged.GetFlightSummaryToCity(req)
		if err != nil {
			fmt.Println(err)
			continue
		}

		for _, flight := range summary.Leaving {
			fmt.Printf("%s => %s for $%d (%s) leaving @ %s and arrving @ %s\n",
				flight.Departure.Airport,
				flight.Arrival.Airport,
				flight.Price,
				flight.Airline,
				flight.Departure.Time.Format("03:04PM"),
				flight.Arrival.Time.Format("03:04PM"))
		}
		fmt.Printf("min leaving price is $%d\n", summary.MinLeavingPrice)

		for _, flight := range summary.Returning {
			fmt.Printf("%s => %s for $%d (%s) leaving @ %s and arriving @ %s\n",
				flight.Departure.Airport,
				flight.Arrival.Airport,
				flight.Price,
				flight.Airline,
				flight.Departure.Time.Format("03:04PM"),
				flight.Arrival.Time.Format("03:04PM"))
		}
		fmt.Printf("min returning price is $%d\n", summary.MinReturningPrice)

		fmt.Printf("min roundtrip price is $%d\n", summary.MinRoundTripPrice)
		return
	}

}
