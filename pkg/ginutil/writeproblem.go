package ginutil

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/iver-wharf/wharf-core/pkg/problem"
)

// WriteProblem writes the Problem as JSON into the output response body
// together with appropriate Content-Type header.
//
// Problem.Type is set to "about:blank" (as recommended by the IETF RFC-7808)
// if left unset, or converts scheme-less URIs to start with
// "https://iver-wharf.github.io/#/".
//
// Problem.Status is set to 500 (Internal Server Error) if left unset.
//
// Problem.Instance is set to the request URI from the gorm.Context if left
// unset.
//
// Problem.Title is set to "Unknown error." if left unset.
//
// Problem.Detail is unaltered.
//
// Problem.Errors is set to the errors set to gin.Context.Errors if left empty.
func WriteProblem(c *gin.Context, prob problem.Response) {
	if prob.Type == "" {
		prob.Type = "about:blank"
	} else if u, err := url.Parse(prob.Type); err == nil {
		prob.Type = problem.ConvertURLToAbsDocsURL(*u).String()
	}
	if prob.Status == 0 {
		prob.Status = http.StatusInternalServerError
	}
	if prob.Instance == "" && c.Request != nil {
		prob.Instance = c.Request.RequestURI
	}
	if prob.Title == "" {
		prob.Title = "Unknown error."
	}
	if len(prob.Errors) == 0 && len(c.Errors) > 0 {
		prob.Errors = c.Errors.Errors()
	}
	c.Header("Content-Type", problem.HTTPContentType)
	c.JSON(prob.Status, prob)
}

// WriteProblemError is a shorthand for adding an error via gin.Context.Error
// and writing the problem using WriteProblem.
func WriteProblemError(c *gin.Context, err error, prob problem.Response) {
	c.Error(err)
	WriteProblem(c, prob)
}

// WriteBodyReadError uses WriteProblemError to write a 400 "Bad Request"
// response with the type "/prob/api/unexpected-body-read-error".
//
// Meant to be used on unexpected error when reading the raw HTTP request body.
func WriteBodyReadError(c *gin.Context, err error, detail string) {
	WriteProblemError(c, err, problem.Response{
		Type:   "/prob/api/unexpected-body-read-error",
		Title:  "Error reading request body.",
		Status: http.StatusBadRequest,
		Detail: detail,
	})
}

// WriteMultipartFormReadError uses WriteProblemError to write a 400
// "Bad Request" response with the type
// "/prob/api/unexpected-multipart-read-error".
//
// Meant to be used on unexpected error when reading a multipart/form-data
// request using gin.Context.MultipartForm().
func WriteMultipartFormReadError(c *gin.Context, err error, detail string) {
	WriteProblemError(c, err, problem.Response{
		Type:   "/prob/api/unexpected-multipart-read-error",
		Title:  "Error reading multipart data.",
		Status: http.StatusBadRequest,
		Detail: detail,
	})
}

// WriteDBReadError uses WriteProblemError to write a 502 "Bad Gateway" response
// with the type "/prob/api/unexpected-db-read-error".
//
// Meant to be used on unexpected error responses when doing a SELECT or other
// read operation towards the database.
func WriteDBReadError(c *gin.Context, err error, detail string) {
	WriteProblemError(c, err, problem.Response{
		Type:   "/prob/api/unexpected-db-read-error",
		Title:  "Error reading from database.",
		Status: http.StatusBadGateway,
		Detail: detail,
	})
}

// WriteDBWriteError uses WriteProblemError to write a 502 "Bad Gateway"
// response with the type "/prob/api/unexpected-db-write-error".
//
// Meant to be used on unexpected error responses when doing a CREATE, UPDATE or
// other write operation towards the database.
func WriteDBWriteError(c *gin.Context, err error, detail string) {
	WriteProblemError(c, err, problem.Response{
		Type:   "/prob/api/unexpected-db-read-error",
		Title:  "Error writing to database.",
		Status: http.StatusBadGateway,
		Detail: detail,
	})
}

// WriteDBNotFound uses WriteProblem to write a 404 "Not Found" response with
// the type "/prob/api/record-not-found".
//
// Meant to be used when fetching a specific item from the database but it was
// not found so this response is returned instead.
func WriteDBNotFound(c *gin.Context, detail string) {
	WriteProblem(c, problem.Response{
		Type:   "/prob/api/record-not-found",
		Title:  "Record not found.",
		Status: http.StatusBadGateway,
		Detail: detail,
	})
}

// WriteInvalidParamError uses WriteProblemError to write a 400 "Bad Request"
// response with the type "/prob/api/invalid-param".
//
// Meant to be used when parsing parameters in an endpoint handler.
func WriteInvalidParamError(c *gin.Context, err error, paramName, detail string) {
	WriteProblemError(c, err, problem.Response{
		Type:     "/prob/api/invalid-param",
		Title:    "Invalid API parameter.",
		Detail:   detail,
		Status:   http.StatusBadRequest,
		Instance: fmt.Sprintf("%s#%s", c.Request.RequestURI, paramName),
	})
}

// WriteInvalidBindError uses WriteProblemError to write a 400 "Bad Request"
// response with the type "/prob/api/invalid-param".
//
// Meant to be used when binding parameters in an endpoint handler.
func WriteInvalidBindError(c *gin.Context, err error, detail string) {
	WriteProblemError(c, err, problem.Response{
		Type:     "/prob/api/invalid-param",
		Title:    "Invalid API parameter.",
		Detail:   detail,
		Status:   http.StatusBadRequest,
		Instance: c.Request.RequestURI,
	})
}

// WriteAPIClientReadError uses WriteProblemError to write a 502 "Bad Gateway"
// response with the type "/prob/api-client/unexpected-read-error".
//
// Meant to be used on unexpected error when reading data using the Wharf API.
func WriteAPIClientReadError(c *gin.Context, err error, detail string) {
	WriteProblemError(c, err, problem.Response{
		Type:   "/prob/api-client/unexpected-read-error",
		Title:  "Unexpected API client read error.",
		Status: http.StatusBadGateway,
		Detail: detail,
	})
}

// WriteAPIClientWriteError uses WriteProblemError to write a 502 "Bad Gateway"
// response with the type "/prob/api-client/unexpected-write-error".
//
// Meant to be used on unexpected error when writing data using the Wharf API.
func WriteAPIClientWriteError(c *gin.Context, err error, detail string) {
	WriteProblemError(c, err, problem.Response{
		Type:   "/prob/api-client/unexpected-write-error",
		Title:  "Unexpected API client write error.",
		Status: http.StatusBadGateway,
		Detail: detail,
	})
}

// WriteProviderResponseError uses WriteProblemError to write a
// 502 "Bad Gateway" response with the type
// "/prob/provider/unexpected-response-format".
//
// Meant to be used on unexpected error when a provider plugin fails to parse
// or interpret a response from the remote provider.
func WriteProviderResponseError(c *gin.Context, err error, detail string) {
	WriteProblemError(c, err, problem.Response{
		Type:   "/prob/provider/unexpected-response-format",
		Title:  "Unexpected provider response format.",
		Status: http.StatusBadGateway,
		Detail: detail,
	})
}

// WriteFetchBuildDefinitionError uses WriteProblemError to write a
// 502 "Bad Gateway" response with the type
// "/prob/provider/fetch-build-definition".
//
// Meant to be used on error when the provider plugin fails to fetch the
// build definition from the remote provider.
func WriteFetchBuildDefinitionError(c *gin.Context, err error, detail string) {
	WriteProblemError(c, err, problem.Response{
		Type:   "/prob/provider/fetch-build-definition",
		Title:  "Error fetching build definition.",
		Status: http.StatusBadGateway,
		Detail: detail,
	})
}

// WriteComposingProviderDataError uses WriteProblemError to write a
// 502 "Bad Gateway" response with the type "/prob/provider/composing-provider-data".
//
// Meant to be used by the provider plugins on error when composing the
// provider object to submit to the Wharf API, such as when it fails to parse
// URLs received from the remote provider.
func WriteComposingProviderDataError(c *gin.Context, err error, detail string) {
	WriteProblemError(c, err, problem.Response{
		Type:   "/prob/provider/composing-provider-data",
		Title:  "Error composing provider data.",
		Status: http.StatusBadGateway,
		Detail: detail,
	})
}

// WriteTriggerError uses WriteProblemError to write a 502 "Bad Gateway"
// response with the type "/prob/api-client/unexpected-trigger-error".
//
// Meant to be used when unexpectedly failing to trigger a new build indirectly
// from a Wharf API client, such as from a Wharf provider plugin.
func WriteTriggerError(c *gin.Context, err error, detail string) {
	WriteProblemError(c, err, problem.Response{
		Type:   "/prob/api-client/unexpected-trigger-error",
		Title:  "Unexpected trigger error.",
		Status: http.StatusBadGateway,
		Detail: detail,
	})
}

// WriteUnauthorizedError uses WriteProblemError to write a 401 "Unauthorized"
// response with the type "/prob/api/unauthorized".
//
// Meant to be used for failed authentication.
func WriteUnauthorizedError(c *gin.Context, err error, detail string) {
	WriteProblemError(c, err, problem.Response{
		Type:   "/prob/api/unauthorized",
		Title:  "Unauthorized.",
		Status: http.StatusUnauthorized,
		Detail: detail,
	})
}

// WriteUnauthorized uses WriteProblem to write a 401 "Unauthorized"
// response with the type "/prob/api/unauthorized".
//
// Meant to be used for failed authentication.
func WriteUnauthorized(c *gin.Context, detail string) {
	WriteProblem(c, problem.Response{
		Type:   "/prob/api/unauthorized",
		Title:  "Unauthorized.",
		Status: http.StatusUnauthorized,
		Detail: detail,
	})
}
