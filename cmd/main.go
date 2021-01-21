package main

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/minormending/go-skiplagged/models"
	"github.com/minormending/go-skiplagged/skiplagged"
)

func main() {
	os.Setenv("HTTP_PROXY", "http://localhost:8888")

	logFile, err := os.OpenFile("output.log", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer logFile.Close()
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)

	req, err := models.NewRequest("NYC", "", "2021-02-18", "2021-02-22", 1)
	if err != nil {
		panic(err)
	}
	req.WithMaxPrice(200).
		WithLeavingCriteria(8, 19).
		WithReturningCriteria(11, 22)

	cities := []*skiplagged.CitySummary{}
	if len(req.TripCity) == 0 {
		cities, err = skiplagged.GetCitySummaryLeavingCity(req)
		if err != nil {
			panic(err)
		}
		for _, city := range cities {
			log.Printf("%s (%s) is $%d\n", city.FullName, city.Name, city.MinRoundTripPrice)
		}
	} else {
		cities = append(cities, &skiplagged.CitySummary{
			Name: req.TripCity,
		})
	}

	for _, city := range cities {
		req.TripCity = city.Name
		summary, err := skiplagged.GetFlightSummaryToCity(req)
		if err != nil {
			log.Println(err)
			continue
		}

		if len(summary.Leaving) == 0 || len(summary.Returning) == 0 {
			log.Printf("did not find viable flights to %s (%s)", summary.FullName, summary.Name)
			continue
		}

		log.Printf("Found flights to %s with min rountrip $%d", summary.FullName, summary.MinRoundTripPrice)
		for _, flight := range summary.Leaving {
			log.Printf("%s => %s (%s) for $%d (%s) leaving @ %s and arrving @ %s\n",
				flight.Departure.Airport,
				summary.FullName,
				flight.Arrival.Airport,
				flight.Price,
				flight.Airline,
				flight.Departure.Time.Format("03:04PM"),
				flight.Arrival.Time.Format("03:04PM"))
		}
		log.Printf("min leaving price is $%d to %s (%s)\n", summary.MinLeavingPrice, summary.FullName, summary.Name)

		for _, flight := range summary.Returning {
			log.Printf("%s (%s) => %s for $%d (%s) leaving @ %s and arriving @ %s\n",
				summary.FullName,
				flight.Departure.Airport,
				flight.Arrival.Airport,
				flight.Price,
				flight.Airline,
				flight.Departure.Time.Format("03:04PM"),
				flight.Arrival.Time.Format("03:04PM"))
		}
		log.Printf("%s (%s) min returning price is $%d\n", summary.FullName, summary.Name, summary.MinReturningPrice)

		time.Sleep(time.Second * 2)
	}

}
