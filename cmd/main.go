package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

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
	help            = flag.Bool("help", false, "print help infomation")
)

var (
	infoLogger *log.Logger
)

func init() {
	infoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
}

func saveJSON(filename string, summaries []*skiplagged.CitySummary) error {
	if len(filename) > 0 {
		jsonfile, err := os.OpenFile("summary.json", os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			return err
		}
		defer jsonfile.Close()

		err = formatters.ToJSON(jsonfile, summaries)
		if err != nil {
			return err
		}
	}
	return nil
}

func saveMarkdown(filename string, summaries []*skiplagged.CitySummary) error {
	if len(filename) > 0 {
		markdown, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			return err
		}
		defer markdown.Close()

		err = formatters.ToMarkdown(markdown, summaries)
		if err != nil {
			return err
		}
	}
	return nil
}

func logCitySummaries(summaries []*skiplagged.CitySummary) {
	for _, summary := range summaries {
		infoLogger.Printf("%s (%s) is $%d\n", summary.FullName, summary.Name, summary.MinRoundTripPrice)
		for _, flight := range summary.Leaving {
			infoLogger.Printf("%s => %s (%s) for $%d (%s) leaving @ %s and arrving @ %s\n",
				flight.Departure.Airport,
				summary.FullName,
				flight.Arrival.Airport,
				flight.Price,
				flight.Airline,
				flight.Departure.Time.Format("03:04PM"),
				flight.Arrival.Time.Format("03:04PM"),
			)
		}
		infoLogger.Printf("min leaving price is $%d to %s (%s)\n", summary.MinLeavingPrice, summary.FullName, summary.Name)

		for _, flight := range summary.Returning {
			infoLogger.Printf("%s (%s) => %s for $%d (%s) leaving @ %s and arriving @ %s\n",
				summary.FullName,
				flight.Departure.Airport,
				flight.Arrival.Airport,
				flight.Price,
				flight.Airline,
				flight.Departure.Time.Format("03:04PM"),
				flight.Arrival.Time.Format("03:04PM"),
			)
		}
		infoLogger.Printf("%s (%s) min returning price is $%d\n", summary.FullName, summary.Name, summary.MinReturningPrice)
	}
}

func usage() {
	fmt.Println(`Usage: skiplagged [OPTIONS] ORIGIN START_DATE END_DATE

Gets flight information from the Skiplagged API.

Arguments:
	ORIGIN		departure city or airport
	START_DATE	departure date, yyyy-MM-dd
	END_DATE	return date, yyyy-MM-dd
	
Options:`)
	flag.PrintDefaults()
}

func main() {
	flag.Parse()
	if *help {
		usage()
		return
	}

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
	if len(*outputJSON) > 0 || len(*outputMD) > 0 {
		infoLogger.SetOutput(ioutil.Discard)
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
	}

	if *skipWorldwide == false {
		summaries = skiplagged.GetAllFlightSummariesToCity(req, summaries)
	}
	logCitySummaries(summaries)

	err = saveJSON(*outputJSON, summaries)
	if err != nil {
		panic(err)
	}

	err = saveMarkdown(*outputMD, summaries)
	if err != nil {
		panic(err)
	}
}
