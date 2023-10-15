package internal

import (
	"testing"
	"time"

	"github.com/sbiemont/gocsv/lib"
	. "github.com/smartystreets/goconvey/convey"
)

type customUnmarshal struct {
	private string
}

func (it *customUnmarshal) UnmarshalCSV(s string) error {
	it.private = s
	return nil
}

func testUnmarshal[T any](ct CacheTags[T], cm CacheUnmarshaler, inputs []string) T {
	var ts T
	err := Unmarshal(ct, cm, inputs, &ts)
	So(err, ShouldBeNil)
	return ts
}

func TestUnmarshalCSV(t *testing.T) {
	Convey("boolean", t, func() {
		type testStruct struct {
			Bool1     bool  `csv:"0"`
			Bool2     bool  `csv:"1"`
			BoolEmpty bool  `csv:"2,omitempty"`
			Ptr       *bool `csv:"3"`
			PtrNil    *bool `csv:"4,omitempty"`
		}

		ct, err := NewCacheTags[testStruct]()
		So(err, ShouldBeNil)
		cm := NewCacheUnmarshaler()
		ts := testUnmarshal(ct, cm, []string{
			"true",
			"false",
			"",
			"true",
			"",
		})
		ptr := true
		So(ts, ShouldResemble, testStruct{
			Bool1:     true,
			Bool2:     false,
			BoolEmpty: false,
			Ptr:       &ptr,
			PtrNil:    nil,
		})
	})

	Convey("string", t, func() {
		type testStruct struct {
			Str1     string  `csv:"0"`
			Str2     string  `csv:"1"`
			StrEmpty string  `csv:"2,omitempty"`
			Ptr      *string `csv:"3"`
			PtrNil   *string `csv:"4,omitempty"`
		}

		ct, err := NewCacheTags[testStruct]()
		So(err, ShouldBeNil)
		cm := NewCacheUnmarshaler()
		ts := testUnmarshal(ct, cm, []string{
			"str",
			"",
			"",
			"ptr",
			"",
		})
		ptr := "ptr"
		So(ts, ShouldResemble, testStruct{
			Str1:     "str",
			Str2:     "",
			StrEmpty: "",
			Ptr:      &ptr,
			PtrNil:   nil,
		})
	})

	Convey("integer", t, func() {
		type testStruct struct {
			Int   int   `csv:"0"`
			Int8  int8  `csv:"1"`
			Int16 int16 `csv:"2"`
			Int32 int32 `csv:"3"`
			Int64 int64 `csv:"4"`
			Ptr   *int  `csv:"5"`
		}

		ct, err := NewCacheTags[testStruct]()
		So(err, ShouldBeNil)
		cm := NewCacheUnmarshaler()
		ts := testUnmarshal(ct, cm, []string{
			"42",
			"-127",
			"32767",
			"43",
			"44",
			"0xff", // 255
		})
		ptr := 255
		So(ts, ShouldResemble, testStruct{
			Int:   42,
			Int8:  -127,
			Int16: 32767,
			Int32: 43,
			Int64: 44,
			Ptr:   &ptr,
		})
	})

	Convey("float", t, func() {
		type testStruct struct {
			Flt32  float32  `csv:"0"`
			Flt64  float64  `csv:"1"`
			Ptr    *float64 `csv:"2"`
			PtrNil *float64 `csv:"3,omitempty"`
		}

		ct, err := NewCacheTags[testStruct]()
		So(err, ShouldBeNil)
		cm := NewCacheUnmarshaler()
		ts := testUnmarshal(ct, cm, []string{
			"3.40282e+38",
			"1.79769e+308",
			"42",
			"",
		})
		ptr := 42.0
		So(ts, ShouldResemble, testStruct{
			Flt32:  3.40282e+38,
			Flt64:  1.79769e+308,
			Ptr:    &ptr,
			PtrNil: nil,
		})
	})

	Convey("csv duration", t, func() {
		type testStruct struct {
			Dur1   lib.Duration  `csv:"0"`
			Dur2   *lib.Duration `csv:"1"`
			PtrNil *lib.Duration `csv:"2,omitempty"`
		}

		ptr := lib.Duration(time.Minute)
		ct, err := NewCacheTags[testStruct]()
		So(err, ShouldBeNil)
		cm := NewCacheUnmarshaler()
		ts := testUnmarshal(ct, cm, []string{
			"1h",
			"1m",
			"",
		})
		So(ts, ShouldResemble, testStruct{
			Dur1:   lib.Duration(time.Hour),
			Dur2:   &ptr,
			PtrNil: nil,
		})

		// Use cache
		ts2 := testUnmarshal(ct, cm, []string{
			"2h",
			"2m",
			"",
		})
		ptr2 := lib.Duration(2 * time.Minute)
		So(ts2, ShouldResemble, testStruct{
			Dur1:   lib.Duration(2 * time.Hour),
			Dur2:   &ptr2,
			PtrNil: nil,
		})
	})

	Convey("custom struct", t, func() {
		type testStruct struct {
			Custom1   customUnmarshal  `csv:"0"`
			Custom2   *customUnmarshal `csv:"1"`
			CustomNil *customUnmarshal `csv:"2,omitempty"`
		}

		ptr := customUnmarshal{private: "custom 2"}
		ct, err := NewCacheTags[testStruct]()
		So(err, ShouldBeNil)
		cm := NewCacheUnmarshaler()
		ts := testUnmarshal(ct, cm, []string{
			"custom 1",
			"custom 2",
			"",
		})
		So(ts, ShouldResemble, testStruct{
			Custom1:   customUnmarshal{private: "custom 1"},
			Custom2:   &ptr,
			CustomNil: nil,
		})
	})

	Convey("when using marshal text", t, func() {
		type testStruct struct {
			Time1     time.Time  `csv:"0"`
			Time2     *time.Time `csv:"1"`
			TimeEmpty time.Time  `csv:"2,omitempty"`
			TimeNil   *time.Time `csv:"3,omitempty"`
		}

		ptr := time.Date(2023, 3, 4, 0, 0, 0, 0, time.UTC)
		ct, err := NewCacheTags[testStruct]()
		So(err, ShouldBeNil)
		cm := NewCacheUnmarshaler()
		ts := testUnmarshal(ct, cm, []string{
			"2023-02-03T10:11:12Z",
			"2023-03-04T00:00:00Z",
			"",
			"",
		})
		So(ts, ShouldResemble, testStruct{
			Time1:     time.Date(2023, 2, 3, 10, 11, 12, 0, time.UTC),
			Time2:     &ptr,
			TimeEmpty: time.Time{},
			TimeNil:   nil,
		})
	})
}
