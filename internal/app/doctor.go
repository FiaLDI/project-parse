package app

import (
	"context"
	"fmt"
	"os"
	"runtime"

	"github.com/FiaLDI/project-parse/internal/version"
)

// DoctorCheck is a single health check result.
type DoctorCheck struct {
	Name     string `json:"name"`
	OK       bool   `json:"ok"`
	Critical bool   `json:"critical"`
	Message  string `json:"message"`
}

// DoctorResult aggregates environment diagnostics.
type DoctorResult struct {
	Version string        `json:"version"`
	Go      string        `json:"go"`
	Checks  []DoctorCheck `json:"checks"`
}

// Doctor inspects the local environment and configuration.
func (a *App) Doctor(ctx context.Context) (DoctorResult, error) {
	_ = ctx
	res := DoctorResult{
		Version: version.String(),
		Go:      runtime.Version(),
		Checks:  make([]DoctorCheck, 0, 6),
	}

	cfgErr := a.cfg.Validate()
	res.Checks = append(res.Checks, DoctorCheck{
		Name:     "config",
		OK:       cfgErr == nil,
		Critical: true,
		Message:  checkMsg(cfgErr == nil, "configuration is valid", fmt.Sprintf("%v", cfgErr)),
	})

	res.Checks = append(res.Checks, DoctorCheck{
		Name:     "logger",
		OK:       a.log != nil,
		Critical: true,
		Message:  checkMsg(a.log != nil, "logger configured", "logger is nil"),
	})

	cwd, err := os.Getwd()
	res.Checks = append(res.Checks, DoctorCheck{
		Name:     "workdir",
		OK:       err == nil,
		Critical: true,
		Message:  checkMsg(err == nil, cwd, fmt.Sprintf("getwd: %v", err)),
	})

	outDir := a.cfg.Report.OutputDir
	if outDir == "" {
		outDir = "./.project-parser"
	}
	writable := dirWritable(outDir)
	res.Checks = append(res.Checks, DoctorCheck{
		Name:     "output_dir",
		OK:       writable,
		Critical: true,
		Message:  checkMsg(writable, outDir+" is writable (or creatable)", outDir+" is not writable"),
	})

	res.Checks = append(res.Checks, DoctorCheck{
		Name:     "scanner",
		OK:       a.scanner != nil,
		Critical: false,
		Message:  checkMsg(a.scanner != nil, "scanner injected", "scanner not wired yet (pending stage)"),
	})

	res.Checks = append(res.Checks, DoctorCheck{
		Name:     "registry",
		OK:       a.registry != nil,
		Critical: false,
		Message:  checkMsg(a.registry != nil, "plugin registry injected", "registry not wired yet (pending stage)"),
	})

	return res, nil
}

// HasCriticalFailures reports whether any critical doctor check failed.
func (r DoctorResult) HasCriticalFailures() bool {
	for _, c := range r.Checks {
		if c.Critical && !c.OK {
			return true
		}
	}
	return false
}

func checkMsg(ok bool, good, bad string) string {
	if ok {
		return good
	}
	return bad
}

func dirWritable(path string) bool {
	if err := os.MkdirAll(path, 0o755); err != nil {
		return false
	}
	f, err := os.CreateTemp(path, ".project-parser-doctor-*")
	if err != nil {
		return false
	}
	name := f.Name()
	_ = f.Close()
	_ = os.Remove(name)
	return true
}
