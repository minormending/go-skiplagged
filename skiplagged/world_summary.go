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

	byName := map[string]*CitySummary{}
	for _, trip := range manifest.Trips {
		price := trip.Cost / 100
		if req.Criteria.MaxPrice > 0 && price > req.Criteria.MaxPrice {
			continue
		}
		city := manifest.Cities[trip.City]
		fullName := fmt.Sprintf("%s, %s", city.Name, city.Region)
		if other, ok := byName[fullName]; !ok || other.MinRoundTripPrice > price {
			byName[fullName] = &CitySummary{
				Name:              trip.City,
				FullName:          fullName,
				MinRoundTripPrice: price,
			}
		}
	}

	summaries := []*CitySummary{}
	for _, summary := range byName {
		summaries = append(summaries, summary)
	}
	return summaries, nil
}
