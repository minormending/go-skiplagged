package formatters

import (
	"encoding/json"
	"io"

	"github.com/minormending/go-skiplagged/skiplagged"
)

// ToJSON writes the object out to the stream
func ToJSON(wr io.Writer, summaries []*skiplagged.CitySummary) error {
	enc := json.NewEncoder(wr)
	enc.SetIndent("", "\t")
	return enc.Encode(summaries)
}
