package example

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/sbiemont/gocsv"
	"github.com/sbiemont/gocsv/lib"
)

func TestSimpleExample(t *testing.T) {
	var content = [][]string{
		{"1", "John Doe", "unused", "2023-06-01", "1.71"},
		{"2", "Jane Doe", "unused", "2023-05-12", ""},
	}

	// Define mapping (column 2 is unused)
	// Unmapped columns will be exported empty
	type row struct {
		ID     int      `csv:"0"`
		Name   string   `csv:"1"`
		Date   lib.Date `csv:"3"`
		Height *float64 `csv:"4,omitempty"`
	}

	heightJohnDoe := 1.71

	johnDoe := row{
		ID:     1,
		Name:   "John Doe",
		Date:   lib.Date(time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC)),
		Height: &heightJohnDoe,
	}
	janeDoe := row{
		ID:     2,
		Name:   "Jane Doe",
		Date:   lib.Date(time.Date(2023, 5, 12, 0, 0, 0, 0, time.UTC)),
		Height: nil,
	}

	Convey("decode", t, func() {
		res, err := gocsv.Decode[row](content)
		So(err, ShouldBeNil)
		So(res, ShouldResemble, []row{
			johnDoe,
			janeDoe,
		})
	})

	Convey("encode", t, func() {
		res, err := gocsv.Encode([]row{
			johnDoe,
			janeDoe,
		})
		So(err, ShouldBeNil)
		So(res, ShouldResemble, [][]string{
			{"1", "John Doe", "", "2023-06-01", "1.710000"},
			{"2", "Jane Doe", "", "2023-05-12", ""},
		})
	})
}
