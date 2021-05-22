package main_test

// The following comment is used by 'go generate ./...' command. DO NOT DELETE!!!
//go:generate mockgen -destination ../../mock/mock_transport/mock_transport.go github.com/conradhodge/travel-api-client/transport API

import (
	"errors"
	"testing"
	"time"

	main "github.com/conradhodge/next-tram/lambda/get-next-tram"
	"github.com/conradhodge/next-tram/mock/mock_transport"
	"github.com/conradhodge/travel-api-client/transport"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetNextTram(t *testing.T) {
	// Replace the TimeNow function temporarily
	originalTimeNow := main.TimeNow
	defer func() { main.TimeNow = originalTimeNow }()

	main.TimeNow = func() time.Time {
		now, _ := time.Parse(time.RFC3339, "2021-05-22T20:00:00+01:00")
		return now
	}

	firstTime, _ := time.Parse(time.RFC3339, "2021-05-22T12:34:56Z")
	secondTime, _ := time.Parse(time.RFC3339, "2021-05-22T12:36:56Z")

	tests := []struct {
		name               string
		aimedDeparture     *time.Time
		expectedDeparture  *time.Time
		getNextTramTimeErr error
		expectedMessage    string
	}{
		{
			name:              "Aimed matches expected time",
			aimedDeparture:    &firstTime,
			expectedDeparture: &firstTime,
			expectedMessage:   "Your next flying magic carpet to Xanadu is due at 12:34PM",
		},
		{
			name:              "Aimed different than expected time",
			aimedDeparture:    &firstTime,
			expectedDeparture: &secondTime,
			expectedMessage:   "Your next flying magic carpet to Xanadu is due at 12:34PM, but is expected at 12:36PM",
		},
		{
			name:               "GetNextTramTime error",
			getNextTramTimeErr: errors.New("bang"),
			expectedMessage:    "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// We need a controller
			// https://github.com/golang/mock
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Mock the API
			mockAPI := mock_transport.NewMockAPI(ctrl)

			nextTramInfo := &transport.DepartureInfo{
				VehicleMode:           "magic carpet",
				LineName:              "flying",
				DirectionName:         "Xanadu",
				AimedDepartureTime:    tc.aimedDeparture,
				ExpectedDepartureTime: tc.expectedDeparture,
			}

			mockAPI.
				EXPECT().
				GetNextDepartureTime("123456", main.TimeNow()).
				Return(nextTramInfo, tc.getNextTramTimeErr).
				AnyTimes()

			message, err := main.GetNextTram(mockAPI, "123456")

			if tc.getNextTramTimeErr != nil {
				assert.EqualError(t, err, tc.getNextTramTimeErr.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, message, tc.expectedMessage)
		})
	}
}
