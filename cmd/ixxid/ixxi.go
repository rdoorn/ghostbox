package main

import (
	"time"

	_ "net/http/pprof"

	"github.com/rdoorn/ixxi/internal/cmd"
	"github.com/rdoorn/ixxi/internal/handler"
	// Only enabled for profiling
)

// version is set during makefile
var version string
var versionBuild string
var versionSha string

// Initialize package
func init() {
	handler.Version = version
	handler.VersionBuild = versionBuild
	handler.VersionSha = versionSha
	handler.StartTime = time.Now()
}

// main start
func main() {
	cmd.Execute()
}
