package check

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strings"
	"syscall"
	"time"
)

const certTLSDialTimeout = 5 * time.Second

type tlsinfo struct {
	SAN     []string
	Version uint16
	Expiry  time.Time
}

func (t tlsinfo) Summary() string {
	var ss []string

	ver := tlsFormat(t.Version)
	if oldTlsVersion(t.Version) {
		ver += "!!"
	}
	ss = append(ss, ver)

	exp := "exp. " + t.Expiry.Format("2006/01/02")
	if expiredCert(t.Expiry) {
		exp += "!!"
	}
	ss = append(ss, exp)

	ss = append(ss, t.SAN...)
	return strings.Join(ss, ", ")
}

func (t tlsinfo) Json() ([]byte, error) {
	return json.Marshal(t)
}

// Tls finds out TLS information by connecting to the ipaddr and TCP port 443.
func Tls(ipaddr net.IP) (Check, error) {
	result := Check{
		Description: "tls",
		Type:        InfoAndIsMalicious,
	}

	address := net.JoinHostPort(ipaddr.String(), "443")
	conn, err := tls.DialWithDialer(&net.Dialer{Timeout: certTLSDialTimeout}, "tcp", address, &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		// Ignore ECONNREFUSED error.
		var s syscall.Errno
		if errors.As(err, &s) {
			if s == syscall.ECONNREFUSED {
				return result, nil
			}
		}

		return result, newCheckError(fmt.Errorf("connect to %s: %v", address, err))
	}
	defer conn.Close()

	// search only unique dns names
	dnsSet := make(map[string]struct{})
	var dnsNames []string
	var expiry time.Time
	for _, cert := range conn.ConnectionState().PeerCertificates {
		for i, dnsName := range cert.DNSNames {
			if i == 0 || cert.NotAfter.Before(expiry) {
				expiry = cert.NotAfter
			}
			if _, ok := dnsSet[dnsName]; ok {
				continue
			}
			dnsNames = append(dnsNames, dnsName)
			dnsSet[dnsName] = struct{}{}
		}
	}

	t := tlsinfo{
		SAN:     dnsNames,
		Version: conn.ConnectionState().Version,
		Expiry:  expiry,
	}

	result.IpAddrInfo = t

	if oldTlsVersion(conn.ConnectionState().Version) || expiredCert(t.Expiry) {
		result.IpAddrIsMalicious = true
	}

	return result, nil
}

func oldTlsVersion(tlsVersion uint16) bool {
	if tlsVersion == tls.VersionTLS12 || tlsVersion == tls.VersionTLS13 {
		return false
	}
	return true
}

func expiredCert(expiryDate time.Time) bool {
	return expiryDate.Before(time.Now())
}

func tlsFormat(tlsVersion uint16) string {
	switch tlsVersion {
	case 0:
		return ""
	case tls.VersionTLS10:
		return "TLS 1.0"
	case tls.VersionTLS11:
		return "TLS 1.1"
	case tls.VersionTLS12:
		return "TLS 1.2"
	case tls.VersionTLS13:
		return "TLS 1.3"
	default:
		return "TLS Version %d (unknown)"
	}
}
