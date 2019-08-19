package main

import (
	"time"

	_ "net/http/pprof"

	"github.com/rdoorn/ghostbox/internal/cmd"
	"github.com/rdoorn/ghostbox/internal/ghostbox"
	// Only enabled for profiling
)

// version is set during makefile
var version string
var versionBuild string
var versionSha string

// Initialize package
func init() {
	ghostbox.Version = version
	ghostbox.VersionBuild = versionBuild
	ghostbox.VersionSha = versionSha
	ghostbox.StartTime = time.Now()
}

// main start
func main() {
	cmd.Execute()
}
