package ginutil_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/iver-wharf/wharf-core/v2/pkg/ginutil"
	"github.com/iver-wharf/wharf-core/v2/pkg/logger"
	"github.com/iver-wharf/wharf-core/v2/pkg/logger/consolepretty"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func ExampleDefaultLoggerHandler() {
	logger.AddOutput(logger.LevelDebug, consolepretty.Default)

	r := gin.New()

	r.Use(ginutil.DefaultLoggerHandler)
	gin.DefaultWriter = ginutil.DefaultLoggerWriter
	gin.DefaultErrorWriter = ginutil.DefaultLoggerWriter

	// Faking a request here
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	r.ServeHTTP(w, req)
}
