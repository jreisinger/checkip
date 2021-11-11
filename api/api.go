package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/jreisinger/checkip/pkg/check"
	"github.com/jreisinger/checkip/pkg/checker"
)

func Serve(addr string, path string) {
	s := server{addr, path}
	http.HandleFunc(path, s.handle)
	log.Printf("starting API server at %s:%s", addr, path)
	log.Fatal(http.ListenAndServe(addr, nil))
}

type server struct {
	addr string
	path string
}

func (s server) handle(w http.ResponseWriter, r *http.Request) {
	ip := strings.TrimPrefix(r.URL.Path, s.path)
	ipaddr := net.ParseIP(ip)
	if ipaddr == nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "wrong IP address: %q", ip)
		return
	}

	results := check.Run(checker.DefaultCheckers, ipaddr)
	b, err := json.Marshal(results)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "marshalling JSON: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", b)
}
