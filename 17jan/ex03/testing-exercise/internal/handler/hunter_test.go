package handler_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testdoubles/internal/handler"
	"testdoubles/internal/hunter"
	"testdoubles/internal/prey"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	url_base = "/hunter"
)

type TestCases struct {
	description  string
	input        string
	expectedCode int
	expectedBody string
}

type TestCasesMock struct {
	TestCases
	expectedMockCalls int
}

func TestHunterHandler_ConfigurePrey(t *testing.T) {
	type testCasesConfigurePrey struct {
		TestCases
		stubPrey func() *prey.PreyStub
	}

	testCases := []testCasesConfigurePrey{
		{
			TestCases: TestCases{
				description: "case 1: success - configured prey",
				input: `{
					"speed": 100.0,
					"position": {
						"X": 10.0,
						"Y": 20.0,
						"Z": 30.0
					}
				}`,
				expectedCode: http.StatusCreated,
				expectedBody: `{ "message": "configured prey" }`,
			},
			stubPrey: func() *prey.PreyStub {
				return prey.NewPreyStub()
			},
		},
		{
			TestCases: TestCases{
				description: "case 2: error - invalid prey body",
				input: `{
					"speed": invalid,
					"position": {
						"X": 10.0,
						"Y": 20.0,
						"Z": 30.0
					}
				}`,
				expectedCode: http.StatusBadRequest,
				expectedBody: `{ "message": "invalid body", "status": "Bad Request" }`,
			},
			stubPrey: func() *prey.PreyStub {
				return prey.NewPreyStub()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			// handler and its dependencies
			stubPrey := tc.stubPrey()
			hd := handler.NewHunter(nil, stubPrey)
			hdFunc := hd.ConfigurePrey()

			// http request and response
			request := httptest.NewRequest("POST", url_base+"/configure-prey", io.NopCloser(strings.NewReader(tc.input)))
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()

			// WHEN
			hdFunc(response, request)

			// THEN
			require.Equal(t, tc.expectedCode, response.Code)
			require.JSONEq(t, tc.expectedBody, response.Body.String())
		})
	}
}

func TestHunterHandler_ConfigureHunter(t *testing.T) {
	type testCasesConfigureHunter struct {
		TestCasesMock
		mock func() *hunter.HunterMock
	}

	testCases := []testCasesConfigureHunter{
		{
			TestCasesMock: TestCasesMock{
				TestCases: TestCases{
					description: "case 1: success - configured hunter",
					input: `{
						"speed": 100.0,
						"position": {
							"X": 10.0,
							"Y": 20.0,
							"Z": 30.0
						}
					}`,
					expectedCode: http.StatusCreated,
					expectedBody: `{
						"message": "configured hunter"
					}`,
				},
				expectedMockCalls: 1,
			},
			mock: func() *hunter.HunterMock {
				return hunter.NewHunterMock()
			},
		},
		{
			TestCasesMock: TestCasesMock{
				TestCases: TestCases{
					description: "case 2: error - invalid hunter body",
					input: `{
						"speed": invalid,
						"position": {
							"X": 10.0,
							"Y": 20.0,
							"Z": 30.0
						}
					}`,
					expectedCode: http.StatusBadRequest,
					expectedBody: `{
						"message": "invalid body",
						"status": "Bad Request"
					}`,
				},
				expectedMockCalls: 0,
			},
			mock: func() *hunter.HunterMock {
				return hunter.NewHunterMock()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			// handler and its dependencies
			mockHunter := hunter.NewHunterMock()
			hd := handler.NewHunter(mockHunter, nil)
			hdFunc := hd.ConfigureHunter()

			// http request and response
			request := httptest.NewRequest("POST", url_base+"/configure-hunter", io.NopCloser(strings.NewReader(tc.input)))
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()

			// WHEN
			hdFunc(response, request)

			// THEN
			require.Equal(t, tc.expectedCode, response.Code)
			require.Equal(t, tc.expectedMockCalls, mockHunter.Calls.Configure)
			require.JSONEq(t, tc.expectedBody, response.Body.String())
		})
	}
}

func TestHunterHandler_Hunt(t *testing.T) {
	type testCasesHunt struct {
		TestCasesMock
		mock func() *hunter.HunterMock
	}
	testCases := []testCasesHunt{
		{
			TestCasesMock: TestCasesMock{
				TestCases: TestCases{
					description:  "case 1: success - hunter hunt prey",
					input:        "",
					expectedCode: http.StatusOK,
					expectedBody: `{
						"message": "hunt completed",
						"status":  "hunter hunted the prey in 10.000000 seconds"
					}`,
				},
				expectedMockCalls: 1,
			},
			mock: func() *hunter.HunterMock {
				mockHunter := hunter.NewHunterMock()
				mockHunter.HuntFunc = func(pr prey.Prey) (duration float64, err error) {
					return 10.0, nil
				}
				return mockHunter
			},
		},
		{
			TestCasesMock: TestCasesMock{
				TestCases: TestCases{
					description:  "case 2: success - hunter cannot hunt prey",
					input:        "",
					expectedCode: http.StatusOK,
					expectedBody: `{
						"message": "hunt completed",
						"status":  "hunter can not hunt the prey after 20.000000 seconds"
					}`,
				},
				expectedMockCalls: 1,
			},
			mock: func() *hunter.HunterMock {
				mockHunter := hunter.NewHunterMock()
				mockHunter.HuntFunc = func(pr prey.Prey) (duration float64, err error) {
					return 20.0, hunter.ErrCanNotHunt
				}
				return mockHunter
			},
		},
		{
			TestCasesMock: TestCasesMock{
				TestCases: TestCases{
					description:  "case 3: error - internal server error",
					input:        "",
					expectedCode: http.StatusInternalServerError,
					expectedBody: `{
						"message": "error hunting",
						"status": "Internal Server Error"
					}`,
				},
				expectedMockCalls: 1,
			},
			mock: func() *hunter.HunterMock {
				mockHunter := hunter.NewHunterMock()
				mockHunter.HuntFunc = func(pr prey.Prey) (duration float64, err error) {
					return 0.0, errors.New("internal server error")
				}
				return mockHunter
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			// handler and its dependencies
			mockHunter := tc.mock()

			hd := handler.NewHunter(mockHunter, nil)
			hdFunc := hd.Hunt()

			// http request and response
			request := httptest.NewRequest("POST", url_base+"/hunt", nil)
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()

			// WHEN
			hdFunc(response, request)

			// THEN
			require.Equal(t, tc.expectedCode, response.Code)
			require.Equal(t, tc.expectedMockCalls, mockHunter.Calls.Hunt)
			require.JSONEq(t, tc.expectedBody, response.Body.String())
		})
	}
}
