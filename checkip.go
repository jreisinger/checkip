// Checkip checks an IP address using various public services. An IP address is
// checked by running one or more Checkers. There are two kinds of Checkers. An
// InfoChecker just gathers some useful information about the IP address. A
// SecChecker says whether the IP address is considered malicious or not.
package checkip

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"regexp"
	"sync"

	"github.com/logrusorgru/aurora"
)

// Checker runs a check of an IP address.
type Checker interface {
	Check(ip net.IP) error
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

// Run runs checkers concurrently checking the ipaddr.
func Run(checkers []Checker, ipaddr net.IP) Result {
	var res Result

	var wg sync.WaitGroup
	for _, c := range checkers {
		wg.Add(1)
		go func(c Checker) {
			defer wg.Done()
			if err := c.Check(ipaddr); err != nil {
				res.Errors = append(res.Errors, redactSecrets(err.Error()))
			}
		}(c)
	}
	wg.Wait()

	var total, malicious int
	for _, c := range checkers {
		switch ip := c.(type) {
		case InfoChecker:
			res.Infos = append(res.Infos, ip.Info())
		case SecChecker:
			total++
			if ip.IsMalicious() {
				malicious++
			}
		}

	}
	res.ProbabilityMalicious = float64(malicious) / float64(total)

	return res
}

func redactSecrets(s string) string {
	key := regexp.MustCompile(`(key|pass|password)=\w+`)
	return key.ReplaceAllString(s, "${1}=REDACTED")
}

// Result holds the result of running a check.
type Result struct {
	Infos                []string
	ProbabilityMalicious float64
	Errors               []string
}

func (res Result) Print() error {
	for _, i := range res.Infos {
		fmt.Println(i)
	}

	var msg string

	switch {
	case res.ProbabilityMalicious < 0.15:
		msg = fmt.Sprint(aurora.Green("Malicious"))
	case res.ProbabilityMalicious < 0.50:
		msg = fmt.Sprint(aurora.Yellow("Malicious"))
	default:
		msg = fmt.Sprint(aurora.Red("Malicious"))
	}

	_, err := fmt.Printf("%s\t%.0f%%\n", msg, res.ProbabilityMalicious*100)
	return err
}

func (res Result) PrintJSON() error {
	enc := json.NewEncoder(os.Stdout)
	return enc.Encode(&res)
}
