package gosysinfo

import (
	"errors"
	"path/filepath"
	"reflect"
	"testing"
)

func TestGetBlockDeviceInfo(t *testing.T) {
	old := checkBlockDevicePathFn
	defer func() { checkBlockDevicePathFn = old }()
	checkBlockDevicePathFn = func(devname string) (string, error) {
		if devname != "nvme0n1" {
			return "", ErrBlockDeviceNotFound
		}
		return filepath.Join(blockBase, devname), nil
	}

	reader := mockReader{
		filepath.Join(blockBase, "nvme0n1", "size"):          "2000409264",
		filepath.Join(blockBase, "nvme0n1", "device", "model"): "WD PC SN740",
	}

	got, err := GetBlockDeviceInfo(reader, "nvme0n1")
	if err != nil {
		t.Fatalf("GetBlockDeviceInfo() error = %v", err)
	}
	want := &BlockDeviceInfo{
		Size:  "2000409264",
		Model: "WD PC SN740",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetBlockDeviceInfo() = %v, want %v", got, want)
	}
}

func TestGetBlockDeviceInfoNotFound(t *testing.T) {
	_, err := GetBlockDeviceInfo(mockReader{}, "missing0")
	if !errors.Is(err, ErrBlockDeviceNotFound) {
		t.Fatalf("GetBlockDeviceInfo() error = %v, want %v", err, ErrBlockDeviceNotFound)
	}
}

func TestIsPartition(t *testing.T) {
	tests := []struct {
		name string
		dev  string
		want bool
	}{
		{name: "whole device", dev: "nvme0n1", want: false},
		{name: "partition", dev: "nvme0n1p1", want: true},
		{name: "short name", dev: "sda", want: false},
		{name: "sda partition", dev: "sda1", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isPartition(tt.dev); got != tt.want {
				t.Errorf("isPartition(%q) = %v, want %v", tt.dev, got, tt.want)
			}
		})
	}
}
