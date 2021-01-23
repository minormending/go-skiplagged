package skiplagged

import (
	"errors"
	"fmt"
	"time"

	"github.com/minormending/go-skiplagged/clients"
	"github.com/minormending/go-skiplagged/models"
)

// GetFlightSummaryToCity filters a city
func GetFlightSummaryToCity(req *models.Request) (*CitySummary, error) {
	manifest, err := clients.GetFlightsToCity(req)
	if err != nil {
		return nil, errors.New("unable to get flights to city")
	}

	city := manifest.Info.To
	summary := CitySummary{
		Name:              req.TripCity,
		FullName:          fmt.Sprintf("%s, %s", city.City, city.State),
		MinRoundTripPrice: 0,
		MinLeavingPrice:   0,
		MinReturningPrice: 0,
		Leaving:           []*models.Flight{},
		Returning:         []*models.Flight{},
	}

	for _, outbound := range manifest.Itineraries.Outbound {
		flight, err := flightMeetsLeavingCriteria(manifest.Flights, outbound, req)
		if err != nil {
			continue
		}

		price := outbound.OneWayPrice / 100.0
		if summary.MinLeavingPrice == 0 || price < summary.MinLeavingPrice {
			summary.MinLeavingPrice = price
		}

		leg := flight.Segments[0]
		summary.Leaving = append(summary.Leaving, &models.Flight{
			Price:        price,
			Airline:      manifest.Airlines[leg.Airline].Name,
			FlightNumber: leg.FlightNumber,
			Duration:     time.Duration(leg.Duration),
			Departure:    leg.Departure,
			Arrival:      leg.Arrival,
		})
	}

	for _, inbound := range manifest.Itineraries.Inbound {
		flight, err := flightMeetsReturningCriteria(manifest.Flights, inbound, req)
		if err != nil {
			continue
		}

		price := inbound.OneWayPrice / 100.0
		if summary.MinReturningPrice == 0 || price < summary.MinReturningPrice {
			summary.MinReturningPrice = price
		}

		leg := flight.Segments[0]
		summary.Returning = append(summary.Returning, &models.Flight{
			Price:        price,
			Airline:      manifest.Airlines[leg.Airline].Name,
			FlightNumber: leg.FlightNumber,
			Duration:     time.Duration(leg.Duration),
			Departure:    leg.Departure,
			Arrival:      leg.Arrival,
		})
	}

	if len(summary.Leaving) > 0 && len(summary.Returning) > 0 {
		summary.MinRoundTripPrice = summary.MinLeavingPrice + summary.MinReturningPrice

		if req.Criteria.MaxPrice > 0 && summary.MinRoundTripPrice > req.Criteria.MaxPrice {
			summary.Leaving = []*models.Flight{}
			summary.Returning = []*models.Flight{}
		}
	}
	return &summary, nil
}
