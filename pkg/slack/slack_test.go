package slack

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	gifMocks "github.com/king-jam/gotd/pkg/gif/mocks"
	"github.com/stretchr/testify/mock"
)

type inputValues struct {
	user                  string
	command               string
	gifServiceGood        bool
	verificationTokenGood bool
}

type expectedValues struct {
	statusCode int
	respBody   string
}

type testCase struct {
	inputs   inputValues
	expected expectedValues
}

func TestHTTPHandler(t *testing.T) {
	tests := map[string]testCase{
		"valid setup": {
			inputs:   inputValues{"U5T9HLMAN", "/gotd", true, true},
			expected: expectedValues{200, "Requested GIF\nwww.link.com\nGIF Successfully posted to GOTD"},
		},
		"bad token": {
			inputs:   inputValues{"U5T9HLMAN", "/gotd", true, false},
			expected: expectedValues{200, "Requested GIF\nwww.link.com\nunable to validate slack token"},
		},
		"bad gif service": {
			inputs:   inputValues{"U5T9HLMAN", "/gotd", false, true},
			expected: expectedValues{200, "Requested GIF\nwww.link.com\nerror setting"},
		},
		"invalid user": {
			inputs:   inputValues{"WRONG USER", "/gotd", true, true},
			expected: expectedValues{200, "Requested GIF\nwww.link.com\nuser not authorized"},
		},
		"bad command": {
			inputs:   inputValues{"U5T9HLMAN", "/wrongcommand", true, true},
			expected: expectedValues{200, "Requested GIF\nwww.link.com\ninvalid slash command sent"},
		},
	}
	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			handlertest(t, test)
		})
	}
}

func handlertest(t *testing.T, test testCase) {
	gifSvcMock := new(gifMocks.Service)

	if test.inputs.gifServiceGood {
		gifSvcMock.On("Set", mock.Anything, mock.Anything).Return(nil)
	} else {
		gifSvcMock.On("Set", mock.Anything, mock.Anything).Return(errors.New("some error"))
	}

	testToken := "mytesttoken"

	var verificationToken string
	if test.inputs.verificationTokenGood {
		verificationToken = testToken
	} else {
		verificationToken = "WRONG"
	}

	form := url.Values{}
	form.Add("token", testToken)
	form.Add("team_id", "DOESNTMATTER")
	form.Add("team_domain", "DOESNTMATTER")
	form.Add("channel_id", "DOESNTMATTER")
	form.Add("channel_name", "DOESNTMATTER")
	form.Add("user_id", test.inputs.user)
	form.Add("user_name", "DOESNTMATTER")
	form.Add("command", test.inputs.command)
	form.Add("text", "www.link.com")

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("POST", "/slack", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := New(gifSvcMock, verificationToken)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != test.expected.statusCode {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	if rr.Body.String() != test.expected.respBody {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), test.expected.respBody)
	}
}
