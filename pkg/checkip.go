// Checkip checks an IP address using various public services. An IP address is
// checked by running one or more Checkers. There are two kinds of Checkers. An
// InfoChecker just gathers some useful information about the IP address. A
// SecChecker says whether the IP address is considered malicious or not.
package checkip

import (
	"encoding/json"
	"fmt"
	"net"
	"regexp"
	"sync"
)

// Checker runs a check of an IP address. String() returns checker's name.
type Checker interface {
	Check(ip net.IP) error
	fmt.Stringer
}

// InfoChecker finds information about an IP address.
type InfoChecker interface {
	Info() string
	Checker
}

// SecChecker checks an IP address from the security point of view.
type SecChecker interface {
	IsMalicious() bool
	Checker
}

type InfoSecChecker interface {
	InfoChecker
	SecChecker
}

// Result holds the result of a check.
type Result struct {
	Name        string
	Type        string
	Data        Checker
	Info        string
	IsMalicious bool
	Err         error `json:"-"` // omit error from marshalling - https://bit.ly/2ZZOM7C
	ErrMsg      string
}

// Run runs checkers concurrently checking the ipaddr.
func Run(checkers []Checker, ipaddr net.IP) []Result {
	var res []Result

	var wg sync.WaitGroup
	for _, chk := range checkers {
		wg.Add(1)
		go func(c Checker) {
			defer wg.Done()
			var errMsg string
			err := c.Check(ipaddr)
			if err != nil {
				errMsg = redactSecrets(err.Error())
			}
			switch v := c.(type) {
			case InfoSecChecker:
				r := Result{Name: v.String(), Type: "InfoSec", Data: v, Info: v.Info(), IsMalicious: v.IsMalicious(), Err: err, ErrMsg: errMsg}
				res = append(res, r)
			case InfoChecker:
				r := Result{Name: v.String(), Type: "Info", Data: v, Info: v.Info(), Err: err, ErrMsg: errMsg}
				res = append(res, r)
			case SecChecker:
				r := Result{Name: c.String(), Type: "Sec", Data: v, IsMalicious: v.IsMalicious(), Err: err, ErrMsg: errMsg}
				res = append(res, r)
			}

		}(chk)
	}
	wg.Wait()

	return res
}

// JSON marshals results into JSON format.
func JSON(results []Result) ([]byte, error) {
	type Overview struct {
		Infos []struct {
			Name string
			Info string
		}
		ProbabilityMalicious float64
	}

	var o Overview

	for _, d := range results {
		i := struct {
			Name string
			Info string
		}{d.Name, d.Info}
		o.Infos = append(o.Infos, i)
	}
	o.ProbabilityMalicious = probabilityMalicious(results)

	j := struct {
		Overview
		Results []Result
	}{
		Overview: o,
		Results:  results,
	}

	return json.Marshal(&j)
}

func probabilityMalicious(results []Result) float64 {
	var malicious, totalSec float64
	for _, r := range results {
		if r.Type == "Sec" || r.Type == "InfoSec" {
			totalSec++
			if r.IsMalicious {
				malicious++
			}
		}
	}
	return malicious / totalSec
}

func redactSecrets(s string) string {
	key := regexp.MustCompile(`(key|pass|password)=\w+`)
	return key.ReplaceAllString(s, "${1}=REDACTED")
}

func na(s string) string {
	if s == "" {
		return "n/a"
	}
	return s
}

func nonEmpty(strings ...string) []string {
	var ss []string
	for _, s := range strings {
		if s != "" {
			ss = append(ss, s)
		}
	}
	return ss
}
