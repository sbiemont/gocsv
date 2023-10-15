package internal

import (
	"testing"
	"time"

	"github.com/sbiemont/gocsv/lib"
	. "github.com/smartystreets/goconvey/convey"
)

type customMarshal struct {
	private string
}

func (it customMarshal) MarshalCSV() (string, error) {
	return it.private, nil
}

func testMarshal[T any](ct CacheTags[T], cm CacheMarshaler, item T) []string {
	res, err := Marshal(ct, cm, item)
	So(err, ShouldBeNil)
	return res
}

func TestMarshalCSV(t *testing.T) {
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
		cm := NewCacheMarshaler()
		ptr := true
		res := testMarshal(ct, cm, testStruct{
			Bool1:     true,
			Bool2:     false,
			BoolEmpty: false,
			Ptr:       &ptr,
			PtrNil:    nil,
		})
		So(res, ShouldResemble, []string{
			"true",
			"false",
			"false",
			"true",
			"",
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
		cm := NewCacheMarshaler()
		ptr := "ptr"
		ts := testMarshal(ct, cm, testStruct{
			Str1:     "str",
			Str2:     "",
			StrEmpty: "",
			Ptr:      &ptr,
			PtrNil:   nil,
		})
		So(ts, ShouldResemble, []string{
			"str",
			"",
			"",
			"ptr",
			"",
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
		cm := NewCacheMarshaler()
		ptr := 255
		ts := testMarshal(ct, cm, testStruct{
			Int:   42,
			Int8:  -127,
			Int16: 32767,
			Int32: 43,
			Int64: 44,
			Ptr:   &ptr,
		})
		So(ts, ShouldResemble, []string{
			"42",
			"-127",
			"32767",
			"43",
			"44",
			"255",
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
		cm := NewCacheMarshaler()
		ptr := 42.0
		ts := testMarshal(ct, cm, testStruct{
			Flt32:  3.40282,
			Flt64:  1.79769,
			Ptr:    &ptr,
			PtrNil: nil,
		})
		So(ts, ShouldResemble, []string{
			"3.402820",
			"1.797690",
			"42.000000",
			"",
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
		cm := NewCacheMarshaler()
		ts := testMarshal(ct, cm, testStruct{
			Dur1:   lib.Duration(time.Hour),
			Dur2:   &ptr,
			PtrNil: nil,
		})
		So(ts, ShouldResemble, []string{
			"1h0m0s",
			"1m0s",
			"",
		})

		// Use cache
		ptr2 := lib.Duration(2 * time.Minute)
		ts2 := testMarshal(ct, cm, testStruct{
			Dur1:   lib.Duration(2 * time.Hour),
			Dur2:   &ptr2,
			PtrNil: nil,
		})
		So(ts2, ShouldResemble, []string{
			"2h0m0s",
			"2m0s",
			"",
		})
	})

	Convey("custom struct", t, func() {
		type testStruct struct {
			Custom1   customMarshal  `csv:"0"`
			Custom2   *customMarshal `csv:"1"`
			CustomNil *customMarshal `csv:"2,omitempty"`
		}

		ptr := customMarshal{private: "custom 2"}
		ct, err := NewCacheTags[testStruct]()
		So(err, ShouldBeNil)
		cm := NewCacheMarshaler()
		ts := testMarshal(ct, cm, testStruct{
			Custom1:   customMarshal{private: "custom 1"},
			Custom2:   &ptr,
			CustomNil: nil,
		})
		So(ts, ShouldResemble, []string{
			"custom 1",
			"custom 2",
			"",
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
		cm := NewCacheMarshaler()
		ts := testMarshal(ct, cm, testStruct{
			Time1:     time.Date(2023, 2, 3, 10, 11, 12, 0, time.UTC),
			Time2:     &ptr,
			TimeEmpty: time.Time{},
			TimeNil:   nil,
		})
		So(ts, ShouldResemble, []string{
			"2023-02-03T10:11:12Z",
			"2023-03-04T00:00:00Z",
			"0001-01-01T00:00:00Z", // should use ptr
			"",
		})
	})
}
