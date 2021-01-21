package clients

// Location is a location
type Location struct {
	City     string   `json:"city"`
	State    string   `json:"state"`
	Airports []string `json:"airports"`
}
