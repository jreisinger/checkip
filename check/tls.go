package check

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"net"
	"strings"
	"syscall"
	"time"

	"github.com/jreisinger/checkip"
)

const certTLSDialTimeout = 5 * time.Second

type tlsinfo struct {
	SAN     []string
	Version string
	Expiry  time.Time
}

func (t tlsinfo) Summary() string {
	var ss []string
	ss = append(ss, t.Version)
	ss = append(ss, t.Expiry.Format("2006/01/02"))
	ss = append(ss, t.SAN...)
	return strings.Join(ss, ", ")
}

func (t tlsinfo) Json() ([]byte, error) {
	return json.Marshal(t)
}

// Tls finds out TLS version and SANs by connecting to the ipaddr and TCP port
// 443.
func Tls(ipaddr net.IP) (checkip.Result, error) {
	result := checkip.Result{
		Name: "cert",
		Type: checkip.TypeInfo,
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

		return result, newCheckError(err)
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
		Version: tlsFormat(conn.ConnectionState().Version),
		Expiry:  expiry,
	}

	result.Info = t

	return result, nil
}

func tlsFormat(tlsVersion uint16) string {
	switch tlsVersion {
	case 0:
		return ""
	case tls.VersionSSL30:
		return "SSLv3 - Deprecated!"
	case tls.VersionTLS10:
		return "TLS 1.0 - Deprecated!"
	case tls.VersionTLS11:
		return "TLS 1.1 - Deprecated!"
	case tls.VersionTLS12:
		return "TLS 1.2"
	case tls.VersionTLS13:
		return "TLS 1.3"
	default:
		return "TLS Version %d (unknown)"
	}
}