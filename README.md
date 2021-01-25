# go-skiplagged
 
Get flight information from the Skiplagged API using custom price and depature/arrival time filters. Provides possible flight options given on a starting city and trip dates and can output the results in JSON or a formatted markdown file.

# Install

```
go get -u github.com/minormending/go-skiplagged
```

# Usage

```
Usage: skiplagged [OPTIONS] ORIGIN START_DATE END_DATE

Gets flight information from the Skiplagged API.

Arguments:
        ORIGIN          departure city or airport
        START_DATE      departure date, yyyy-MM-dd
        END_DATE        return date, yyyy-MM-dd

Options:
  -exclude string
        exclude airports from the trip
  -help
        print help infomation
  -leaveAfter int
        initial departure flight must be after this hour
  -leaveBefore int
        initial departure flight must be before this hour
  -maxPrice int
        maximum price for the entire trip
  -outjson string
        save trip results as json with the specified filename.
  -outmd string
        save trip results as markdown with the specified filename.
  -overwrite
        overwrite existing output file.
  -proxy
        sets the http proxy for requests
  -returnAfter int
        destination return flight must be after this hour
  -returnBefore int
        destination return flight must be before this hour
  -skipworldwide
        skip compute trips for all cities
  -to string
        destination city or airport, optional
  -travelers int
        amount of travelers for the trip (default 1)
```

# Example

```
skiplagged --to=AUS --maxPrice=200 --leaveBefore=13 --returnAfter=12 --returnBefore=19 NYC 2021-03-04 2021-03-08

Austin, Texas (AUS) is $96
EWR => Austin, Texas (AUS) for $48 (United Airlines) leaving @ 07:30AM and arrving @ 10:18AM
EWR => Austin, Texas (AUS) for $48 (United Airlines) leaving @ 11:00AM and arrving @ 01:52PM
EWR => Austin, Texas (AUS) for $49 (JetBlue Airways) leaving @ 06:45AM and arrving @ 09:46AM
JFK => Austin, Texas (AUS) for $114 (JetBlue Airways) leaving @ 09:41AM and arrving @ 12:41PM
min leaving price is $48 to Austin, Texas (AUS)
Austin, Texas (AUS) => EWR for $48 (United Airlines) leaving @ 07:45AM and arriving @ 12:29PM
Austin, Texas (AUS) => EWR for $48 (United Airlines) leaving @ 11:05AM and arriving @ 04:05PM
Austin, Texas (AUS) => JFK for $114 (JetBlue Airways) leaving @ 01:23PM and arriving @ 05:50PM
Austin, Texas (AUS) min returning price is $48
```

This find all flights costing less than $200 from any airport in New York City (NYC) to the airport AUS in Austin, TX departing NYC on March 4th before 1pm and returning to NYC between 12pm and 7pm on March 8th.

The results are printed to the console, but we can easily generate a markdown file with the `--outmd=summary.md` option.
![image](https://user-images.githubusercontent.com/54426030/105773491-f8f87200-5f31-11eb-9443-1f6eb8633ef1.png)
