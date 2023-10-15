package lib

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBool(t *testing.T) {
	Convey("marshal", t, func() {
		Convey("when true", func() {
			res, err := Bool(true).MarshalCSV()
			So(err, ShouldBeNil)
			So(res, ShouldEqual, "true")
		})

		Convey("when false", func() {
			res, err := Bool(false).MarshalCSV()
			So(err, ShouldBeNil)
			So(res, ShouldEqual, "false")
		})
	})

	Convey("unmarshal", t, func() {
		Convey("when true", func() {
			var res Bool
			err := res.UnmarshalCSV("true")
			So(err, ShouldBeNil)
			So(res, ShouldEqual, Bool(true))
		})

		Convey("when false", func() {
			var res Bool
			err := res.UnmarshalCSV("false")
			So(err, ShouldBeNil)
			So(res, ShouldEqual, Bool(false))
		})
	})
}

func TestDuration(t *testing.T) {
	Convey("marshal", t, func() {
		res, err := Duration(time.Hour + 2*time.Minute + 3*time.Second).MarshalCSV()
		So(err, ShouldBeNil)
		So(res, ShouldEqual, "1h2m3s")
	})

	Convey("unmarshal", t, func() {
		Convey("when ok", func() {
			var res Duration
			err := res.UnmarshalCSV("1h2m3s")
			So(err, ShouldBeNil)
			So(res, ShouldEqual, Duration(time.Hour+2*time.Minute+3*time.Second))
		})

		Convey("when ko", func() {
			var res Duration
			err := res.UnmarshalCSV("oups")
			So(err, ShouldBeError, `time: invalid duration "oups"`)
			So(res, ShouldBeZeroValue)
		})
	})
}

func TestDate(t *testing.T) {
	Convey("marshal", t, func() {
		res, err := Date(time.Date(2023, 12, 2, 3, 4, 5, 6, time.UTC)).MarshalCSV()
		So(err, ShouldBeNil)
		So(res, ShouldEqual, "2023-12-02")
	})

	Convey("unmarshal", t, func() {
		Convey("when ok", func() {
			var res Date
			err := res.UnmarshalCSV("2023-12-02")
			So(err, ShouldBeNil)
			So(res, ShouldEqual, Date(time.Date(2023, 12, 2, 0, 0, 0, 0, time.UTC)))
		})

		Convey("when ko", func() {
			var res Date
			err := res.UnmarshalCSV("oups")
			So(err, ShouldBeError, `parsing time "oups" as "2006-01-02": cannot parse "oups" as "2006"`)
			So(res, ShouldBeZeroValue)
		})
	})
}
