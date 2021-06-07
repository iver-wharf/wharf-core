package ginutils

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iver-wharf/wharf-core/pkg/problem"
)

// RecoverProblemHandle writes a HTTP "Internal Server Error" problem response.
// Meant to be used with the gin-gonic panic recover middleware.
func RecoverProblemHandle(c *gin.Context, err interface{}) {
	WriteProblem(c, problem.Response{
		Type:   "/prob/api/internal-server-error",
		Title:  "Internal server error.",
		Status: http.StatusInternalServerError,
		Detail: fmt.Sprintf("Unhandled error: %s", err),
	})
}
