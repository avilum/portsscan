package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptrace"

	"crypto/tls"
	// "net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

type PortScanner struct {
	ip           string
	lock         *semaphore.Weighted
	portsMapping map[int]bool
}

func ScanPort(ip string, port int, timeout time.Duration, portsMapping map[int]bool) {
	target := fmt.Sprintf("%s:%d", ip, port)

	// WASI - supprt for TCP/UDP session - not supported within browsers.
	//conn, err := net.DialTimeout("tcp", target, timeout)

	// HTTP session - supported from browsers API
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := http.Client{
		Transport: tr,
		Timeout:   timeout,
	}
	target = fmt.Sprintf("http://%s", target)
	req, err := http.NewRequest("GET", target, nil)
	if err != nil {
		fmt.Print("Failed to initiate request ", err)
	}

	trace := &httptrace.ClientTrace{
		DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
			fmt.Printf("DNS Info: %+v\n", dnsInfo)
		},
		GotConn: func(connInfo httptrace.GotConnInfo) {
			fmt.Printf("Got Conn: %+v\n", connInfo)
		},
		GotFirstResponseByte: func() {
			fmt.Printf("Got first byte!")
		},
	}

	req.Header.Add("js.fetch:mode", "no-cors")
	// req.Header.Add("Access-Control-Allow-Origin", "0.0.0.0")
	// req.Header.Add("Access-Control-Allow-Credentials", "true")
	// req.Header.Add("Access-Control-Allow-Methods", "GET, PUT, POST, HEAD, TRACE, DELETE, PATCH, COPY, HEAD, LINK, OPTIONS")
	// Access-Control-Request-Method: POST
	//
	//resp, err := client.Get(target)
	fmt.Println("(GO request): ", fmt.Sprintf("%+v", req))
	//resp, err := client.Do(req)

	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	//if _, err := http.DefaultTransport.RoundTrip(req); err != nil {
	if _, err := client.Do(req); err != nil {
		fmt.Println(err)
		fmt.Println("(GO error): ", err.Error())
		if strings.Contains(strings.ToLower(err.Error()), "exceeded while awaiting") ||
			strings.Contains(strings.ToLower(err.Error()), "ssl") ||
			strings.Contains(strings.ToLower(err.Error()), "cors") ||
			strings.Contains(strings.ToLower(err.Error()), "invalid") ||
			strings.Contains(strings.ToLower(err.Error()), "protocol") {
			fmt.Println(port, "<filtered (open)>")
			portsMapping[port] = true
			return
		} else {
			fmt.Println(port, "<closed>")
			portsMapping[port] = false
			return
		}
	}

	// defer resp.Body.Close()
	fmt.Println(port, "<open>")
	portsMapping[port] = true
	return
}

func (ps *PortScanner) Start(f, l int, timeout time.Duration, portsMapping map[int]bool) {
	wg := sync.WaitGroup{}
	for port := f; port <= l; port++ {
		ps.lock.Acquire(context.TODO(), 1)
		wg.Add(1)
		go func(port int) {
			defer ps.lock.Release(1)
			defer wg.Done()
			ScanPort(ps.ip, port, timeout, portsMapping)
		}(port)
	}
	time.Sleep(5 * time.Second)
	wg.Wait()
}

func main() {
	portsMapping := make(map[int]bool)
	ps := &PortScanner{
		ip:           "0.0.0.0",
		lock:         semaphore.NewWeighted(5),
		portsMapping: portsMapping,
	}
	//	ps.Start(1, 65535, 100*time.Millisecond)
	ps.Start(4999, 5002, 10*time.Millisecond, portsMapping)
	fmt.Println("Finished. Ports Mapping:")
	fmt.Println(portsMapping)
}
