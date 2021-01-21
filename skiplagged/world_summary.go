package skiplagged

import (
	"errors"
	"fmt"

	"github.com/minormending/go-skiplagged/clients"
	"github.com/minormending/go-skiplagged/models"
)

// GetCitySummaryLeavingCity provides possible options from city
func GetCitySummaryLeavingCity(req *models.Request) ([]*CitySummary, error) {
	manifest, err := clients.GetWorldwideFlightsFromCity(req)
	if err != nil {
		return nil, errors.New("unable to get flights from city")
	}

	summary := []*CitySummary{}
	for _, trip := range manifest.Trips {
		price := trip.Cost / 100
		if req.Criteria.MaxPrice > 0 && price > req.Criteria.MaxPrice {
			continue
		}
		city := manifest.Cities[trip.City]
		summary = append(summary, &CitySummary{
			Name:              trip.City,
			FullName:          fmt.Sprintf("%s, %s", city.Name, city.Region),
			MinRoundTripPrice: price,
		})
	}

	return summary, nil
}
