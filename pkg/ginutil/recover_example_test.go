package ginutil_test

import (
	"github.com/gin-gonic/gin"
	"github.com/iver-wharf/wharf-core/v2/pkg/ginutil"
)

func ExampleRecoverProblem() {
	r := gin.New()

	r.Use(
		ginutil.RecoverProblem,
	)
}

func ExampleRecoverProblemHandle() {
	r := gin.New()

	r.Use(
		gin.CustomRecovery(ginutil.RecoverProblemHandle),
	)
}
