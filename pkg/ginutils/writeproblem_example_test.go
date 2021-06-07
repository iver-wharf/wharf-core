package ginutils_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/iver-wharf/wharf-core/pkg/ginutils"
	"github.com/iver-wharf/wharf-core/pkg/problem"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func indentedBodyFromResponse(resp *http.Response) string {
	body, _ := io.ReadAll(resp.Body)
	var indentedBodyBuff bytes.Buffer
	json.Indent(&indentedBodyBuff, body, "", "  ")
	return indentedBodyBuff.String()
}

func ExampleWriteProblem() {
	var prob = problem.Response{
		Type:     "https://iver-wharf.github.io/#/prob/build/run/invalid-input",
		Title:    "Invalid input variable for build.",
		Status:   400,
		Detail:   "Build requires input variable 'myInput' to be of type 'string', but got 'int' instead.",
		Instance: "/projects/12345/builds/run/6789",
		Errors: []string{
			"strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	ginutils.WriteProblem(c, prob)

	resp := w.Result()

	fmt.Println("HTTP/1.1", resp.Status)
	fmt.Println("Content-Type:", resp.Header.Get("Content-Type"))
	fmt.Println()
	fmt.Println(indentedBodyFromResponse(resp))

	// Output:
	// HTTP/1.1 400 Bad Request
	// Content-Type: application/problem+json
	//
	// {
	//   "type": "https://iver-wharf.github.io/#/prob/build/run/invalid-input",
	//   "title": "Invalid input variable for build.",
	//   "status": 400,
	//   "detail": "Build requires input variable 'myInput' to be of type 'string', but got 'int' instead.",
	//   "instance": "/projects/12345/builds/run/6789",
	//   "errors": [
	//     "strconv.ParseUint: parsing \"-1\": invalid syntax"
	//   ]
	// }
}

func ExampleWriteProblem_emptyResponse() {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/example/request/uri?foo=bar", nil)

	ginutils.WriteProblem(c, problem.Response{})

	resp := w.Result()

	fmt.Println("HTTP/1.1", resp.Status)
	fmt.Println("Content-Type:", resp.Header.Get("Content-Type"))
	fmt.Println()
	fmt.Println(indentedBodyFromResponse(resp))

	// Output:
	// HTTP/1.1 500 Internal Server Error
	// Content-Type: application/problem+json
	//
	// {
	//   "type": "about:blank",
	//   "title": "Unknown error.",
	//   "status": 500,
	//   "detail": "",
	//   "instance": "/example/request/uri?foo=bar",
	//   "errors": null
	// }
}
