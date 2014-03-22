package embd

import "testing"

func TestPinMapLookup(t *testing.T) {
	var tests = []struct {
		key interface{}
		cap int
		id  string

		found bool
	}{
		{"10", CapAnalog, "P1_1", true},
		{10, CapAnalog, "P1_1", true},
		{"10", CapNormal, "P1_2", true},
		{"P1_2", CapNormal, "P1_2", true},
		{"P1_2", CapAnalog, "P1_2", true},
		{"GPIO10", CapNormal, "P1_2", true},
		{key: "NOTTHERE", found: false},
	}
	var pinMap = PinMap{
		&PinDesc{ID: "P1_1", Aliases: []string{"AN1", "10"}, Caps: CapAnalog},
		&PinDesc{ID: "P1_2", Aliases: []string{"10", "GPIO10"}, Caps: CapNormal},
	}
	for _, test := range tests {
		pd, found := pinMap.Lookup(test.key, test.cap)
		if found != test.found {
			t.Errorf("Outcome mismatch for %v: got found = %v, expected found = %v", test.key, found, test.found)
			continue
		}
		if !found {
			continue
		}
		if pd.ID != test.id {
			var capStr string
			switch test.cap {
			case CapNormal:
				capStr = "CapNormal"
			case CapAnalog:
				capStr = "CapAnalog"
			default:
				t.Fatalf("Unknown cap %v", test.cap)
			}
			t.Errorf("Looking up %q with %v: got %v, want %v", test.key, capStr, pd.ID, test.id)
		}
	}
}
