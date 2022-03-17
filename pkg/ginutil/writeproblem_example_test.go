package ginutil_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/iver-wharf/wharf-core/v2/pkg/ginutil"
	"github.com/iver-wharf/wharf-core/v2/pkg/problem"
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

	ginutil.WriteProblem(c, prob)

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
