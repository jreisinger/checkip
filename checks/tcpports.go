package checks

import (
	"encoding/json"
	"fmt"
	"net"
	"sort"
	"strings"
	"time"

	"github.com/jreisinger/checkip/check"
	"github.com/jreisinger/nmapservices"
)

// TcpPorts represents a TCP ports.
type TcpPort struct {
	Name   string // service name like ssh
	Number int16
}

// OpenTcpPorts are the TCP ports that were found open on the given IP address.
type OpenTcpPorts []TcpPort

type byPortNumber OpenTcpPorts

func (x byPortNumber) Len() int           { return len(x) }
func (x byPortNumber) Less(i, j int) bool { return x[i].Number < x[j].Number }
func (x byPortNumber) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

func (t OpenTcpPorts) Summary() string {
	var out []string
	sort.Sort(byPortNumber(t))
	for _, p := range t {
		out = append(out, fmt.Sprintf("%d (%s)", p.Number, p.Name))
	}
	return strings.Join(out, ", ")
}

func (t OpenTcpPorts) JsonString() (string, error) {
	b, err := json.Marshal(t)
	return string(b), err
}

// TcpPorts tries to connect to the 1000 TCP ports that are most often found
// open on Internet hosts. Then it reports which of those ports are open on the
// given IP address.
func TcpPorts(ipaddr net.IP) (check.Result, error) {
	openports, err := scan(ipaddr, 1000)
	if err != nil {
		return check.Result{}, check.NewError(err)
	}

	return check.Result{
		Name: "Open TCP ports",
		Type: check.TypeInfo,
		Info: OpenTcpPorts(openports),
	}, nil
}

func scan(ip net.IP, top int) (OpenTcpPorts, error) {
	ports := make(chan TcpPort)
	results := make(chan TcpPort)

	for i := 0; i < 100; i++ {
		go scanner(ip, ports, results)
	}

	var openports OpenTcpPorts

	services, err := nmapservices.Get()
	if err != nil {
		return openports, err
	}

	var topPorts []TcpPort

	for _, s := range services.Tcp().Top(top) {
		topPorts = append(topPorts, TcpPort{Number: s.Port, Name: s.Name})
	}

	go func() {
		for _, port := range topPorts {
			ports <- port
		}
	}()

	for range topPorts {
		port := <-results
		if port.Number != 0 {
			openports = append(openports, port)
		}
	}

	return openports, nil
}

func scanner(ip net.IP, ports, results chan TcpPort) {
	for port := range ports {
		addr := fmt.Sprintf("%s:%d", ip, port.Number)
		conn, err := net.DialTimeout("tcp", addr, time.Second*2)
		if err != nil {
			results <- TcpPort{Number: 0}
			continue
		}
		conn.Close()
		results <- port
	}
}
