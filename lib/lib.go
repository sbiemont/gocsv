package lib

// Marshaler defines the unique method for marshaling a CSV field
type Marshaler interface {
	MarshalCSV() (string, error)
}

// Unmarshaler defines the unique method for unmarshaling a CSV field
type Unmarshaler interface {
	UnmarshalCSV(string) error
}
