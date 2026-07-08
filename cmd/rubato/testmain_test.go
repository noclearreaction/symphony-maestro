//go:build smoke

package main

import (
	"os"
	"testing"
)

// TestMain allows the test binary to double as a rubato server subprocess.
// When RUBATO_TEST_SUBPROCESS=1 is set, main() is called directly and the
// process acts as a rubato instance. Otherwise the test suite runs normally.
func TestMain(m *testing.M) {
	if os.Getenv("RUBATO_TEST_SUBPROCESS") == "1" {
		main()
		return
	}
	os.Exit(m.Run())
}
