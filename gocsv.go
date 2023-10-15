package gocsv

import (
	"fmt"

	"github.com/sbiemont/gocsv/internal"
)

// Decode a csv struct into the given type of data
func Decode[T any](data [][]string) ([]T, error) {
	// Prepare data
	ct, err := internal.NewCacheTags[T]()
	if err != nil {
		return nil, err
	}
	cm := internal.NewCacheUnmarshaler()
	res := make([]T, len(data))

	// Read all
	for i, row := range data {
		err := internal.Unmarshal(ct, cm, row, &res[i])
		if err != nil {
			return nil, fmt.Errorf("row %d: %w", i, err)
		}
	}
	return res, nil
}

// Encode into a csv struct
func Encode[T any](data []T) ([][]string, error) {
	// Prepare data
	ct, err := internal.NewCacheTags[T]()
	if err != nil {
		return nil, err
	}
	cm := internal.NewCacheMarshaler()
	res := make([][]string, len(data))

	// Read all
	for i, item := range data {
		row, err := internal.Marshal(ct, cm, item)
		if err != nil {
			return nil, fmt.Errorf("row %d: %w", i, err)
		}
		res[i] = row
	}
	return res, nil
}
