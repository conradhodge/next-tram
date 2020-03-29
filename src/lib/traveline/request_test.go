package traveline_test

// The following comment is used by 'go generate ./...' command. DO NOT DELETE!!!
//go:generate mockgen -destination ../mock/mock_traveline/mock_traveline.go github.com/conradhodge/next-tram/src/lib/traveline API

import (
	"errors"
	"testing"
	"time"

	"github.com/conradhodge/next-tram/src/lib/mock/mock_traveline"
	"github.com/conradhodge/next-tram/src/lib/traveline"
	"github.com/golang/mock/gomock"
)

func TestGetNextTram(t *testing.T) {
	now := time.Now()
	nextTramTime, _ := time.Parse(time.RFC3339, "2020-03-30T12:34:56.911+01:00")

	tests := []struct {
		name                 string
		naptanCode           string
		when                 time.Time
		request              string
		response             string
		nextTramTime         time.Time
		buildError           error
		parseError           error
		sendError            error
		expectedNextTramTime string
		expectedError        error
	}{
		{
			name:                 "Valid next tram time",
			naptanCode:           "123456789",
			when:                 now,
			request:              "<Siri><ServiceRequest></ServiceRequest></Siri>",
			response:             "<Siri><ServiceDelivery></ServiceDelivery></Siri>",
			nextTramTime:         nextTramTime,
			expectedNextTramTime: "12:34PM",
		},
		{
			name:          "Error building request",
			naptanCode:    "123456789",
			when:          now,
			buildError:    errors.New("Error building request"),
			expectedError: errors.New("Error building request"),
		},
		{
			name:          "Error sending request",
			naptanCode:    "123456789",
			when:          now,
			request:       "<Siri><ServiceRequest></ServiceRequest></Siri>",
			response:      "<Siri><ServiceDelivery></ServiceDelivery></Siri>",
			sendError:     errors.New("Error sending request"),
			expectedError: errors.New("Error sending request"),
		},
		{
			name:          "Error parsing request",
			naptanCode:    "123456789",
			when:          now,
			request:       "<Siri><ServiceRequest></ServiceRequest></Siri>",
			response:      "<Siri><ServiceDelivery></ServiceDelivery></Siri>",
			parseError:    errors.New("Error parsing response"),
			expectedError: errors.New("Error parsing response"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// We need a controller
			// https://github.com/golang/mock
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Mock the API
			mockAPI := mock_traveline.NewMockAPI(ctrl)

			mockAPI.
				EXPECT().
				BuildRequest(gomock.Any(), gomock.Eq(test.naptanCode), gomock.Eq(test.when)).
				Return(test.request, test.buildError).
				AnyTimes()

			mockAPI.
				EXPECT().
				Send(gomock.Eq(test.request)).
				Return(test.response, test.sendError).
				AnyTimes()

			mockAPI.
				EXPECT().
				ParseResponse(gomock.Eq(test.response)).
				Return(test.nextTramTime, test.parseError).
				AnyTimes()

			req := traveline.Request{API: mockAPI}

			nextTramTime, err := req.GetNextTram(test.naptanCode, test.when)

			if test.expectedError != nil {
				if err == nil {
					t.Fatalf("Expected error '%s'; got no error", test.expectedError)
				}
				if err.Error() != test.expectedError.Error() {
					t.Fatalf("Expected error '%s'; got '%s'", test.expectedError.Error(), err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("Expected no error; got '%s'", err)
				}
			}

			if nextTramTime != test.expectedNextTramTime {
				t.Fatalf("Expected response: %s, got: %s", test.expectedNextTramTime, nextTramTime)
			}
		})
	}
}
