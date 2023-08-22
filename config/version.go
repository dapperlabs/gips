package config

import (
	"time"

	"github.com/carlmjohnson/versioninfo"
)

type VersionInfo struct {
	Revision   string    `json:"revision"`
	DirtyBuild bool      `json:"dirty_build"`
	LastCommit time.Time `json:"last_commit"`
}

func GetVersionInfo() VersionInfo {
	v := VersionInfo{
		Revision:   versioninfo.Revision,
		DirtyBuild: versioninfo.DirtyBuild,
		LastCommit: versioninfo.LastCommit,
	}
	return v
}
