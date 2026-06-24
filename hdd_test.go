package gosysinfo

import (
	"errors"
	"path/filepath"
	"reflect"
	"testing"
)

func TestGetBlockDeviceInfo(t *testing.T) {
	old := checkBlockDevicePathFn
	oldChild := childPartitionsFn
	defer func() {
		checkBlockDevicePathFn = old
		childPartitionsFn = oldChild
	}()
	checkBlockDevicePathFn = func(devname string) (string, error) {
		if devname != "nvme0n1" {
			return "", ErrBlockDeviceNotFound
		}
		return filepath.Join(blockBase, devname), nil
	}
	childPartitionsFn = func(string) ([]string, error) {
		return nil, nil
	}

	reader := mockReader{
		filepath.Join(blockBase, "nvme0n1", "size"):           "2000409264",
		filepath.Join(blockBase, "nvme0n1", "device", "model"): "WD PC SN740",
	}

	got, err := GetBlockDeviceInfo(reader, "nvme0n1")
	if err != nil {
		t.Fatalf("GetBlockDeviceInfo() error = %v", err)
	}
	want := &BlockDeviceInfo{
		Size:           "2000409264",
		Model:          "WD PC SN740",
		Encrypted:      false,
		EncryptionType: "",
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

func TestGetBlockDeviceInfoDMEncrypted(t *testing.T) {
	old := checkBlockDevicePathFn
	defer func() { checkBlockDevicePathFn = old }()
	checkBlockDevicePathFn = func(devname string) (string, error) {
		if devname != "dm-0" {
			return "", ErrBlockDeviceNotFound
		}
		return filepath.Join(blockBase, devname), nil
	}

	reader := mockReader{
		filepath.Join(blockBase, "dm-0", "size"):     "998088704",
		filepath.Join(blockBase, "dm-0", "dm", "uuid"): "CRYPT-LUKS2-94b5a5b67d654944b4bf84d5ccc955a9-root",
	}

	got, err := GetBlockDeviceInfo(reader, "dm-0")
	if err != nil {
		t.Fatalf("GetBlockDeviceInfo() error = %v", err)
	}
	want := &BlockDeviceInfo{
		Size:           "998088704",
		Model:          "",
		Encrypted:      true,
		EncryptionType: "LUKS2",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetBlockDeviceInfo() = %v, want %v", got, want)
	}
}

func TestGetBlockDeviceInfoPartitionEncrypted(t *testing.T) {
	old := checkBlockDevicePathFn
	defer func() { checkBlockDevicePathFn = old }()
	partPath := filepath.Join(blockClassBase, "nvme0n1p2")
	checkBlockDevicePathFn = func(devname string) (string, error) {
		if devname != "nvme0n1p2" {
			return "", ErrBlockDeviceNotFound
		}
		return partPath, nil
	}

	udevContent := "E:ID_FS_TYPE=crypto_LUKS\nE:ID_FS_VERSION=2\n"
	reader := mockReader{
		filepath.Join(partPath, "size"): "998105088",
		filepath.Join(partPath, "dev"):  "259:2",
		filepath.Join(udevDataBase, "b259:2"): udevContent,
	}

	got, err := GetBlockDeviceInfo(reader, "nvme0n1p2")
	if err != nil {
		t.Fatalf("GetBlockDeviceInfo() error = %v", err)
	}
	want := &BlockDeviceInfo{
		Size:           "998105088",
		Model:          "",
		Encrypted:      true,
		EncryptionType: "LUKS2",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetBlockDeviceInfo() = %v, want %v", got, want)
	}
}

func TestGetBlockDeviceInfoDiskWithEncryptedPartition(t *testing.T) {
	old := checkBlockDevicePathFn
	oldChild := childPartitionsFn
	defer func() {
		checkBlockDevicePathFn = old
		childPartitionsFn = oldChild
	}()

	diskPath := filepath.Join(blockBase, "nvme0n1")
	partPath := filepath.Join(blockClassBase, "nvme0n1p2")
	checkBlockDevicePathFn = func(devname string) (string, error) {
		switch devname {
		case "nvme0n1":
			return diskPath, nil
		case "nvme0n1p2":
			return partPath, nil
		default:
			return "", ErrBlockDeviceNotFound
		}
	}
	childPartitionsFn = func(path string) ([]string, error) {
		if path == diskPath {
			return []string{"nvme0n1p1", "nvme0n1p2"}, nil
		}
		return nil, nil
	}

	reader := mockReader{
		filepath.Join(diskPath, "size"):           "2000409264",
		filepath.Join(diskPath, "device", "model"): "WD PC SN740",
		filepath.Join(partPath, "dev"):            "259:2",
		filepath.Join(blockClassBase, "nvme0n1p1", "dev"): "259:1",
		filepath.Join(udevDataBase, "b259:1"):     "E:ID_FS_TYPE=vfat\n",
		filepath.Join(udevDataBase, "b259:2"):     "E:ID_FS_TYPE=crypto_LUKS\nE:ID_FS_VERSION=2\n",
	}

	got, err := GetBlockDeviceInfo(reader, "nvme0n1")
	if err != nil {
		t.Fatalf("GetBlockDeviceInfo() error = %v", err)
	}
	want := &BlockDeviceInfo{
		Size:           "2000409264",
		Model:          "WD PC SN740",
		Encrypted:      true,
		EncryptionType: "LUKS2",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetBlockDeviceInfo() = %v, want %v", got, want)
	}
}

func TestEncryptionFromDMUUID(t *testing.T) {
	tests := []struct {
		uuid    string
		wantEnc bool
		wantTyp string
	}{
		{uuid: "CRYPT-LUKS2-94b5a5b67d654944b4bf84d5ccc955a9-root", wantEnc: true, wantTyp: "LUKS2"},
		{uuid: "CRYPT-LUKS-abc123-root", wantEnc: true, wantTyp: "LUKS"},
		{uuid: "CRYPT-PLAIN-abc123", wantEnc: true, wantTyp: "PLAIN"},
		{uuid: "", wantEnc: false},
		{uuid: "not-crypt", wantEnc: false},
	}
	for _, tt := range tests {
		t.Run(tt.uuid, func(t *testing.T) {
			gotEnc, gotTyp := encryptionFromDMUUID(tt.uuid)
			if gotEnc != tt.wantEnc || gotTyp != tt.wantTyp {
				t.Errorf("encryptionFromDMUUID(%q) = (%v, %q), want (%v, %q)", tt.uuid, gotEnc, gotTyp, tt.wantEnc, tt.wantTyp)
			}
		})
	}
}

func TestParseUdevProperties(t *testing.T) {
	content := "S:disk/by-id/test\nE:ID_FS_TYPE=crypto_LUKS\nE:ID_FS_VERSION=2\n"
	got := parseUdevProperties(content)
	want := map[string]string{
		"ID_FS_TYPE":    "crypto_LUKS",
		"ID_FS_VERSION": "2",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("parseUdevProperties() = %v, want %v", got, want)
	}
}

func TestEncryptionTypeFromUdev(t *testing.T) {
	tests := []struct {
		name    string
		props   map[string]string
		wantEnc bool
		wantTyp string
	}{
		{
			name:    "luks2",
			props:   map[string]string{"ID_FS_TYPE": "crypto_LUKS", "ID_FS_VERSION": "2"},
			wantEnc: true,
			wantTyp: "LUKS2",
		},
		{
			name:    "luks1",
			props:   map[string]string{"ID_FS_TYPE": "crypto_LUKS", "ID_FS_VERSION": "1"},
			wantEnc: true,
			wantTyp: "LUKS",
		},
		{
			name:    "raw crypto type",
			props:   map[string]string{"ID_FS_TYPE": "crypto_LUKS"},
			wantEnc: true,
			wantTyp: "crypto_LUKS",
		},
		{
			name:    "not encrypted",
			props:   map[string]string{"ID_FS_TYPE": "vfat"},
			wantEnc: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEnc, gotTyp := encryptionTypeFromUdev(tt.props)
			if gotEnc != tt.wantEnc || gotTyp != tt.wantTyp {
				t.Errorf("encryptionTypeFromUdev() = (%v, %q), want (%v, %q)", gotEnc, gotTyp, tt.wantEnc, tt.wantTyp)
			}
		})
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
