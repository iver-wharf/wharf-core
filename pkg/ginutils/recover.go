package ginutils

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iver-wharf/wharf-core/pkg/problem"
)

func RecoverHandle(c *gin.Context, err interface{}) {
	WriteProblem(c, problem.Response{
		Type:   "/prob/api/internal-server-error",
		Title:  "Internal server error.",
		Status: http.StatusInternalServerError,
		Detail: fmt.Sprintf("Unhandled error: %s", err),
	})
}
