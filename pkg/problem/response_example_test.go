package problem_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/iver-wharf/wharf-core/v2/pkg/problem"
)

func ExampleParseHTTPResponse() {
	req := httptest.NewRequest("GET", "/projects/123", nil)

	// Faking a HTTP response here
	req.Response = &http.Response{
		Body: io.NopCloser(strings.NewReader(`
{
  "type": "https://iver-wharf.github.io/#/prob/build/run/invalid-input",
  "title": "Invalid input variable for build.",
  "status": 400,
  "detail": "Build requires input variable 'myInput' to be of type 'string', but got 'int' instead.",
  "instance": "/projects/12345/builds/run/6789",
  "errors": [
    "strconv.ParseUint: parsing \"-1\": invalid syntax"
  ]
}
`)),
		Header: make(http.Header),
	}
	req.Response.Header.Add("Content-Type", problem.HTTPContentType)

	if problem.IsHTTPResponse(req.Response) {
		p, err := problem.ParseHTTPResponse(req.Response)
		if err != nil {
			panic(err)
		}

		fmt.Println(p.String())
	}

	// Output:
	// {(problem) HTTP 400, https://iver-wharf.github.io/#/prob/build/run/invalid-input
	//     Title: Invalid input variable for build.
	//    Detail: Build requires input variable 'myInput' to be of type 'string', but got 'int' instead.
	//  Error(s): [strconv.ParseUint: parsing "-1": invalid syntax]
	//  Instance: /projects/12345/builds/run/6789 }
}
