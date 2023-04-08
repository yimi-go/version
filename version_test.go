package version

import (
	"encoding/json"
	"fmt"
	"runtime"
	"runtime/debug"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func genVersion(unknown, devel bool) *version {
	return &version{
		BuildInfoReader: func() (info *debug.BuildInfo, ok bool) {
			if unknown {
				return nil, false
			}
			main := debug.Module{
				Path: "host/user/repo",
			}
			if !devel {
				main.Version = "v0.0.1"
				main.Sum = "some hash hex"
			} else {
				main.Version = "(devel)"
			}
			info = &debug.BuildInfo{
				GoVersion: runtime.Version(),
				Path:      "host/user/repo/main/pkg",
				Main:      main,
			}
			if devel {
				info.Settings = append(info.Settings, debug.BuildSetting{Key: "vcs", Value: "git"})
				info.Settings = append(info.Settings, debug.BuildSetting{Key: "vcs.revision", Value: "abc...xyz"})
				info.Settings = append(info.Settings, debug.BuildSetting{Key: "vcs.time", Value: "2022-04-03T16:59:50Z"})
				info.Settings = append(info.Settings, debug.BuildSetting{Key: "vcs.modified", Value: "false"})
			}
			return info, true
		},
	}
}

func TestGet(t *testing.T) {
	tests := []struct {
		name    string
		unknown bool
		devel   bool
	}{
		{
			name:    "unknown",
			unknown: true,
		},
		{
			name: "install",
		},
		{
			name:  "devel",
			devel: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dv := defaultVersion
			defer func() {
				defaultVersion = dv
			}()
			v := genVersion(tt.unknown, tt.devel)
			defaultVersion = v
			bi, ok := v.BuildInfoReader()
			require.NotEqual(t, tt.unknown, ok)
			info := Get()
			if !ok {
				assert.Equal(t, "unknown", info.Version)
				assert.Equal(t, "unknown", info.VCS)
				assert.Equal(t, "unknown", info.Revision)
				assert.Equal(t, "unknown", info.BuildTime)
				assert.Equal(t, "unknown", info.VCSModified)
			} else {
				assert.Equal(t, bi.Main.Version, info.Version)
				var vcs, revision, buildTime, modified string
				for i := range bi.Settings {
					switch bi.Settings[i].Key {
					case "vcs":
						vcs = bi.Settings[i].Value
					case "vcs.revision":
						revision = bi.Settings[i].Value
					case "vcs.time":
						buildTime = bi.Settings[i].Value
					case "vcs.modified":
						modified = bi.Settings[i].Value
					}
				}
				assert.Equal(t, vcs, info.VCS)
				assert.Equal(t, revision, info.Revision)
				assert.Equal(t, buildTime, info.BuildTime)
				assert.Equal(t, modified, info.VCSModified)
			}
			assert.Equal(t, runtime.Version(), info.GoVersion)
			assert.Equal(t, runtime.Compiler, info.Compiler)
			assert.Equal(t, fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH), info.Platform)
		})
	}
}

func TestInfo_ToJSON(t *testing.T) {
	tests := []struct {
		name    string
		unknown bool
		devel   bool
	}{
		{
			name:    "unknown",
			unknown: true,
		},
		{
			name: "install",
		},
		{
			name:  "devel",
			devel: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dv := defaultVersion
			defer func() {
				defaultVersion = dv
			}()
			v := genVersion(tt.unknown, tt.devel)
			defaultVersion = v
			info := Get()
			s := info.ToJSON()
			t.Logf("json: %s", s)
			u := Info{}
			err := json.Unmarshal([]byte(s), &u)
			require.Nilf(t, err, "ToJSON: %s", err)
			assert.Equal(t, info, u)
		})
	}
}

func TestInfo_String(t *testing.T) {
	tests := []struct {
		name    string
		unknown bool
		devel   bool
	}{
		{
			name:    "unknown",
			unknown: true,
		},
		{
			name: "install",
		},
		{
			name:  "devel",
			devel: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dv := defaultVersion
			defer func() {
				defaultVersion = dv
			}()
			v := genVersion(tt.unknown, tt.devel)
			defaultVersion = v
			info := Get()
			s := info.String()
			t.Logf("\n%s", s)
		})
	}
}
