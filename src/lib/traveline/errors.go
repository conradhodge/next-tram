package traveline

// NoTimesFoundError indicates that no departure times can be found
type NoTimesFoundError struct{}

func (e NoTimesFoundError) Error() string {
	return "No next departure times found"
}
