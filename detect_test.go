package embd

import "testing"

func TestKernelVersionParse(t *testing.T) {
	var tests = []struct {
		versionStr          string
		major, minor, patch int
	}{
		{
			"3.8.2",
			3, 8, 2,
		},
		{
			"3.8.10+",
			3, 8, 10,
		},
	}
	for _, test := range tests {
		major, minor, patch, err := parseVersion(test.versionStr)
		if err != nil {
			t.Errorf("Failed parsing %q: %v", test.versionStr, err)
			continue
		}
		if major != test.major || minor != test.minor || patch != test.patch {
			t.Errorf("Parse of %q: got (%v, %v, %v) want (%v, %v, %v)", test.versionStr, major, minor, patch, test.major, test.minor, test.patch)
		}
	}
}
