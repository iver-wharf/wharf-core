package gormutil_test

import (
	"fmt"

	"github.com/iver-wharf/wharf-core/pkg/gormutil"
	"github.com/iver-wharf/wharf-core/pkg/logger"
	"github.com/iver-wharf/wharf-core/pkg/logger/consolepretty"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ExampleDefaultLogger() {
	logger.AddOutput(logger.LevelDebug, consolepretty.Default)

	db, err := gorm.Open(postgres.Open("host=localhost"), &gorm.Config{
		DryRun:               true,
		DisableAutomaticPing: true,
		Logger:               gormutil.DefaultLogger,
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
}
