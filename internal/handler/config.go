package handler

import (
	"fmt"
	"sort"
	"time"

	"github.com/rdoorn/gohelper/filehelper"
	"github.com/rdoorn/gohelper/logging"
	"github.com/rdoorn/gohelper/tlsconfig"
	"github.com/rdoorn/ixxi/internal/ixdb"
)

const (
	// Name is the application name
	Name string = "ixxi"
	// DefaultPasswordSalt is the Password salt used if not provided
	DefaultPasswordSalt string = "iP3x.!fNBe;Ajf0="

	// KB is one kilo byte
	KB int64 = 1024
	// MB is one mega byte
	MB int64 = 1024 * 1024
	// GB is one giga byte
	GB int64 = 1024 * 1024 * 1024
)

var (
	// Version of application
	Version string
	// VersionBuild number
	VersionBuild string
	// VersionSha git commit of build
	VersionSha string
	// StartTime of application
	StartTime time.Time
	// ReloadTime last time a reload was successfull
	ReloadTime time.Time
	// FailedReloadTime last time a reload failed
	FailedReloadTime time.Time
	// FailedReloadError last time a reload failed
	FailedReloadError string

	// FreeTeerStorageLimitMB is the storage limit of newly registered non-paying users
	FreeTeerStorageLimitMB = 1 * GB
)

// Config holds your main config
type Config struct {
	LogLevel     string              `mapstructure:"log_level"`
	LogOutput    []string            `mapstructure:"log_output"`
	PidFile      string              `mapstructure:"pid_file"`
	Listener     string              `mapstructure:"listener"`
	Port         int                 `mapstructure:"port"`
	TLS          tlsconfig.TLSConfig `mapstructure:"tls"`
	Users        ixdb.DB             `mapstructure:"user_provider"`
	File         ixdb.DB             `mapstructure:"file_provider"`
	DB           ixdb.DB             `mapstructure:"db_provider"`
	PasswordSalt string              `mapstructure:"password_salt"`
}

func (c *Config) Verify() error {
	// check log level
	if _, err := logging.ToLevel(c.LogLevel); err != nil {
		return err
	}

	// check log output is writable
	sort.Strings(c.LogOutput) // order it here for testing later
	for _, output := range c.LogOutput {
		switch output {
		case "stderr", "stdout":
		default:
			if err := filehelper.IsWritable(output); err != nil {
				return fmt.Errorf("cannot write to log output %s: %s", output, err)
			}
			if filehelper.IsDir(output) {
				return fmt.Errorf("cannot write to log output %s: target is a directory not a file", output)
			}
		}
	}

	// check pid file is writable
	if err := filehelper.IsWritable(c.PidFile); err != nil {
		return fmt.Errorf("cannot write to log output %s: %s", c.PidFile, err)
	}
	if filehelper.IsDir(c.PidFile) {
		return fmt.Errorf("cannot write to log output %s: target is a directory not a file", c.PidFile)
	}

	if c.TLS.CertificateFile != "" {
		if err := c.TLS.Valid(); err != nil {
			return fmt.Errorf("tls configuration failed: %s", err)
		}
	}

	if err := c.Users.Verify(); err != nil {
		return fmt.Errorf("Invalid user store: %s", err)
	}
	if err := c.File.Verify(); err != nil {
		return fmt.Errorf("Invalid file store: %s", err)
	}
	if err := c.DB.Verify(); err != nil {
		return fmt.Errorf("Invalid database store: %s", err)
	}

	if c.PasswordSalt == "" {
		c.PasswordSalt = DefaultPasswordSalt
	}
	return nil
}

func (c *Config) Addr() string {
	return fmt.Sprintf("%s:%d", c.Listener, c.Port)
}
