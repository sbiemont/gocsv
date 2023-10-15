package gocsv

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// See https://data.nantesmetropole.fr/explore/dataset/512042839_composteurs-quartier-nantes-metropole/table
var filename = "example/512042839_composteurs-quartier-nantes-metropole.csv"

type Geolocation struct {
	Lon float64
	Lat float64
}

func (it Geolocation) MarshalCSV() (string, error) {
	return fmt.Sprintf("%.20f, %.20f", it.Lon, it.Lat), nil
}

func (it *Geolocation) UnmarshalCSV(s string) error {
	geo := strings.SplitN(s, ",", 2)

	lon, err := strconv.ParseFloat(strings.Trim(geo[0], " "), 64)
	if err != nil {
		return err
	}
	lat, err := strconv.ParseFloat(strings.Trim(geo[1], " "), 64)
	if err != nil {
		return err
	}
	it.Lon = lon
	it.Lat = lat
	return nil
}

func TestCSV(t *testing.T) {
	type testStruct struct {
		ID          int         `csv:"0"`
		Name        string      `csv:"1"`
		Category    string      `csv:"2"`
		Year        *int        `csv:"3,omitempty"`
		Address     string      `csv:"4"`
		Place       string      `csv:"5"`
		Link        string      `csv:"6"`
		Geolocation Geolocation `csv:"7"`
	}

	Convey("read all", t, func() {
		// Open file
		f, err := os.Open(filename)
		So(err, ShouldBeNil)
		defer f.Close()

		// Read csv
		csvReader := csv.NewReader(f)
		csvReader.Comma = ';'
		data, err := csvReader.ReadAll()
		So(err, ShouldBeNil)

		// Read struct
		res, err := Decode[testStruct](data[1:]) // remove headers (first row)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)

		// Write all
		data2, err := Encode(res)
		So(err, ShouldBeNil)
		So(data2, ShouldNotBeNil)
	})
}
