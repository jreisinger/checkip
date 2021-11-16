package checks

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/go-ping/ping"
	"github.com/jreisinger/checkip/check"
)

func CheckPing(ipaddr net.IP) (check.Result, error) {
	pinger, err := ping.NewPinger(ipaddr.String())
	if err != nil {
		return check.Result{}, check.NewError(err)
	}
	pinger.Count = 3
	pinger.Timeout = time.Second * 5
	err = pinger.Run() // Blocks until finished.
	if err != nil {
		return check.Result{}, check.NewError(err)
	}

	ps := pinger.Statistics() // get send/receive/duplicate/rtt stats

	return check.Result{
		Name:            "ping",
		Type:            check.TypeInfo,
		Info:            stats(*ps),
		IPaddrMalicious: false,
	}, nil
}

type stats ping.Statistics

func (s stats) Summary() string {
	return fmt.Sprintf("%.0f%% packet loss, sent %d, recv %d, avg round-trip %d ms", s.PacketLoss, s.PacketsSent, s.PacketsRecv, s.AvgRtt.Milliseconds())
}

func (s stats) JsonString() (string, error) {
	b, err := json.Marshal(s)
	return string(b), err
}
