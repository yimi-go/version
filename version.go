// Package version supplies version information collected at build time.
package version

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"runtime"
	"runtime/debug"
	"text/template"
)

// Info contains versioning information.
type Info struct {
	Version     string `json:"version,omitempty"`
	VCS         string `json:"vcs,omitempty"`
	Revision    string `json:"revision,omitempty"`
	BuildTime   string `json:"build_time,omitempty"`
	VCSModified string `json:"vcs_modified,omitempty"`
	GoVersion   string `json:"go_version"`
	Compiler    string `json:"compiler"`
	Platform    string `json:"platform"`
}

//go:embed table.gohtml
var tableTemplate string

// String returns info as a human-friendly version string.
func (info Info) String() string {
	buf := &bytes.Buffer{}
	// I'm sure it would not return an error.
	tmpl := template.Must(template.New("version").Parse(tableTemplate))
	_ = tmpl.Execute(buf, info)
	return buf.String()
}

// ToJSON returns the JSON string of version information.
func (info Info) ToJSON() string {
	s, _ := json.Marshal(info)

	return string(s)
}

type version struct {
	BuildInfoReader func() (info *debug.BuildInfo, ok bool)
}

// Get returns the overall codebase version. It's for detecting
// what code a binary was built from.
func (v *version) Get() Info {
	info, ok := v.BuildInfoReader()
	if !ok {
		return Info{
			Version:     "unknown",
			VCS:         "unknown",
			Revision:    "unknown",
			BuildTime:   "unknown",
			VCSModified: "unknown",
			GoVersion:   runtime.Version(),
			Compiler:    runtime.Compiler,
			Platform:    fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		}
	}
	var vcs, revision, buildTime, modified string
	for i := range info.Settings {
		switch info.Settings[i].Key {
		case "vcs":
			vcs = info.Settings[i].Value
		case "vcs.revision":
			revision = info.Settings[i].Value
		case "vcs.time":
			buildTime = info.Settings[i].Value
		case "vcs.modified":
			modified = info.Settings[i].Value
		}
	}
	return Info{
		Version:     info.Main.Version,
		VCS:         vcs,
		Revision:    revision,
		BuildTime:   buildTime,
		VCSModified: modified,
		GoVersion:   runtime.Version(),
		Compiler:    runtime.Compiler,
		Platform:    fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

var defaultVersion = &version{
	BuildInfoReader: debug.ReadBuildInfo,
}

func Get() Info {
	return defaultVersion.Get()
}
