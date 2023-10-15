package internal

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCacheTags(t *testing.T) {
	Convey("custom", t, func() {
		type custom struct {
			Prop1 float64 `csv:"10,omitempty"`
			Prop2 int     `csv:"20"`
			Prop3 string
		}

		cache, err := NewCacheTags[custom]()
		So(err, ShouldBeNil)
		So(cache, ShouldResemble, CacheTags[custom]{
			0: {
				col:       10,
				omitEmpty: true,
			},
			1: {
				col:       20,
				omitEmpty: false,
			},
		})
	})

	Convey("when ko", t, func() {
	})
}
