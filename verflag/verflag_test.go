package verflag

import (
	"math/rand"
	"os"
	"reflect"
	"testing"

	"github.com/spf13/pflag"
)

func TestVersionValue_IsBoolFlag(t *testing.T) {
	var v versionValue
	if !v.IsBoolFlag() {
		t.Errorf("IsBoolFlag() should return true")
	}
}

func Test_versionValue_Get(t *testing.T) {
	vt := VersionTrue
	vf := VersionFalse
	vr := VersionRaw
	var vn *versionValue = nil
	tests := []struct {
		name string
		v    *versionValue
		want interface{}
	}{
		{
			name: "true",
			v:    &vt,
			want: &vt,
		},
		{
			name: "false",
			v:    &vf,
			want: &vf,
		},
		{
			name: "raw",
			v:    &vr,
			want: &vr,
		},
		{
			name: "nil",
			v:    vn,
			want: vn,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.Get(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_versionValue_Set(t *testing.T) {
	v := VersionFalse
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		v       versionValue
		args    args
		wantErr bool
	}{
		{
			name:    "true",
			v:       v,
			args:    args{s: "true"},
			wantErr: false,
		},
		{
			name:    "false",
			v:       v,
			args:    args{s: "false"},
			wantErr: false,
		},
		{
			name:    "raw",
			v:       v,
			args:    args{s: "raw"},
			wantErr: false,
		},
		{
			name:    "invalid",
			v:       v,
			args:    args{s: "invalid"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.v.Set(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_versionValue_String(t *testing.T) {
	vt := VersionTrue
	vf := VersionFalse
	vr := VersionRaw
	var vo versionValue
	for vo >= vt && vo <= vr {
		vo = versionValue(rand.Int())
	}
	var vn *versionValue = nil
	tests := []struct {
		name string
		v    *versionValue
		want string
	}{
		{
			name: "true",
			v:    &vt,
			want: "true",
		},
		{
			name: "false",
			v:    &vf,
			want: "false",
		},
		{
			name: "raw",
			v:    &vr,
			want: "raw",
		},
		{
			name: "nil",
			v:    vn,
			want: "",
		},
		{
			name: "other",
			v:    &vo,
			want: "false",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVersionValue_Type(t *testing.T) {
	v := versionValue(rand.Int())
	if got := v.Type(); got != "version" {
		t.Errorf("Type() = %v, want %v", got, "version")
	}
}

func TestAddFlags(t *testing.T) {
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	AddFlags(fs)
	err := fs.Parse([]string{"--version=true"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	flag := fs.Lookup("version")
	if flag == nil {
		t.Fatalf("flag not found")
	}
	if !flag.Changed {
		t.Fatalf("flag not changed")
	}
	if flag.Value.String() != "true" {
		t.Fatalf("flag value unexpected: %v", flag.Value)
	}
}

func TestPrintAndExitIfRequested_false(t *testing.T) {
	origin := os.Stdout
	defer func() {
		os.Stdout = origin
	}()
	_, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	os.Stdout = w

	*versionFlag = VersionFalse
	PrintAndExitIfRequested()
}

func TestPrintAndExitIfRequested_true(t *testing.T) {
	origin := os.Stdout
	defer func() {
		os.Stdout = origin
	}()
	_, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	os.Stdout = w

	*versionFlag = VersionTrue
	defer func() {
		if err := recover(); err == nil {
			t.Fatalf("expected panic")
		}
	}()
	PrintAndExitIfRequested()
}

func TestPrintAndExitIfRequested_raw(t *testing.T) {
	origin := os.Stdout
	defer func() {
		os.Stdout = origin
	}()
	_, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	os.Stdout = w

	*versionFlag = VersionRaw
	defer func() {
		if err := recover(); err == nil {
			t.Fatalf("expected panic")
		}
	}()
	PrintAndExitIfRequested()
}
