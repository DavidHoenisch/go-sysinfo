package gosysinfo

import (
	"errors"
	"path/filepath"
	"reflect"
	"testing"
)

func TestGetBatteryInfo(t *testing.T) {
	oldPath := checkBatteryPathFn
	defer func() { checkBatteryPathFn = oldPath }()
	checkBatteryPathFn = func(name string) (string, error) {
		if name != "BAT0" {
			return "", ErrBatteryNotFound
		}
		return filepath.Join(powerSupplyBase, name), nil
	}

	base := filepath.Join(powerSupplyBase, "BAT0")
	reader := mockReader{
		filepath.Join(base, "type"):                            "Battery",
		filepath.Join(base, "status"):                          "Discharging",
		filepath.Join(base, "capacity"):                        "87",
		filepath.Join(base, "capacity_level"):                  "Normal",
		filepath.Join(base, "energy_now"):                      "59800000",
		filepath.Join(base, "energy_full"):                     "68633000",
		filepath.Join(base, "energy_full_design"):              "73000000",
		filepath.Join(base, "power_now"):                       "12500000",
		filepath.Join(base, "voltage_now"):                     "16800000",
		filepath.Join(base, "voltage_min_design"):              "15939000",
		filepath.Join(base, "manufacturer"):                    "ASUSTeK",
		filepath.Join(base, "model_name"):                      "ASUS Battery",
		filepath.Join(base, "serial_number"):                   "ABC123",
		filepath.Join(base, "technology"):                      "Li-ion",
		filepath.Join(base, "present"):                         "1",
		filepath.Join(base, "cycle_count"):                     "42",
		filepath.Join(base, "charge_control_end_threshold"):    "100",
	}

	got, err := GetBatteryInfo(reader, "BAT0")
	if err != nil {
		t.Fatalf("GetBatteryInfo() error = %v", err)
	}
	want := &BatteryInfo{
		Status:                    "Discharging",
		Capacity:                  "87",
		CapacityLevel:             "Normal",
		EnergyNow:                 "59800000",
		EnergyFull:                "68633000",
		EnergyFullDesign:          "73000000",
		PowerNow:                  "12500000",
		VoltageNow:                "16800000",
		VoltageMinDesign:          "15939000",
		Manufacturer:              "ASUSTeK",
		ModelName:                 "ASUS Battery",
		SerialNumber:              "ABC123",
		Technology:                "Li-ion",
		Present:                   "1",
		CycleCount:                "42",
		ChargeControlEndThreshold: "100",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetBatteryInfo() = %v, want %v", got, want)
	}
}

func TestGetBatteryInfoNotFound(t *testing.T) {
	_, err := GetBatteryInfo(mockReader{}, "missing0")
	if !errors.Is(err, ErrBatteryNotFound) {
		t.Fatalf("GetBatteryInfo() error = %v, want %v", err, ErrBatteryNotFound)
	}
}

func TestGetBatteryInfoNotBattery(t *testing.T) {
	oldPath := checkBatteryPathFn
	defer func() { checkBatteryPathFn = oldPath }()
	checkBatteryPathFn = func(name string) (string, error) {
		return filepath.Join(powerSupplyBase, name), nil
	}

	reader := mockReader{
		filepath.Join(powerSupplyBase, "AC0", "type"): "Mains",
	}

	_, err := GetBatteryInfo(reader, "AC0")
	if !errors.Is(err, ErrBatteryNotFound) {
		t.Fatalf("GetBatteryInfo() error = %v, want %v", err, ErrBatteryNotFound)
	}
}

func TestFilterBatteryNames(t *testing.T) {
	oldReadType := readPowerSupplyType
	defer func() { readPowerSupplyType = oldReadType }()
	readPowerSupplyType = func(name string) string {
		switch name {
		case "BAT0", "hid-0018:04F3:4359.0005-battery":
			return "Battery"
		case "AC0":
			return "Mains"
		default:
			return ""
		}
	}

	got := filterBatteryNames([]string{
		"AC0",
		"BAT0",
		"hid-0018:04F3:4359.0005-battery",
		"ucsi-source-psy-USBC000:001",
	})
	want := []string{"BAT0", "hid-0018:04F3:4359.0005-battery"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("filterBatteryNames() = %v, want %v", got, want)
	}
}
