package clamav

import (
	"strings"

	"github.com/DavidHoenisch/go-sysinfo/internal/probe"
)

type Config struct {
	LocalSocket         string
	DatabaseDirectory   string
	FreshClamDatabaseDir string
}

type DatabaseStats struct {
	SignatureCount string
	Raw            string
}

type Info struct {
	Version        string
	SignatureCount string
	DatabasePath   string
	Running        bool
}

func localSocketPath(r Reader) string {
	content := strings.TrimSpace(r.fs().Read(defaultClamdConf))
	if content == "" {
		return defaultLocalSocket
	}
	values := probe.ParseKeyValueConfig(content)
	if path := values["LocalSocket"]; path != "" {
		return path
	}
	return defaultLocalSocket
}

func databaseDirectory(r Reader) string {
	content := strings.TrimSpace(r.fs().Read(defaultClamdConf))
	if content != "" {
		values := probe.ParseKeyValueConfig(content)
		if path := values["DatabaseDirectory"]; path != "" {
			return path
		}
	}

	content = strings.TrimSpace(r.fs().Read(defaultFreshclamConf))
	if content == "" {
		return ""
	}
	values := probe.ParseKeyValueConfig(content)
	return values["DatabaseDirectory"]
}

func Available(r Reader) bool {
	if pingClamd(r.socket(), localSocketPath(r)) {
		return true
	}
	_, err := r.cmd().Run("clamd", "--version")
	return err == nil
}

func (r Reader) Available() bool {
	return Available(r)
}

func GetVersion(r Reader) string {
	if response, ok := clamdRoundTrip(r.socket(), localSocketPath(r), "VERSION"); ok {
		return response
	}

	out, err := r.cmd().Run("clamd", "--version")
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func (r Reader) GetVersion() string {
	return GetVersion(r)
}

func GetDatabaseStats(r Reader) *DatabaseStats {
	response, ok := clamdRoundTrip(r.socket(), localSocketPath(r), "STATS")
	if !ok {
		return &DatabaseStats{}
	}
	return &DatabaseStats{
		SignatureCount: parseStatsSignatures(response),
		Raw:            response,
	}
}

func (r Reader) GetDatabaseStats() *DatabaseStats {
	return GetDatabaseStats(r)
}

func GetConfig(r Reader) *Config {
	clamdContent := strings.TrimSpace(r.fs().Read(defaultClamdConf))
	freshclamContent := strings.TrimSpace(r.fs().Read(defaultFreshclamConf))

	clamdValues := probe.ParseKeyValueConfig(clamdContent)
	freshclamValues := probe.ParseKeyValueConfig(freshclamContent)

	localSocket := clamdValues["LocalSocket"]
	if localSocket == "" {
		localSocket = defaultLocalSocket
	}

	databaseDirectory := clamdValues["DatabaseDirectory"]
	if databaseDirectory == "" {
		databaseDirectory = freshclamValues["DatabaseDirectory"]
	}

	return &Config{
		LocalSocket:          localSocket,
		DatabaseDirectory:    databaseDirectory,
		FreshClamDatabaseDir: freshclamValues["DatabaseDirectory"],
	}
}

func (r Reader) GetConfig() *Config {
	return GetConfig(r)
}

func GetInfo(r Reader) *Info {
	running := pingClamd(r.socket(), localSocketPath(r))
	return &Info{
		Version:        GetVersion(r),
		SignatureCount: GetDatabaseStats(r).SignatureCount,
		DatabasePath:   databaseDirectory(r),
		Running:        running,
	}
}

func (r Reader) GetInfo() *Info {
	return GetInfo(r)
}
