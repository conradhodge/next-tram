package transport_test

// The following comment is used by 'go generate ./...' command. DO NOT DELETE!!!
//go:generate mockgen -destination ../mock/mock_traveline/mock_traveline.go github.com/conradhodge/next-tram/src/lib/traveline API

import (
	"errors"
	"testing"
	"time"

	"github.com/conradhodge/next-tram/src/lib/matcher"
	"github.com/conradhodge/next-tram/src/lib/mock/mock_traveline"
	"github.com/conradhodge/next-tram/src/lib/transport"
	"github.com/conradhodge/next-tram/src/lib/traveline"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
)

func TestGetNextTram(t *testing.T) {
	now := time.Now()
	nextTramTime, _ := time.Parse(time.RFC3339, "2020-03-30T12:34:56.911+01:00")
	differsNextTramTime, _ := time.Parse(time.RFC3339, "2020-03-30T12:37:56.911+01:00")

	tests := []struct {
		name           string
		naptanCode     string
		when           time.Time
		parseResult    *traveline.MonitoredVehicleJourney
		buildError     error
		parseError     error
		sendError      error
		expectedError  error
		expectedResult *transport.NextTramInfo
	}{
		{
			name:       "Aimed departure time differs from expected departure time",
			naptanCode: "123456789",
			when:       now,
			parseResult: &traveline.MonitoredVehicleJourney{
				VehicleMode:       "magic carpet",
				PublishedLineName: "flying",
				DirectionName:     "Xanadu",
				MonitoredCall: struct {
					AimedDepartureTime    string "xml:\"AimedDepartureTime\""
					ExpectedDepartureTime string "xml:\"ExpectedDepartureTime\""
				}{
					AimedDepartureTime:    "2020-03-30T12:34:56.911+01:00",
					ExpectedDepartureTime: "2020-03-30T12:37:56.911+01:00",
				},
			},
			expectedResult: &transport.NextTramInfo{
				VehicleMode:           "magic carpet",
				LineName:              "flying",
				DirectionName:         "Xanadu",
				AimedDepartureTime:    &nextTramTime,
				ExpectedDepartureTime: &differsNextTramTime,
			},
		},
		{
			name:       "Aimed departure time matches expected departure time",
			naptanCode: "123456789",
			when:       now,
			parseResult: &traveline.MonitoredVehicleJourney{
				VehicleMode:       "magic carpet",
				PublishedLineName: "flying",
				DirectionName:     "Xanadu",
				MonitoredCall: struct {
					AimedDepartureTime    string "xml:\"AimedDepartureTime\""
					ExpectedDepartureTime string "xml:\"ExpectedDepartureTime\""
				}{
					AimedDepartureTime:    "2020-03-30T12:34:56.911+01:00",
					ExpectedDepartureTime: "2020-03-30T12:34:56.911+01:00",
				},
			},
			expectedResult: &transport.NextTramInfo{
				VehicleMode:           "magic carpet",
				LineName:              "flying",
				DirectionName:         "Xanadu",
				AimedDepartureTime:    &nextTramTime,
				ExpectedDepartureTime: &nextTramTime,
			},
		},
		{
			name:       "No expected departure time",
			naptanCode: "123456789",
			when:       now,
			parseResult: &traveline.MonitoredVehicleJourney{
				VehicleMode:       "magic carpet",
				PublishedLineName: "flying",
				DirectionName:     "Xanadu",
				MonitoredCall: struct {
					AimedDepartureTime    string "xml:\"AimedDepartureTime\""
					ExpectedDepartureTime string "xml:\"ExpectedDepartureTime\""
				}{
					AimedDepartureTime: "2020-03-30T12:34:56.911+01:00",
				},
			},
			expectedResult: &transport.NextTramInfo{
				VehicleMode:        "magic carpet",
				LineName:           "flying",
				DirectionName:      "Xanadu",
				AimedDepartureTime: &nextTramTime,
			},
		},
		{
			name:       "Invalid aimed departure time",
			naptanCode: "123456789",
			when:       now,
			parseResult: &traveline.MonitoredVehicleJourney{
				VehicleMode:       "magic carpet",
				PublishedLineName: "flying",
				DirectionName:     "Xanadu",
				MonitoredCall: struct {
					AimedDepartureTime    string "xml:\"AimedDepartureTime\""
					ExpectedDepartureTime string "xml:\"ExpectedDepartureTime\""
				}{
					AimedDepartureTime: "bongo",
				},
			},
			expectedError: &transport.InvalidTimeFoundError{
				Time:   "bongo",
				Reason: `parsing time "bongo" as "2006-01-02T15:04:05Z07:00": cannot parse "bongo" as "2006"`,
			},
		},
		{
			name:       "Invalid expected departure time",
			naptanCode: "123456789",
			when:       now,
			parseResult: &traveline.MonitoredVehicleJourney{
				VehicleMode:       "magic carpet",
				PublishedLineName: "flying",
				DirectionName:     "Xanadu",
				MonitoredCall: struct {
					AimedDepartureTime    string "xml:\"AimedDepartureTime\""
					ExpectedDepartureTime string "xml:\"ExpectedDepartureTime\""
				}{
					AimedDepartureTime:    "2020-03-30T12:34:56.911+01:00",
					ExpectedDepartureTime: "bango",
				},
			},
			expectedError: &transport.InvalidTimeFoundError{
				Time:   "bango",
				Reason: `parsing time "bango" as "2006-01-02T15:04:05Z07:00": cannot parse "bango" as "2006"`,
			},
		},
		{
			name:          "Error building request",
			naptanCode:    "123456789",
			when:          now,
			buildError:    errors.New("build fail"),
			expectedError: errors.New("build fail"),
		},
		{
			name:          "Error sending request",
			naptanCode:    "123456789",
			when:          now,
			sendError:     errors.New("send fail"),
			expectedError: errors.New("send fail"),
		},
		{
			name:          "Error parsing response",
			naptanCode:    "123456789",
			when:          now,
			parseError:    errors.New("parse fail"),
			expectedError: errors.New("parse fail"),
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
				BuildServiceRequest(matcher.IsGUID(), gomock.Eq(test.naptanCode), gomock.Eq(test.when)).
				Return("<request/>", test.buildError).
				AnyTimes()
			mockAPI.
				EXPECT().
				Send(gomock.Eq("<request/>")).
				Return("<response/>", test.sendError).
				AnyTimes()
			mockAPI.
				EXPECT().
				ParseServiceDelivery(gomock.Eq("<response/>")).
				Return(test.parseResult, test.parseError).
				AnyTimes()

			req := transport.NewTraveline(mockAPI)

			result, err := req.GetNextTramTime(test.naptanCode, test.when)

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

			if diff := cmp.Diff(test.expectedResult, result); diff != "" {
				t.Errorf("GetNextTramTime() (-want +got):\n%s", diff)
			}
		})
	}
}
