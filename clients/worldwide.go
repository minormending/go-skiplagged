package clients

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/minormending/go-skiplagged/models"
)

var countryBaseURL = "https://skiplagged.com/api/skipsy.php?from=%s&depart=%s&return=%s&format=v2&counts[adults]=%d&counts[children]=0&_=1611006103100"

// City is a city
type City struct {
	Name     string   `json:"name"`
	Airports []string `json:"airports"`
	Region   string   `json:"region"`
}

// Trip is a trip
type Trip struct {
	City        string `json:"city"`
	Cost        int    `json:"cost"`
	HiddenCity  bool   `json:"hidden_city"`
	RegularCost int    `json:"regular_cost,omitempty"`
}

// CountryResponse is a response
type CountryResponse struct {
	Cities   map[string]City `json:"cities"`
	Airports map[string]struct {
		Name string `json:"name"`
	} `json:"airports"`
	Info struct {
		From Location `json:"from"`
	} `json:"info"`
	Trips    []Trip  `json:"trips"`
	Duration float64 `json:"duration"`
}

// GetWorldwideFlightsFromCity gets the possible cities for a trip
func GetWorldwideFlightsFromCity(req *models.Request) (*CountryResponse, error) {
	url := fmt.Sprintf(countryBaseURL, req.HomeCity, req.LeavingDay.Format("2006-01-02"), req.ReturningDay.Format("2006-01-02"), req.Travelers)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var payload CountryResponse
	if err = json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return nil, err
	}

	return &payload, nil
}
