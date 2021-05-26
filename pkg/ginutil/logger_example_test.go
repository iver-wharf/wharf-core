package ginutil_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/iver-wharf/wharf-core/pkg/ginutil"
	"github.com/iver-wharf/wharf-core/pkg/logger"
	"github.com/iver-wharf/wharf-core/pkg/logger/consolepretty"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func ExampleLoggerWithConfig() {
	r := gin.New()

	r.Use(ginutil.LoggerWithConfig(ginutil.LoggerConfig{
		Logger:      logger.NewScoped("gin"),
		OmitLatency: true, // untestable
	}))

	// Don't forget to set up your outputs!
	conf := consolepretty.Config{
		DisableDate:       true,
		DisableCallerLine: true,
	}
	defer logger.ClearOutputs()
	logger.AddOutput(logger.LevelDebug, consolepretty.New(conf))

	// Faking a request here
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	r.ServeHTTP(w, req)

	// Output:
	// [DEBUG | gin | gin@v1.7.1/logger.go] clientIp=``  method=GET  path=/ping  status=404
}

func ExampleDefaultLoggerHandler() {
	r := gin.New()

	r.Use(ginutil.DefaultLoggerHandler)
	gin.DefaultWriter = ginutil.DefaultLoggerWriter
	gin.DefaultErrorWriter = ginutil.DefaultLoggerWriter

	// Don't forget to set up your outputs!
	defer logger.ClearOutputs()
	logger.AddOutput(logger.LevelDebug, consolepretty.Default)

	// Faking a request here
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	r.ServeHTTP(w, req)
}
