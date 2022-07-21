package version

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"runtime"
	"testing"
	"time"
)

func randVersionInfo() {
	hash := rand.Intn(1000000)
	GitVersion = fmt.Sprintf("v%d.%d.%d-%s+%x",
		rand.Intn(10), rand.Intn(10), rand.Intn(10),
		"master", hash)
	GitCommit = fmt.Sprintf("%x", hash)
	GitTreeState = "dirty"
	BuildDate = time.Now().UTC().Truncate(time.Second).Format("2006-01-02T15:04:05Z")
}

func TestGet(t *testing.T) {
	randVersionInfo()
	info := Get()
	if info.GitVersion != GitVersion {
		t.Errorf("GitVersion: expected %s, got %s", GitVersion, info.GitVersion)
	}
	if info.GitCommit != GitCommit {
		t.Errorf("GitCommit: expected %s, got %s", GitCommit, info.GitCommit)
	}
	if info.GitTreeState != GitTreeState {
		t.Errorf("GitTreeState: expected %s, got %s", GitTreeState, info.GitTreeState)
	}
	if info.BuildDate != BuildDate {
		t.Errorf("BuildDate: expected %s, got %s", BuildDate, info.BuildDate)
	}
	if info.GoVersion != runtime.Version() {
		t.Errorf("GoVersion: expected %s, got %s", runtime.Version(), info.GoVersion)
	}
	if info.Compiler != runtime.Compiler {
		t.Errorf("Compiler: expected %s, got %s", runtime.Compiler, info.Compiler)
	}
	if info.Platform != fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH) {
		t.Errorf("Platform: expected %s/%s, got %s", runtime.GOOS, runtime.GOARCH, info.Platform)
	}
}

func TestInfo_ToJSON(t *testing.T) {
	randVersionInfo()
	info := Get()
	s := info.ToJSON()
	v := Info{}
	err := json.Unmarshal([]byte(s), &v)
	if err != nil {
		t.Errorf("ToJSON: %s", err)
	}
	if !reflect.DeepEqual(info, v) {
		t.Errorf("ToJSON: expected %+v, got %+v", info, v)
	}
}

func TestInfo_String(t *testing.T) {
	randVersionInfo()
	info := Get()
	s := info.String()
	t.Log(s)
}
