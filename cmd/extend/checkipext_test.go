package main

import (
	"fmt"
	"testing"

	"github.com/jreisinger/checkip/check"
)

func setMockConfig(t *testing.T, fn func(key string) (string, error)) {
	origGetConfigValue := check.GetConfigValue
	check.GetConfigValue = fn
	t.Cleanup(func() {
		check.GetConfigValue = origGetConfigValue
	})
}

func setExtendMockConfig(t *testing.T, value string) {
	setMockConfig(t, func(key string) (string, error) {
		if key == "CHECKS" {
			return value, nil
		}
		return "", fmt.Errorf("unexpected key %s received", key)
	})
}

func TestFirstGetListOfChecks(t *testing.T) {

	setExtendMockConfig(t, "")
	list, err := getListOfChecks("", "")
	if len(list) != len(check.Definitions) {
		t.Fatalf("Default list of check : len(list) = %d not have all definitions than %d", len(list), len(check.Definitions))
	}

	// Test whith config "CHECKS: db-ip.com, dns mx, dns name"
	setExtendMockConfig(t, "db-ip.com, dns mx, dns name")
	list, _ = getListOfChecks("", "")
	if len(list) != 3 {
		t.Fatalf("Config CHECK with 3 checks : len(list) = %d not have all 3 checks", len(list))
	}

	// Test whith selected checks
	list, _ = getListOfChecks("tls, ping", "")
	if len(list) != 2 {
		t.Fatalf("only 2 selected checks : len(list) = %d not have all 2 checks", len(list))
	}

	// Test whith selected checks and append
	list, _ = getListOfChecks("tls, ping", "censys.io, shodan.io, spur.io")
	if len(list) != 5 {
		t.Fatalf("only 2 selected checks and 3 append : len(list) = %d not have all 5 checks", len(list))
	}

	// Test whith bad checks
	list, err = getListOfChecks("something", "")
	if len(list) != 0 && err != nil {
		t.Fatalf("no list : len(list) = %d and errror message : %s", len(list), err)
	}

	// Test whith bad append checks
	list, err = getListOfChecks("", "something")
	if len(list) != 3 && err != nil {
		t.Fatalf("Default list : len(list) = %d and errror message : %s", len(list), err)
	}

}
