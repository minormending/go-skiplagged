package main

import (
	"errors"
	"flag"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/minormending/go-skiplagged/formatters"
	"github.com/minormending/go-skiplagged/models"
	"github.com/minormending/go-skiplagged/skiplagged"
)

var (
	proxy           = flag.String("proxy", "", "sets the http proxy for requests")
	toCity          = flag.String("to", "", "destination city or airport, optional")
	skipWorldwide   = flag.Bool("skipworldwide", false, "skip compute trips for all cities")
	travelers       = flag.Int("travelers", 1, "amount of travelers for the trip")
	maxPrice        = flag.Int("maxPrice", 0, "maximum price for the entire trip")
	leaveAfter      = flag.Int("leaveAfter", 0, "initial departure flight must be after this hour")
	leaveBefore     = flag.Int("leaveBefore", 0, "initial departure flight must be before this hour")
	returnAfter     = flag.Int("returnAfter", 0, "destination return flight must be after this hour")
	returnBefore    = flag.Int("returnBefore", 0, "destination return flight must be before this hour")
	excludeAirports = flag.String("exclude", "", "exclude airports from the trip")
	outputMD        = flag.String("outmd", "", "save trip results as markdown with the specified filename.")
	outputJSON      = flag.String("outjson", "", "save trip results as json with the specified filename.")
)

func saveJSON(filename string, summaries []*skiplagged.CitySummary) error {
	jsonfile, err := os.OpenFile("summary.json", os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer jsonfile.Close()

	err = formatters.ToJSON(jsonfile, summaries)
	if err != nil {
		return err
	}
	return nil
}

func saveMarkdown(filename string, summaries []*skiplagged.CitySummary) error {
	markdown, err := os.OpenFile(filename, os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer markdown.Close()

	err = formatters.ToMarkdown(markdown, summaries)
	if err != nil {
		return err
	}
	return nil
}

func analyzeCities(req *models.Request, cities []*skiplagged.CitySummary) []*skiplagged.CitySummary {
	summaries := []*skiplagged.CitySummary{}
	for _, city := range cities {
		req.TripCity = city.Name
		summary, err := skiplagged.GetFlightSummaryToCity(req)
		if err != nil {
			log.Println(err)
			continue
		}

		if len(summary.Leaving) == 0 || len(summary.Returning) == 0 {
			log.Printf("did not find viable flights to %s (%s)", summary.FullName, summary.Name)
			continue
		}

		summaries = append(summaries, summary)
		time.Sleep(time.Second * 2)
	}
	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].MinRoundTripPrice < summaries[j].MinRoundTripPrice
	})
	return summaries
}

func main() {
	flag.Parse()
	fromCity := flag.Arg(0)
	start := flag.Arg(1)
	end := flag.Arg(2)
	if len(fromCity) == 0 {
		panic(errors.New("must specify a origin city or airport"))
	} else if len(start) == 0 || len(end) == 0 {
		panic(errors.New("must specify a start and end date"))
	}

	if len(*proxy) > 0 {
		os.Setenv("HTTP_PROXY", *proxy)
	}

	req, err := models.NewRequest(fromCity, *toCity, start, end, *travelers)
	if err != nil {
		panic(err)
	}
	req.WithMaxPrice(*maxPrice).
		WithLeavingCriteria(*leaveAfter, *leaveBefore).
		WithReturningCriteria(*returnAfter, *returnBefore).
		WithExcludeAirportsCriteria(strings.Split(*excludeAirports, ","))

	summaries := []*skiplagged.CitySummary{}
	if len(*toCity) > 0 {
		summaries = append(summaries, &skiplagged.CitySummary{Name: *toCity})
	} else {
		summaries, err = skiplagged.GetCitySummaryLeavingCity(req)
		if err != nil {
			panic(err)
		}
		for _, city := range summaries {
			log.Printf("%s (%s) is $%d\n", city.FullName, city.Name, city.MinRoundTripPrice)
		}
	}

	if *skipWorldwide == false {
		summaries = analyzeCities(req, summaries)
	}

	if len(*outputJSON) > 0 {
		saveJSON(*outputJSON, summaries)
		if err != nil {
			panic(err)
		}
	}

	if len(*outputMD) > 0 {
		saveMarkdown(*outputMD, summaries)
		if err != nil {
			panic(err)
		}
	}
}
