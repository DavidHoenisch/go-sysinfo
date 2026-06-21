package gosysinfo

import (
	"reflect"
	"testing"
)

func TestGetRAM(t *testing.T) {
	reader := mockReader{
		procMemInfo: `MemTotal:       31939132 kB
MemFree:         7445796 kB
MemAvailable:   23866632 kB
Buffers:              32 kB
`,
	}

	got := GetRAM(reader)
	want := &RAMInfo{
		MemTotal:     "31939132 kB",
		MemFree:      "7445796 kB",
		MemAvailable: "23866632 kB",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetRAM() = %v, want %v", got, want)
	}
}

func TestGetRAEmpty(t *testing.T) {
	got := GetRAM(mockReader{})
	want := &RAMInfo{}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetRAM() = %v, want %v", got, want)
	}
}
