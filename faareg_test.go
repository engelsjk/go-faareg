package faareg

import "testing"

func TestRegistration(t *testing.T) {
	// Reference: https://registry.faa.gov/AircraftInquiry/Search/NNumberResult?nNumberTxt=265FT
	reg, err := GetRegistration("265FT")
	if err != nil {
		t.Fatalf(err.Error())
	}
	test(t, reg.Aircraft.Registration, "N265FT")
	test(t, reg.Aircraft.SerialNumber, "465-46")
	test(t, reg.Aircraft.ManufacturerName, "ROCKWELL INTERNATIONAL CORP")
	test(t, reg.Aircraft.Model, "NA-265-65")
	test(t, reg.Aircraft.ModeSCodeHex, "A2921F")

}

func test(t *testing.T, got string, wanted string) {
	if got != wanted {
		t.Fatalf("wanted %s, got %s", wanted, got)
	}
}
