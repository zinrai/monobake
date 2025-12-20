package main

import (
	"fmt"
	"strings"

	"github.com/Masterminds/semver/v3"
)

// TagInfo holds parsed Git tag information.
type TagInfo struct {
	Target  string
	Version string
}

// ParseGitRef parses a Git ref (e.g., refs/tags/backend/v1.0.0) into TagInfo.
func ParseGitRef(ref string) (*TagInfo, error) {
	tag := strings.TrimPrefix(ref, "refs/tags/")

	lastSlash := strings.LastIndex(tag, "/")
	if lastSlash == -1 {
		return nil, fmt.Errorf("invalid tag format: missing separator '/'")
	}

	target := tag[:lastSlash]
	version := tag[lastSlash+1:]

	if target == "" {
		return nil, fmt.Errorf("invalid tag format: empty target")
	}

	if version == "" {
		return nil, fmt.Errorf("invalid tag format: empty version")
	}

	v, err := semver.NewVersion(version)
	if err != nil {
		return nil, fmt.Errorf("invalid version: %w", err)
	}

	return &TagInfo{
		Target:  target,
		Version: v.String(),
	}, nil
}
