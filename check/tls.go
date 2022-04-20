package check

import (
	"crypto/tls"
	"encoding/json"
	"net"
	"strings"
	"time"

	"github.com/jreisinger/checkip"
)

const certTLSDialTimeout = 5 * time.Second

type tlsinfo struct {
	SAN     []string
	version string
}

func (t tlsinfo) Summary() string {
	var ss []string
	ss = append(ss, t.version)
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
		Name: "tls",
		Type: checkip.TypeInfo,
	}

	address := net.JoinHostPort(ipaddr.String(), "443")
	conn, err := tls.DialWithDialer(&net.Dialer{Timeout: certTLSDialTimeout}, "tcp", address, &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		return result, newCheckError(err)
	}
	defer conn.Close()

	// search only unique dns names
	dnsSet := make(map[string]struct{})
	var dnsNames []string
	for _, cert := range conn.ConnectionState().PeerCertificates {
		for _, dnsName := range cert.DNSNames {
			if _, ok := dnsSet[dnsName]; ok {
				continue
			}
			dnsNames = append(dnsNames, dnsName)
			dnsSet[dnsName] = struct{}{}
		}
	}

	t := tlsinfo{
		SAN:     dnsNames,
		version: tlsFormat(conn.ConnectionState().Version),
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
