package testutils

import (
	"os"
	"strconv"
	"sync"
	"testing"
)

const (
	IntegralTestKey = "INTEGRAL_TEST"
)

var (
	integralTestEnabled bool = false
	once                sync.Once
)

// SkipIntegral
// skips slow integral tests if env variable is not enabled
func SkipIntegral(t testing.TB) {
	t.Helper()

	if IntegralEnabled() {
		return
	}

	t.Skipf("skipped integral test %s=0", IntegralTestKey)
}

func IntegralEnabled() bool {
	once.Do(func() {
		val, ok := os.LookupEnv(IntegralTestKey)

		if !ok {
			return
		}

		if v, err := strconv.ParseInt(val, 10, 64); err == nil {
			integralTestEnabled = v != 0
		}
	})

	return integralTestEnabled
}
