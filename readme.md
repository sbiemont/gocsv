# gocsv

This project aims to simplify the reading of a csv file content.

It uses an internal cache for reading struct definitions but does not support super fast processing.

## Usage

### Mapping

Define a mapping on the csv columns ; just provide:

* the column number
* the optional `omitempty` property if the field can be `nil`

```go
// define the struct types
type row struct {
  ID    int      `csv:"0"`           // Identifier at column 0
  Name  string   `csv:"3"`           // Name at column 3
  Value *float64 `csv:"4,omitempty"` // Optional value at column 4
}
```

### Decode

To decode a set of rows, call the `gocsv.Decode[T]` function (and provide the generic type of row to decode).

```go
// Read the whole csv file using the built-in reader
file, _ := os.Open(filename)
records, _ := csv.NewReader(file).ReadAll()

// Decode them (except the titles row)
rows, _ := gocsv.Decode[row](records[1:])
```

### Encode

To encode, call the `gocsv.Encode[T]` function (the output type `T` is optional and will be deduced from the parameter)

```go
// Encode records
value := 3.14159
records, err := gocsv.Encode([]row{
  {ID: 1, Name: "John Doe", Value: nil},
  {ID: 2, Name: "Jane Doe", Value: &value},
})

// Write them into a file
file, _ := os.Create("test.csv")
_ = csv.NewWriter(file).WriteAll(records)
```

## Example

See [example](https:..github.com/sbiemont/gocsv/example) directory for more examples

## Custom fields

You can specify your own types that respond to the `MarshalCSV` and `UnmarshalCSV`. Then, use it as a field definition in your mapping structure.

```go
// Date is a mapping on a time.Time with format "YYYY-MM-DD"
type Date time.Time

// MarshalCSV exports the date
func (it Date) MarshalCSV() (string, error) {
  return time.Time(it).Format(time.DateOnly), nil
}

// UnmarshalCSV imports the date
func (it *Date) UnmarshalCSV(s string) error {
  tm, err := time.Parse(time.DateOnly, s)
  *it = Date(tm)
  return err
}
```
