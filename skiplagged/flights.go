package skiplagged

import (
	"fmt"

	"github.com/minormending/go-skiplagged/clients"
	"github.com/minormending/go-skiplagged/models"
)

func flightMeetsLeavingCriteria(flights map[string]clients.Flight, outbound clients.InOutBoundFlight, req *models.Request) (*clients.Flight, error) {
	flight, err := flightMeetsCriteria(flights, outbound, req)
	if err != nil {
		return flight, err
	}

	if !req.Criteria.Leave.After.IsZero() && flight.Segments[0].Departure.Time.Before(req.Criteria.Leave.After) {
		return nil, fmt.Errorf("leaving home too early in the day @ %s", flight.Segments[0].Departure.Time.Format("03:04PM"))
	}
	if !req.Criteria.Leave.Before.IsZero() && flight.Segments[0].Departure.Time.After(req.Criteria.Leave.Before) {
		return nil, fmt.Errorf("leaving home too late in the day @ %s", flight.Segments[0].Departure.Time.Format("03:04PM"))
	}
	for _, airport := range req.Criteria.ExcludeAirports {
		if airport == flight.Segments[0].Departure.Airport {
			return nil, fmt.Errorf("leaving airport has been excluded @ %s", flight.Segments[0].Departure.Airport)
		}
	}
	return flight, nil
}

func flightMeetsReturningCriteria(flights map[string]clients.Flight, inbound clients.InOutBoundFlight, req *models.Request) (*clients.Flight, error) {
	flight, err := flightMeetsCriteria(flights, inbound, req)
	if err != nil {
		return flight, err
	}
	if !req.Criteria.Return.After.IsZero() && flight.Segments[0].Arrival.Time.Before(req.Criteria.Return.After) {
		return nil, fmt.Errorf("arrival home is too early in the day @ %s", flight.Segments[0].Arrival.Time.Format("03:04PM"))
	}
	if !req.Criteria.Return.Before.IsZero() && flight.Segments[0].Arrival.Time.After(req.Criteria.Return.Before) {
		return nil, fmt.Errorf("arrival home is too late in the day @ %s", flight.Segments[0].Arrival.Time.Format("03:04PM"))
	}
	for _, airport := range req.Criteria.ExcludeAirports {
		if airport == flight.Segments[0].Arrival.Airport {
			return nil, fmt.Errorf("arrival airport has been excluded @ %s", flight.Segments[0].Arrival.Airport)
		}
	}
	return flight, nil
}

func flightMeetsCriteria(flights map[string]clients.Flight, bound clients.InOutBoundFlight, req *models.Request) (*clients.Flight, error) {
	if req.Criteria.MaxPrice > 0 {
		flightPrice := bound.OneWayPrice / 100.0
		roundTripPrice := bound.MinRoundTripPrice / 100.0
		if flightPrice == 0 || flightPrice > req.Criteria.MaxPrice {
			return nil, fmt.Errorf("flight price too expensive @ %d", flightPrice)
		} else if roundTripPrice > req.Criteria.MaxPrice {
			return nil, fmt.Errorf("roundtrip price too expensive @ %d", roundTripPrice)
		}
	}
	flight, ok := flights[bound.Flight]
	if !ok {
		return nil, fmt.Errorf("flight not found @ %s", bound.Flight)
	} else if flight.Count > 1 {
		return nil, fmt.Errorf("only non-stop flights allowed @ %d legs", flight.Count)
	} else if len(flight.Segments) > 1 {
		return nil, fmt.Errorf("no hidden city flights allowed @ %d legs", len(flight.Segments))
	}
	return &flight, nil
}
