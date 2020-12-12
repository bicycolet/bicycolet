package throttle_test

import (
	"os"
	"time"

	"github.com/bicycolet/bicycolet/internal/logger/throttle"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

func ExampleNewFilter() {
	logger := log.NewLogfmtLogger(os.Stdout)
	logger = throttle.NewFilter(logger, "component", 10, time.Second)

	for i := 0; i < 100; i++ {
		level.Debug(logger).Log("component", i)
	}
	// Output:
	// level=debug component=0
	// level=debug component=1
	// level=debug component=2
	// level=debug component=3
	// level=debug component=4
	// level=debug component=5
	// level=debug component=6
	// level=debug component=7
	// level=debug component=8
	// level=debug component=9
}
