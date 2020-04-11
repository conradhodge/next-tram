package traveline

import "fmt"

// InvalidTimeFoundError indicates that an invalid departure time was found
type InvalidTimeFoundError struct {
	Time   string
	Reason string
}

func (e InvalidTimeFoundError) Error() string {
	return fmt.Sprintf("Invalid departure time \"%s\" found: %s", e.Time, e.Reason)
}

// NoTimesFoundError indicates that no departure times can be found
type NoTimesFoundError struct{}

func (e NoTimesFoundError) Error() string {
	return fmt.Sprint("No next departure times found")
}
