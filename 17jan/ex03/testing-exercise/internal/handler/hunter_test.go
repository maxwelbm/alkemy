package handler_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testdoubles/internal/handler"
	"testdoubles/internal/hunter"
	"testdoubles/internal/prey"

	"testing"

	"github.com/stretchr/testify/require"
)

func TestHunterConfigurePreyHandler(t *testing.T) {
	type arrange struct {
		mockPrey func() *prey.PreyStub
	}
	type input struct {
		request  func() *http.Request
		response *httptest.ResponseRecorder
	}
	type output struct {
		code    int
		body    string
		headers http.Header
	}
	type testCase struct {
		name    string
		arrange arrange
		input   input
		output  output
	}

	testCases := []testCase{
		{
			name: "case 1: success to configure the prey",
			arrange: arrange{
				mockPrey: func() *prey.PreyStub {
					return prey.NewPreyStub()
				},
			},
			input: input{
				request: func() *http.Request {
					r := httptest.NewRequest(http.MethodPost, "/prey/configure-prey",
						strings.NewReader(`{"speed":10.0,"position":{"X": 1.0,"Y": 2.0,"Z": 3.0}}`),
					)
					r.Header.Set("Content-Type", "application/json")
					return r
				},
				response: httptest.NewRecorder(),
			},
			output: output{
				code:    http.StatusOK,
				body:    `{"message":"prey configured","data":null}`,
				headers: http.Header{"Content-Type": []string{"application/json"}},
			},
		},
		{
			name: "case 2: invalid request body",
			arrange: arrange{
				mockPrey: func() *prey.PreyStub { return nil },
			},
			input: input{
				request: func() *http.Request {
					r := httptest.NewRequest(http.MethodPost, "/prey/configure-prey",
						strings.NewReader(`{"speed":"invalid","position":{"X": 1.0,"Y": 2.0}}`),
					)
					r.Header.Set("Content-Type", "application/json")
					return r
				},
				response: httptest.NewRecorder(),
			},
			output: output{
				code: http.StatusBadRequest,
				body: fmt.Sprintf(
					`{"status":"%s","message":"%s"}`,
					http.StatusText(http.StatusBadRequest),
					"invalid request body",
				),
				headers: http.Header{"Content-Type": []string{"application/json"}},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockPrey := tc.arrange.mockPrey()
			// - handler
			hd := handler.NewHunter(nil, mockPrey)
			hdFunc := hd.ConfigurePrey

			// act
			hdFunc(tc.input.response, tc.input.request())

			// assert
			require.Equal(t, tc.output.code, tc.input.response.Code)
			require.JSONEq(t, tc.output.body, tc.input.response.Body.String())
			require.Equal(t, tc.output.headers, tc.input.response.Header())
		})
	}
}

// Tests for Hunter ConfigureHunter handler.
func TestHunterConfigureHunterHandler(t *testing.T) {
	type arrange struct {
		mockPrey func() *prey.PreyStub
	}
	type input struct {
		request  func() *http.Request
		response *httptest.ResponseRecorder
	}
	type output struct {
		code    int
		body    string
		headers http.Header
	}
	type testCase struct {
		name    string
		arrange arrange
		input   input
		output  output
	}
	testCases := []testCase{
		{
			name: "test - invalid request body",
			arrange: arrange{
				mockPrey: func() *prey.PreyStub {
					return prey.NewPreyStub()
				},
			},
			input: input{
				request: func() *http.Request {
					body := strings.NewReader(`{"speed":"invalid","position":{"X": 1.0,"Y": 2.0}}`)
					r := httptest.NewRequest(http.MethodPost, "/hunter/configure-hunter", body)
					r.Header.Set("Content-Type", "application/json")
					return r
				},
				response: httptest.NewRecorder(),
			},
			output: output{
				code:    http.StatusOK,
				body:    `{"message":"invalid request body","data":null}`,
				headers: http.Header{"Content-Type": []string{"application/json"}},
			},
		},
		{
			name: "test - hunter created successfully",
			arrange: arrange{
				mockPrey: func() *prey.PreyStub {
					return prey.NewPreyStub()
				},
			},
			input: input{
				request: func() *http.Request {
					body := strings.NewReader(`{"speed":"15.0","position":{"X": 1.0,"Y": 2.0}}`)
					r := httptest.NewRequest(http.MethodPost, "/hunter/configure-hunter", body)
					r.Header.Set("Content-Type", "application/json")
					return r
				},
				response: httptest.NewRecorder(),
			},
			output: output{
				code:    http.StatusOK,
				body:    `{"message":"hunter configured","data":null}`,
				headers: http.Header{"Content-Type": []string{"application/json"}},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockHunter := hunter.NewHunterMock()
			hd := handler.NewHunter(mockHunter, nil)

			hunterFunc := hd.ConfigureHunter()
			hunterFunc(tc.input.response, tc.input.request())

			require.Equal(t, tc.output.code, tc.input.response.Code)
			require.JSONEq(t, tc.output.body, tc.input.response.Body.String())
			require.Equal(t, tc.output.headers, tc.input.response.Header())
		})
	}
}

// Tests for Hunter Hunt handler.
func TestHunterHuntHandler(t *testing.T) {
	type arrange struct {
		mockHunter func() *hunter.HunterMock
	}
	type input struct {
		request  func() *http.Request
		response *httptest.ResponseRecorder
	}
	type output struct {
		code    int
		body    string
		headers http.Header
	}
	type testCase struct {
		name    string
		arrange arrange
		input   input
		output  output
	}

	testCases := []testCase{
		// case 1: success to hunt the prey
		{
			name: "case 1: success to hunt the prey",
			arrange: arrange{
				mockHunter: func() *hunter.HunterMock {
					return hunter.NewHunterMock()
				},
			},
			input: input{
				request: func() *http.Request {
					return httptest.NewRequest(http.MethodPost, "/hunt", nil)
				},
				response: httptest.NewRecorder(),
			},
			output: output{
				code:    http.StatusOK,
				body:    `{"message":"prey hunted","data":null}`,
				headers: http.Header{"Content-Type": []string{"application/json"}},
			},
		},
		// case 2: hunter can not hunt the prey
		{
			name: "case 2: hunter can not hunt the prey",
			arrange: arrange{
				mockHunter: func() *hunter.HunterMock {
					mk := hunter.NewHunterMock()
					mk.HuntFunc = func(pr prey.Prey) (float64, error) {
						return 0.0, hunter.ErrCanNotHunt
					}
					return mk
				},
			},
			input: input{
				request: func() *http.Request {
					return httptest.NewRequest(http.MethodPost, "/hunt", nil)
				},
				response: httptest.NewRecorder(),
			},
			output: output{
				code: http.StatusInternalServerError,
				body: fmt.Sprintf(
					`{"status":"%s","message":"%s"}`,
					http.StatusText(http.StatusInternalServerError),
					"can not hunt the prey",
				),
				headers: http.Header{"Content-Type": []string{"application/json"}},
			},
		},
		// case 3: internal server error
		{
			name: "case 3: internal server error",
			arrange: arrange{
				mockHunter: func() *hunter.HunterMock {
					mk := hunter.NewHunterMock()
					mk.HuntFunc = func(pr prey.Prey) (float64, error) {
						return 0.0, errors.New("internal error")
					}
					return mk
				},
			},
			input: input{
				request: func() *http.Request {
					return httptest.NewRequest(http.MethodPost, "/hunt", nil)
				},
				response: httptest.NewRecorder(),
			},
			output: output{
				code: http.StatusInternalServerError,
				body: fmt.Sprintf(
					`{"status":"%s","message":"%s"}`,
					http.StatusText(http.StatusInternalServerError),
					"internal server error",
				),
				headers: http.Header{"Content-Type": []string{"application/json"}},
			},
		},
	}

	// run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			// - hunter: mock
			mockHunter := tc.arrange.mockHunter()
			// - handler
			hd := handler.NewHunter(mockHunter, nil)
			hdFunc := hd.Hunt()

			// act
			hdFunc(tc.input.response, tc.input.request())

			// assert
			require.Equal(t, tc.output.code, tc.input.response.Code)
			require.JSONEq(t, tc.output.body, tc.input.response.Body.String())
			require.Equal(t, tc.output.headers, tc.input.response.Header())
		})
	}
}
