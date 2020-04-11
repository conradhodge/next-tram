package traveline_test

import (
	"testing"

	"github.com/conradhodge/next-tram/src/lib/traveline"
)

func TestInvalidTimeFoundError(t *testing.T) {
	err := traveline.InvalidTimeFoundError{
		Time:   "unknown",
		Reason: "Cannot convert \"unknown\" to time",
	}

	expectedError := "Invalid departure time \"unknown\" found: Cannot convert \"unknown\" to time"

	if err.Error() != expectedError {
		t.Fatalf("Expected error:\n%s\ngot:\n%s", expectedError, err.Error())
	}
}

func TestNoTimesFoundError(t *testing.T) {
	err := traveline.NoTimesFoundError{}

	expectedError := "No next departure times found"

	if err.Error() != expectedError {
		t.Fatalf("Expected error:\n%s\ngot:\n%s", expectedError, err.Error())
	}
}
