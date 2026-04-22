// Package check contains functions that can check an IP address.
package check

import "strings"

// Debug set by main flag
var Debug bool

// GetConfigValue export getConfigValue function for main
var GetConfigValue = func(key string) (string, error) {
	return getConfigValue(key)
}

// ExtDefinitions more definition in extended branch
var ExtDefinitions = []Definition{
	{Name: "Misp", Run: Misp, Cache: CacheNone},
	{Name: "MyDB", Run: MyDB, Cache: CacheNone},
	{Name: "IOCLoc", Run: IOCLoc},
	{Name: "Onyphe", Run: Onyphe, NewInfo: func() IpInfo { return &onyphe{} }},
	{Name: "IpAPI", Run: IpAPI, NewInfo: func() IpInfo { return &ipapi{} }},
}

// InList return selected definitions
func InList(list []string, definitions []Definition) []Definition {
	filtered := make([]Definition, 0, len(definitions))
	for _, definition := range definitions {
		for _, l := range list {
			l = strings.Replace(l, " ", "", -1)
			if definition.Name != l {
				continue
			}
			filtered = append(filtered, definition)
		}
	}
	return filtered
}
