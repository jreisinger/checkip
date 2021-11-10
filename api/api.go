package api

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	checkip "github.com/jreisinger/checkip/pkg"
)

type Checkers []checkip.Checker

func (c Checkers) Handler(w http.ResponseWriter, r *http.Request) {
	ip := strings.Split(r.URL.Path, "/")[3]
	netip := net.ParseIP(ip)
	if netip == nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "wrong IP address: %q", ip)
		return
	}
	results := checkip.Run(c, netip)
	b, err := checkip.JSON(results)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "marshalling JSON: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", b)
}
