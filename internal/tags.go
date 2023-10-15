package internal

import (
	"reflect"
	"slices"
	"strconv"
	"strings"
)

type tag struct {
	col       int
	omitEmpty bool
}

// CacheTags stores the tag data for the ith field
type CacheTags[T any] map[int]tag

// NewCacheTags init the cache using a given type of struct
func NewCacheTags[T any]() (CacheTags[T], error) {
	var item T
	datas := make(map[int]tag)
	typ := reflect.Indirect(reflect.ValueOf(item)).Type()
	for i := 0; i < typ.NumField(); i++ {
		// Get "csv" info => parse `csv:"tag0,tag1,..,tagN"`
		csvTag, ok := typ.Field(i).Tag.Lookup("csv")
		if ok {
			// Fetch tags (the first one is the position)
			tags := strings.Split(csvTag, ",")
			pos, err := strconv.Atoi(tags[0])
			if err != nil {
				return nil, err
			}

			// Out of bounds (omit empty or unknown position)
			omitEmpty := slices.Contains(tags[1:], "omitempty")
			datas[i] = tag{
				col:       pos,
				omitEmpty: omitEmpty,
			}
		}
	}
	return datas, nil
}

// tag returns the ith tag (if found)
func (cache CacheTags[T]) tag(i int) (tag, bool) {
	t, ok := cache[i]
	return t, ok
}

// maxCol found in tags
func (cache CacheTags[T]) maxCol() int {
	maxCol := -1
	for _, data := range cache {
		if data.col > maxCol {
			maxCol = data.col
		}
	}
	return maxCol
}
