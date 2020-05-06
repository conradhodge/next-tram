package traveline_test

import (
	"testing"

	"github.com/conradhodge/next-tram/src/lib/traveline"
)

func TestNoTimesFoundError(t *testing.T) {
	err := traveline.NoTimesFoundError{}

	expectedError := "No next departure times found"

	if err.Error() != expectedError {
		t.Fatalf("Expected error:\n%s\ngot:\n%s", expectedError, err.Error())
	}
}
