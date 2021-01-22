package main

import "testing"

func TestIsAvailable(t *testing.T) {
	type testpair struct {
		checkName string
		available bool
	}
	testpairs := []testpair{
		{"abuseipdb", true},
		{"as", true},
		{"dns", true},
		{"geo", true},
		{"ipsum", true},
		{"otx", true},
		{"threatcrowd", true},
		{"virustotal", true},
		{"foo", false},
		{"", false},
	}
	for _, tp := range testpairs {
		_, ok := isAvailable(tp.checkName)
		if ok != tp.available {
			t.Errorf("'%s' should be %t but is %t", tp.checkName, tp.available, ok)
		}
	}
}
