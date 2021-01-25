package formatters

import (
	"encoding/json"
	"io"

	"github.com/minormending/go-skiplagged/models"
	"github.com/minormending/go-skiplagged/skiplagged"
)

// ToJSON writes the object out to the stream
func ToJSON(wr io.Writer, req *models.Request, summaries []*skiplagged.CitySummary) error {
	enc := json.NewEncoder(wr)
	enc.SetIndent("", "\t")
	return enc.Encode(struct {
		Request *models.Request           `json:"request"`
		Data    []*skiplagged.CitySummary `json:"data"`
	}{
		Request: req,
		Data:    summaries,
	})
}
