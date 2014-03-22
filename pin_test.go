package embd

import "testing"

func TestPinMapLookup(t *testing.T) {
	var tests = []struct {
		key interface{}
		id  string
	}{
		{"10", "P1_1"},
		{10, "P1_1"},
		{"P1_2", "P1_2"},
		{"GPIO10", "P1_2"},
	}
	var pinMap = PinMap{
		&PinDesc{ID: "P1_1", Aliases: []string{"AN1", "10"}},
		&PinDesc{ID: "P1_2", Aliases: []string{"GPIO10"}},
	}
	for _, test := range tests {
		pd, found := pinMap.Lookup(test.key)
		if !found {
			t.Errorf("Could not find a descriptor for %q", test.key)
			continue
		}
		if pd.ID != test.id {
			t.Errorf("Looking up %q: got %v, want %v", test.key, pd.ID, test.id)
		}
	}
}
