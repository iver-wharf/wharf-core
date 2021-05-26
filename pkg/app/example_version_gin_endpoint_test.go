package app_test

import (
	_ "embed"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/iver-wharf/wharf-core/pkg/app"
)

// The version.yaml file should be populated by a CI pipeline build step just
// before building the binary for this application.
//
// For example, assuming you have BUILD_VERSION, BUILD_GIT_COMMIT, and BUILD_REF
// environment variables set before running the following script:
//
// 		#!/bin/sh
//
// 		cat <<EOF > version.yaml
// 		version: ${BUILD_VERSION}
// 		buildGitCommit: ${BUILD_GIT_COMMIT}
// 		buildDate: $(date '+%FT%T%:z')
// 		buildRef: ${BUILD_REF}
// 		EOF

// go:embed version.yaml
var versionFile []byte

// AppVersion is the type holding metadata about this application's version.
var AppVersion app.Version

// getVersionHandler godoc
// @summary Returns the version of this API
// @tags meta
// @success 200 {object} app.Version
// @router /version [get]
func getVersionHandler(c *gin.Context) {
	c.JSON(http.StatusOK, AppVersion)
}

func ExampleVersion_ginEndpoint() {
	if err := app.UnmarshalVersionYAML(versionFile, &AppVersion); err != nil {
		fmt.Println("Failed to read embedded version.yaml file:", err)
		os.Exit(1)
	}

	// If you use swaggo then you can set the API version like so:
	//
	// 		docs.SwaggerInfo.Version = AppVersion.Version
	//
	// More info: https://github.com/swaggo/swag#how-to-use-it-with-gin

	r := gin.Default()
	r.GET("/version", getVersionHandler)

	_ = r.Run()
}
