package gosysinfo

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestListSysfsClassEntries(t *testing.T) {
	root := t.TempDir()
	base := filepath.Join(root, "class")
	target := filepath.Join(root, "target")

	if err := os.Mkdir(target, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(base, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, filepath.Join(base, "link")); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(base, "file"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := listSysfsClassEntries(base)
	if err != nil {
		t.Fatalf("listSysfsClassEntries() error = %v", err)
	}
	want := []string{"link"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("listSysfsClassEntries() = %v, want %v", got, want)
	}
}

func TestListSysfsClassEntriesMissingBase(t *testing.T) {
	_, err := listSysfsClassEntries(filepath.Join(t.TempDir(), "missing"))
	if err == nil {
		t.Fatal("listSysfsClassEntries() error = nil, want non-nil")
	}
}
