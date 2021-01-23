package formatters

import (
	"io"
	"text/template"

	"github.com/minormending/go-skiplagged/skiplagged"
)

const (
	full = `{{range $summary := .}}
# {{$summary.FullName}} (${{$summary.MinRoundTripPrice}})

## Leaving to {{$summary.FullName}}
| Price | Airline | From | To  | Leaving | Arriving |
| ----- | ------- | ---- | --- | ------- | -------- |{{range $flight := $summary.Leaving }}
| ${{$flight.Price}} | {{$flight.Airline}} | {{$flight.Departure.Airport}} | {{$flight.Arrival.Airport}} | {{$flight.Departure.Time.Format "3:04 PM"}} | {{$flight.Arrival.Time.Format "3:04 PM"}} |{{end}}

## Returning from {{$summary.FullName}}
| Price | Airline | From | To  | Leaving | Arriving |
| ----- | ------- | ---- | --- | ------- | -------- |{{range $flight := $summary.Returning }}
| ${{$flight.Price}} | {{$flight.Airline}} | {{$flight.Departure.Airport}} | {{$flight.Arrival.Airport}} | {{$flight.Departure.Time.Format "3:04 PM"}} | {{$flight.Arrival.Time.Format "3:04 PM"}} |{{end}}

{{end}}`
)

// ToMarkdown writes out the list of summaries
func ToMarkdown(wr io.Writer, summaries []*skiplagged.CitySummary) error {
	t, err := template.New("full").Parse(full)
	if err != nil {
		return err
	}
	return t.Execute(wr, summaries)
}
