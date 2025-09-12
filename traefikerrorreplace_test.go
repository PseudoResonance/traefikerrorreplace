package cloudflarewarp_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	plugin "github.com/PseudoResonance/traefikerrorreplace"
)

func TestNew(t *testing.T) {
	cfg := plugin.CreateConfig()
	cfg.MatchStatus = []int{500, 503}
	cfg.ReplaceStatus = 418

	ctx := context.Background()
	testCases := []struct {
		desc         string
		sentCode     int
		expectedCode int
	}{
		{
			desc:         "Test Status Good 200",
			sentCode:     200,
			expectedCode: 200,
		},
		{
			desc:         "Test Status Good 404",
			sentCode:     404,
			expectedCode: 404,
		},
		{
			desc:         "Test Status Bad 500",
			sentCode:     500,
			expectedCode: 418,
		},
		{
			desc:         "Test Status Bad 503",
			sentCode:     503,
			expectedCode: 418,
		},
	}
	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			recorder := httptest.NewRecorder()

			next := http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) {
				rw.WriteHeader(test.sentCode)
			})

			handler, err := plugin.New(ctx, next, cfg, "traefikerrorreplace")
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
			if err != nil {
				t.Fatal(err)
			}

			handler.ServeHTTP(recorder, req)

			if recorder.Result().StatusCode != test.expectedCode {
				t.Errorf("invalid status: %v, expected: %v", strconv.Itoa(recorder.Result().StatusCode), strconv.Itoa(test.expectedCode))
				return
			}
		})
	}
}
