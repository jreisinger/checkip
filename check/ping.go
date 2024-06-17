package check

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/go-ping/ping"
)

// Ping sends five pings (ICMP echo request packets) to the ippaddr and returns
// the statistics.
func Ping(ipaddr net.IP) (Check, error) {
	result := Check{
		Description: "ping",
		Type:        Info,
	}

	pinger, err := ping.NewPinger(ipaddr.String())
	if err != nil {
		return result, newCheckError(err)
	}
	pinger.Count = 5
	pinger.Timeout = time.Second * 5
	err = pinger.Run() // Blocks until finished.
	if err != nil {
		return result, newCheckError(err)
	}

	ps := pinger.Statistics() // get send/receive/duplicate/rtt stats
	result.IpAddrInfo = stats(*ps)

	return result, nil
}

type stats ping.Statistics

func (s stats) Summary() string {
	return fmt.Sprintf("%.0f%% packet loss (%d/%d), avg round-trip %d ms", s.PacketLoss, s.PacketsSent, s.PacketsRecv, s.AvgRtt.Milliseconds())
}

func (s stats) Json() ([]byte, error) {
	return json.Marshal(s)
}
