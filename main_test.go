package main

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegsiter(t *testing.T) {

	tests := []struct {
		description string
		route       string
		method      string

		// Expected output
		expectedError bool
		expectedCode  int
		expectedBody  string
	}{
		{
			"register route",
			"/api/v1/register",
			"POST",
			false,
			200,
			"ok",
		},
	}

	app := Setup()
	for _, test := range tests {
		req, _ := http.NewRequest(
			test.method,
			test.route,
			nil,
		)
		res, err := app.Test(req, -1)
		assert.Equal(t, test.expectedError, err != nil, test.description)

		if test.expectedError {
			continue
		}

		body, err := io.ReadAll(res.Body)

		assert.Nilf(t, err, test.description)

		assert.Equal(t, test.expectedBody, string(body), test.description)
	}
}

func TestLogin(t *testing.T) {

}
