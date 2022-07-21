// Package verflag defines utility functions to handle command line flags
// related to version.
package verflag

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/pflag"

	"github.com/yimi-go/version"
)

type versionValue int

// Define some consts.
const (
	VersionFalse versionValue = iota
	VersionTrue
	VersionRaw
)

const strRawVersion string = "raw"

func (v *versionValue) IsBoolFlag() bool {
	return true
}

func (v *versionValue) Get() interface{} {
	return v
}

func (v *versionValue) Set(s string) error {
	if s == strRawVersion {
		*v = VersionRaw
		return nil
	}
	boolVal, err := strconv.ParseBool(s)
	if boolVal {
		*v = VersionTrue
	} else {
		*v = VersionFalse
	}
	return err
}

func (v *versionValue) String() string {
	if v == nil {
		return ""
	}
	if *v == VersionRaw {
		return strRawVersion
	}
	return fmt.Sprintf("%v", *v == VersionTrue)
}

// Type is the type of the flag as required by the pflag.Value interface.
func (v *versionValue) Type() string {
	return "version"
}

// VersionVar defines a flag with the specified name and usage string.
func VersionVar(p *versionValue, name string, value versionValue, usage string) {
	*p = value
	pflag.Var(p, name, usage)
	// "--version" will be treated as "--version=true"
	pflag.Lookup(name).NoOptDefVal = "true"
}

// Version wraps the VersionVar function.
func Version(name string, value versionValue, usage string) *versionValue {
	p := new(versionValue)
	VersionVar(p, name, value, usage)
	return p
}

const versionFlagName = "version"

var versionFlag = Version(versionFlagName, VersionFalse, "Print version information and quit.")

// AddFlags registers this package's flags on arbitrary FlagSets, such that they point to the
// same value as the global flags.
func AddFlags(fs *pflag.FlagSet) {
	fs.AddFlag(pflag.Lookup(versionFlagName))
}

// PrintAndExitIfRequested will check if the -version flag was passed
// and, if so, print the version and exit.
func PrintAndExitIfRequested() {
	if *versionFlag == VersionRaw {
		fmt.Printf("%#v\n", version.Get())
		os.Exit(0)
	} else if *versionFlag == VersionTrue {
		fmt.Printf("%s\n", version.Get())
		os.Exit(0)
	}
}
