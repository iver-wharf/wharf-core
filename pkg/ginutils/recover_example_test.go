package ginutils_test

import (
	"github.com/gin-gonic/gin"
	"github.com/iver-wharf/wharf-core/pkg/ginutils"
)

func ExampleRecoverProblemHandle() {
	r := gin.New()

	r.Use(
		gin.CustomRecovery(ginutils.RecoverProblemHandle),
	)
}
