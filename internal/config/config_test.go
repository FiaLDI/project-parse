package config

import "testing"

func TestDefaultValidate(t *testing.T) {
	cfg := Default()
	if err := cfg.Validate(); err != nil {
		t.Fatalf("default config invalid: %v", err)
	}
}

func TestValidateRejectsBadLogLevel(t *testing.T) {
	cfg := Default()
	cfg.Log.Level = "verbose"
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for bad log level")
	}
}

func TestEffectiveJobs(t *testing.T) {
	tests := []struct {
		name string
		jobs int
		want int
	}{
		{name: "explicit", jobs: 4, want: 4},
		{name: "zero_uses_cpu", jobs: 0, want: 0}, // checked dynamically
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Default()
			cfg.Scan.Jobs = tt.jobs
			got := cfg.EffectiveJobs()
			if tt.jobs > 0 {
				if got != tt.want {
					t.Fatalf("got %d want %d", got, tt.want)
				}
				return
			}
			if got < 1 {
				t.Fatalf("expected at least 1 job, got %d", got)
			}
		})
	}
}

func TestLoadMissingFileReturnsDefault(t *testing.T) {
	cfg, err := Load("/tmp/project-parser-does-not-exist-xyz.yaml")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Log.Level != "info" {
		t.Fatalf("expected default log level info, got %s", cfg.Log.Level)
	}
}
