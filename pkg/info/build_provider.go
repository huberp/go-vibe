package info

import "time"

// BuildInfoProvider provides build-time information about the application.
type BuildInfoProvider struct {
	Version   string
	CommitSHA string
	BuildTime string
	GoVersion string
}

// NewBuildInfoProvider creates a new BuildInfoProvider with the specified build information.
func NewBuildInfoProvider(version, commitSHA, buildTime, goVersion string) *BuildInfoProvider {
	// Use defaults if not provided
	if version == "" {
		version = "dev"
	}
	if commitSHA == "" {
		commitSHA = "unknown"
	}
	if buildTime == "" {
		buildTime = time.Now().Format(time.RFC3339)
	}
	if goVersion == "" {
		goVersion = "unknown"
	}

	return &BuildInfoProvider{
		Version:   version,
		CommitSHA: commitSHA,
		BuildTime: buildTime,
		GoVersion: goVersion,
	}
}

// Name returns the name of this provider.
func (b *BuildInfoProvider) Name() string {
	return "build"
}

// Info returns build information.
func (b *BuildInfoProvider) Info() (map[string]interface{}, error) {
	return map[string]interface{}{
		"version":    b.Version,
		"commit":     b.CommitSHA,
		"build_time": b.BuildTime,
		"go_version": b.GoVersion,
	}, nil
}
