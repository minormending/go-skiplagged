package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/minormending/go-skiplagged/models"
)

var cityBaseURL = "https://skiplagged.com/api/search.php?from=%s&to=%s&depart=%s&return=%s&poll=true&format=v3&counts[adults]=%d&counts[children]=0&_=1611006103106"

// Flight is a flight
type Flight struct {
	Segments []struct {
		Airline      string `json:"airline"`
		FlightNumber int    `json:"flight_number"`
		Departure    struct {
			Time    time.Time `json:"time"`
			Airport string    `json:"airport"`
		} `json:"departure"`
		Arrival struct {
			Time    time.Time `json:"time"`
			Airport string    `json:"airport"`
		} `json:"arrival"`
		Duration int `json:"duration"`
	} `json:"segments"`
	Duration int    `json:"duration"`
	Count    int    `json:"count"`
	Data     string `json:"data"`
}

// InOutBoundFlight is a
type InOutBoundFlight struct {
	Data              string `json:"data"`
	Flight            string `json:"flight"`
	MinRoundTripPrice int    `json:"min_round_trip_price,omitempty"`
	OneWayPrice       int    `json:"one_way_price,omitempty"`
}

// CityResponse represents flights for a city
type CityResponse struct {
	Airlines map[string]struct {
		Name string `json:"name"`
	} `json:"airlines"`
	Cities map[string]struct {
		Name string `json:"name"`
	} `json:"cities"`
	Airports map[string]struct {
		Name string `json:"name"`
	} `json:"airports"`
	Flights     map[string]Flight `json:"flights"`
	Itineraries struct {
		Outbound []InOutBoundFlight `json:"outbound"`
		Inbound  []InOutBoundFlight `json:"inbound"`
	} `json:"itineraries"`
	Info struct {
		From Location `json:"from"`
		To   Location `json:"to"`
	} `json:"info"`
	Duration float64 `json:"duration"`
}

// GetFlightsToCity get trips for a city
func GetFlightsToCity(req *models.Request) (*CityResponse, error) {
	url := fmt.Sprintf(cityBaseURL, req.HomeCity, req.TripCity, req.LeavingDay.Format("2006-01-02"), req.ReturningDay.Format("2006-01-02"), req.Travelers)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var payload CityResponse
	if err = json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return nil, err
	}

	return &payload, nil
}
