package gosysinfo

import (
	"reflect"
	"testing"
)

func TestGetCPU(t *testing.T) {
	reader := mockReader{
		procCPUInfo: `processor	: 0
vendor_id	: AuthenticAMD
cpu family	: 26
model		: 36
model name	: AMD Ryzen AI 9 HX 370 w/ Radeon 890M
processor	: 1
processor	: 2
`,
	}

	got := GetCPU(reader)
	want := &CPUInfo{
		ModelName: "AMD Ryzen AI 9 HX 370 w/ Radeon 890M",
		VendorID:  "AuthenticAMD",
		CPUFamily: "26",
		Model:     "36",
		CoreCount: "3",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetCPU() = %v, want %v", got, want)
	}
}

func TestGetCPUEmpty(t *testing.T) {
	got := GetCPU(mockReader{})
	want := &CPUInfo{}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetCPU() = %v, want %v", got, want)
	}
}
