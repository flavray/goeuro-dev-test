package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
)

// Record represents an object coming from the webservice
type Record struct {
	ID       int    `json:"_id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Position struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"geo_position"`
}

// CSV returns a CSV representation of the record
func (r Record) CSV() string {
	return fmt.Sprintf("%d,%s,%s,%f,%f", r.ID, r.Name, r.Type, r.Position.Latitude, r.Position.Longitude)
}

// API base URL
var BaseURL string

// FetchCity fetches city data from the webservice
func FetchCity(city string) []Record {
	response, err := http.Get(BaseURL + city)

	if err != nil {
		Exit(err)
	}

	var records []Record

	decoder := json.NewDecoder(response.Body)

	if err = decoder.Decode(&records); err != nil {
		Exit(err)
	}

	return records
}

// WriteRecords writes record as CSV
func WriteRecords(records []Record, out *os.File) {
	writer := bufio.NewWriter(out)

	fmt.Fprintln(writer, "_id,name,type,latitude,longitude")

	for _, record := range records {
		fmt.Fprintln(writer, record.CSV())
	}

	writer.Flush()
}

// Exit exits the program, displaying the error that made the program end early
func Exit(err error) {
	fmt.Println(err.Error())
	os.Exit(1)
}

func main() {
	BaseURL = "http://api.goeuro.com/api/v2/position/suggest/en/"

	if len(os.Args) != 2 {
		Exit(errors.New("Usage: ./goeuro-dev-test <city>"))
	} else {
		WriteRecords(FetchCity(os.Args[1]), os.Stdout)
	}
}
