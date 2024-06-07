package check

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type awsIpRanges struct {
	IsOn bool
	Info struct {
		IpPrefix           string   `json:"ip_prefix"`
		Region             string   `json:"region"`
		Services           []string `json:"services"`
		NetworkBorderGroup string   `json:"network_border_group"`
	}
}

// Json implements checkip.Info
func (a awsIpRanges) Json() ([]byte, error) {
	return json.Marshal(&a)
}

// Summary implements checkip.Info
func (a awsIpRanges) Summary() string {
	if a.IsOn {
		return fmt.Sprintf("%t, prefix: %s, region: %s, sevices: %v",
			a.IsOn, a.Info.IpPrefix, a.Info.Region, a.Info.Services)
	}
	return fmt.Sprintf("%t", a.IsOn)
}

// IsOnAWS checks if ipaddr belongs to AWS. If so it provides info about the IP
// address. It gets the info from https://ip-ranges.amazonaws.com/ip-ranges.json
func IsOnAWS(ipaddr net.IP) (Check, error) {
	result := Check{
		Description: "is on AWS",
	}
	resp := struct {
		Prefixes []struct {
			IpPrefix           string `json:"ip_prefix"`
			Region             string `json:"region"`
			Service            string `json:"service"`
			NetworkBorderGroup string `json:"network_border_group"`
		} `json:"prefixes"`
	}{}

	filename, err := getCachePath("aws-ip-ranges.json")
	if err != nil {
		return result, err
	}
	if err := updateFile(filename, "https://ip-ranges.amazonaws.com/ip-ranges.json", ""); err != nil {
		return result, err
	}

	f, err := os.Open(filename)
	if err != nil {
		return result, err
	}
	dec := json.NewDecoder(f)
	if err := dec.Decode(&resp); err != nil {
		return result, err
	}

	var a awsIpRanges
	for _, prefix := range resp.Prefixes {
		_, network, err := net.ParseCIDR(prefix.IpPrefix)
		if err != nil {
			return result, fmt.Errorf("parse CIDR %q: %v", prefix.IpPrefix, err)
		}
		if network.Contains(ipaddr) {
			a.IsOn = true
			a.Info.IpPrefix = prefix.IpPrefix
			a.Info.NetworkBorderGroup = prefix.NetworkBorderGroup
			a.Info.Region = prefix.Region
			a.Info.Services = append(a.Info.Services, prefix.Service)
		}

	}
	result.IpAddrInfo = a
	return result, nil
}
