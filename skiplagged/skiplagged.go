package skiplagged

import "github.com/minormending/go-skiplagged/models"

// CitySummary is
type CitySummary struct {
	Name              string
	MinRoundTripPrice int
	MinLeavingPrice   int
	MinReturningPrice int
	Leaving           []*models.Flight
	Returning         []*models.Flight
}
