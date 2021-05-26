package gormutil_test

import (
	"fmt"

	"github.com/iver-wharf/wharf-core/pkg/gormutil"
	"github.com/iver-wharf/wharf-core/pkg/logger"
	"github.com/iver-wharf/wharf-core/pkg/logger/consolepretty"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ExampleNewLogger() {
	defer logger.ClearOutputs()
	logger.AddOutput(logger.LevelDebug, consolepretty.New(consolepretty.Config{
		DisableDate:       true,
		DisableCallerLine: true,
	}))

	db, err := gorm.Open(postgres.Open("host=localhost"), &gorm.Config{
		DryRun:               true,
		DisableAutomaticPing: true,
		Logger: gormutil.NewLogger(gormutil.LoggerConfig{
			Logger: logger.NewScoped("GORM"),
		}),
	})

	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}

	type User struct {
		ID   int
		Name string `gorm:"size:256"`
	}

	db.Find(&User{}, 1)

	// Sample output:
	// [DEBUG | GORM | gorm@v1.21.10/callbacks.go] rows=0  elapsed=89.768µs  sql=`SELECT * FROM "users" WHERE "users"."id" = 1`
}
