// Package version supplies version information collected at build time.
package version

import (
	"bytes"
	"encoding/json"
	"fmt"
	"runtime"
	"text/template"
)

var (
	// GitVersion is semantic version.
	GitVersion = "v0.0.0-master+$Format:%h$"
	// BuildDate in ISO8601 format, output of $(date -u +'%Y-%m-%dT%H:%M:%SZ').
	BuildDate = "1970-01-01T00:00:00Z"
	// GitCommit is sha1 for git, output of $(git rev-parse HEAD).
	GitCommit = "$Format:%H$"
	// GitTreeState is state of git tree, either "clean" or "dirty".
	GitTreeState = ""
)

// Info contains versioning information.
type Info struct {
	GitVersion   string `json:"gitVersion"`
	GitCommit    string `json:"gitCommit"`
	GitTreeState string `json:"gitTreeState"`
	BuildDate    string `json:"buildDate"`
	GoVersion    string `json:"goVersion"`
	Compiler     string `json:"compiler"`
	Platform     string `json:"platform"`
}

// language=GoTemplate
const tableTemplate = `  gitVersion: {{.GitVersion}}
   gitCommit: {{.GitCommit}}
gitTreeState: {{.GitTreeState}}
   buildDate: {{.BuildDate}}
   goVersion: {{.GoVersion}}
    compiler: {{.Compiler}}
    platform: {{.Platform}}`

// String returns info as a human-friendly version string.
func (info Info) String() string {
	buf := &bytes.Buffer{}
	// I'm sure it would not return an error.
	tmpl, _ := template.New("version").Parse(tableTemplate)
	_ = tmpl.Execute(buf, info)
	return buf.String()
}

// ToJSON returns the JSON string of version information.
func (info Info) ToJSON() string {
	s, _ := json.Marshal(info)

	return string(s)
}

// Get returns the overall codebase version. It's for detecting
// what code a binary was built from.
func Get() Info {
	return Info{
		GitVersion:   GitVersion,
		GitCommit:    GitCommit,
		GitTreeState: GitTreeState,
		BuildDate:    BuildDate,
		GoVersion:    runtime.Version(),
		Compiler:     runtime.Compiler,
		Platform:     fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}
