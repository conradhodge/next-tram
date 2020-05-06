package transport_test

import (
	"testing"

	"github.com/conradhodge/next-tram/src/lib/transport"
)

func TestInvalidTimeFoundError(t *testing.T) {
	err := transport.InvalidTimeFoundError{
		Time:   "unknown",
		Reason: "Cannot convert \"unknown\" to time",
	}

	expectedError := "Invalid departure time \"unknown\" found: Cannot convert \"unknown\" to time"

	if err.Error() != expectedError {
		t.Fatalf("Expected error:\n%s\ngot:\n%s", expectedError, err.Error())
	}
}
