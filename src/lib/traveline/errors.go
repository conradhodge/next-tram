package traveline

import "fmt"

// NoTimesFoundError indicates that no departure times can be found
type NoTimesFoundError struct{}

func (e *NoTimesFoundError) Error() string {
	return fmt.Sprint("No next departure times found")
}
