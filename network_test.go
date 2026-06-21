package gosysinfo

import (
	"errors"
	"path/filepath"
	"reflect"
	"testing"
)

func TestGetNetworkConnectionInfo(t *testing.T) {
	restore := overrideNetworkPathCheck(t, "eth0")
	defer restore()

	reader := mockReader{
		filepath.Join(netBase, "eth0", "address"): "aa:bb:cc:dd:ee:ff",
	}

	got, err := GetNetworkConnectionInfo(reader, "eth0")
	if err != nil {
		t.Fatalf("GetNetworkConnectionInfo() error = %v", err)
	}
	want := &NetworkConnectionInfo{MACAddress: "aa:bb:cc:dd:ee:ff"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetNetworkConnectionInfo() = %v, want %v", got, want)
	}
}

func TestGetNetworkConnectionInfoNotFound(t *testing.T) {
	_, err := GetNetworkConnectionInfo(mockReader{}, "missing0")
	if !errors.Is(err, ErrNetworkNotFound) {
		t.Fatalf("GetNetworkConnectionInfo() error = %v, want %v", err, ErrNetworkNotFound)
	}
}

func TestGetNetworkIPInfoIPv6(t *testing.T) {
	restore := overrideNetworkPathCheck(t, "wlan0")
	defer restore()

	reader := mockReader{
		procNetIfInet6: "fe80000000000000e2e258fffee7581a 03 40 20 80 wlan0\nfe8000000000000074d41de8f8619095 04 40 20 80 tailscale0\n",
		procNetFibTrie: "",
		procNetRoute:   "",
	}

	got, err := GetNetworkIPInfo(reader, "wlan0")
	if err != nil {
		t.Fatalf("GetNetworkIPInfo() error = %v", err)
	}
	if len(got.IPv6) != 1 {
		t.Fatalf("GetNetworkIPInfo() IPv6 len = %d, want 1", len(got.IPv6))
	}
	if got.IPv6[0] != "fe80::e2e2:58ff:fee7:581a" {
		t.Errorf("GetNetworkIPInfo() IPv6[0] = %q", got.IPv6[0])
	}
}

func TestGetNetworkIPInfoIPv4(t *testing.T) {
	restore := overrideNetworkPathCheck(t, "eth0")
	defer restore()

	reader := mockReader{
		procNetIfInet6: "",
		procNetFibTrie: "Local:\n  +-- 0.0.0.0/0 3 1 5\n     +-- 192.168.10.0/24 2 0 2\n        |-- 192.168.10.20\n           /32 host LOCAL\n",
		procNetRoute:   "Iface\tDestination\tGateway \tFlags\tRefCnt\tUse\tMetric\tMask\t\tMTU\tWindow\tIRTT\neth0\t000AA8C0\t00000000\t0001\t0\t0\t0\t00FFFFFF\t0\t0\t0\n",
	}

	got, err := GetNetworkIPInfo(reader, "eth0")
	if err != nil {
		t.Fatalf("GetNetworkIPInfo() error = %v", err)
	}
	if len(got.IPv4) != 1 || got.IPv4[0] != "192.168.10.20" {
		t.Errorf("GetNetworkIPInfo() IPv4 = %v, want [192.168.10.20]", got.IPv4)
	}
}

func TestGetNetworkIPInfoNotFound(t *testing.T) {
	_, err := GetNetworkIPInfo(mockReader{}, "missing0")
	if !errors.Is(err, ErrNetworkNotFound) {
		t.Fatalf("GetNetworkIPInfo() error = %v, want %v", err, ErrNetworkNotFound)
	}
}

func TestGetNetworkStatisticsReadsPerFilePath(t *testing.T) {
	restore := overrideNetworkPathCheck(t, "eth0")
	defer restore()

	reader := mockReader{
		filepath.Join(netBase, "eth0", "mtu"): "1500",
	}

	got, err := GetNetworkStatistics(reader, "eth0")
	if err != nil {
		t.Fatalf("GetNetworkStatistics() error = %v", err)
	}
	if got.Mtu != "1500" {
		t.Errorf("GetNetworkStatistics() Mtu = %q, want %q", got.Mtu, "1500")
	}
	if got.Address != "" {
		t.Errorf("GetNetworkStatistics() Address = %q, want empty when mtu path only mocked", got.Address)
	}
}

func TestGetNetworkStatisticsNotFound(t *testing.T) {
	_, err := GetNetworkStatistics(mockReader{}, "missing0")
	if !errors.Is(err, ErrNetworkNotFound) {
		t.Fatalf("GetNetworkStatistics() error = %v, want %v", err, ErrNetworkNotFound)
	}
}

func overrideNetworkPathCheck(t *testing.T, ifname string) func() {
	t.Helper()
	old := checkNetworkPathByNameFn
	checkNetworkPathByNameFn = func(name string) (string, error) {
		if name != ifname {
			return "", ErrNetworkNotFound
		}
		return filepath.Join(netBase, name), nil
	}
	return func() { checkNetworkPathByNameFn = old }
}
