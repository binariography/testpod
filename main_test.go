package main

import (
	"testing"

	"github.com/pyroscope-io/client/pyroscope"
)

func TestPyroscopeConfiguration(t *testing.T) {
	// We test if the config struct accepts our constants
	// This ensures that if we update the SDK and constants change, the CI fails.
	cfg := pyroscope.Config{
		ApplicationName: "testpod.test",
		ServerAddress:   "http://localhost:4040",
		ProfileTypes: []pyroscope.ProfileType{
			pyroscope.ProfileCPU,
			pyroscope.ProfileInuseSpace,
		},
	}

	if cfg.ApplicationName != "testpod.test" {
		t.Errorf("Expected application name testpod.test, got %s", cfg.ApplicationName)
	}
}
