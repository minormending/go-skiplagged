package models

import (
	"errors"
	"time"
)

// Request represents a request
type Request struct {
	HomeCity     string    `json:"home_city"`
	TripCity     string    `json:"trip_city,omitempty"`
	LeavingDay   time.Time `json:"leaving_day"`
	ReturningDay time.Time `json:"returning_day"`
	Travelers    int       `json:"travelers"`
	Criteria     criteria  `json:"criteria"`
}

type betweenTime struct {
	After  time.Time `json:"after"`
	Before time.Time `json:"before"`
}

type criteria struct {
	MaxPrice        int          `json:"max_price"`
	ExcludeAirports []string     `json:"exclude_airports"`
	Leave           *betweenTime `json:"leave_between"`
	Return          *betweenTime `json:"return_between"`
}

// NewRequest creates a skiplagged API request based on the specified parameters.
func NewRequest(fromCity, toCity, fromDay, toDay string, travelers int) (*Request, error) {
	fromTime, err := time.ParseInLocation("2006-01-02", fromDay, time.Local)
	if err != nil {
		return nil, errors.New("starting date format invalid, should be yyyy-MM-dd")
	}
	toTime, err := time.ParseInLocation("2006-01-02", toDay, time.Local)
	if err != nil {
		return nil, errors.New("ending date format invalid, should be yyyy-MM-dd")
	}
	return &Request{
		HomeCity:     fromCity,
		TripCity:     toCity,
		LeavingDay:   fromTime,
		ReturningDay: toTime,
		Travelers:    travelers,
		Criteria:     criteria{},
	}, nil
}

// WithMaxPrice sets the maximum flight price
func (r *Request) WithMaxPrice(price int) *Request {
	r.Criteria.MaxPrice = price
	return r
}

// WithLeavingCriteria sets the departure flight criteria from the home city
func (r *Request) WithLeavingCriteria(afterHour, beforeHour int) *Request {
	var after, before time.Time
	if afterHour != 0 {
		after = r.LeavingDay.Add(time.Hour * time.Duration(afterHour))
	}
	if beforeHour != 0 {
		before = r.LeavingDay.Add(time.Hour * time.Duration(beforeHour))
	}
	r.Criteria.Leave = &betweenTime{
		After:  after,
		Before: before,
	}
	return r
}

// WithReturningCriteria sets the departure flight criteria back to the home city
func (r *Request) WithReturningCriteria(afterHour, beforeHour int) *Request {
	var after, before time.Time
	if afterHour != 0 {
		after = r.ReturningDay.Add(time.Hour * time.Duration(afterHour))
	}
	if beforeHour != 0 {
		before = r.ReturningDay.Add(time.Hour * time.Duration(beforeHour))
	}
	r.Criteria.Return = &betweenTime{
		After:  after,
		Before: before,
	}
	return r
}

// WithExcludeAirportsCriteria sets the excluded airports
func (r *Request) WithExcludeAirportsCriteria(airports []string) *Request {
	for _, airport := range airports {
		if len(airport) > 0 {
			r.Criteria.ExcludeAirports = append(r.Criteria.ExcludeAirports, airport)
		}
	}
	return r
}
