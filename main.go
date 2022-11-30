package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptrace"

	"crypto/tls"
	"strconv"
	"strings"
	"sync"
	"syscall/js"
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

	// TODO: WASI - add supprt for TCP/UDP session. only HTTP/HTTPS is supported in WASM JS today.
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
			fmt.Println("DNS Info: %+v\n", dnsInfo)
		},
		GotConn: func(connInfo httptrace.GotConnInfo) {
			fmt.Println("Got Conn: %+v\n", connInfo)
		},
		GotFirstResponseByte: func() {
			fmt.Println("Got first byte!")
		},
	}

	// IMPORTANT - enables better HTTP(S) discovery, because many browsers block CORS by default.
	req.Header.Add("js.fetch:mode", "no-cors")
	fmt.Println("(GO request): ", fmt.Sprintf("%+v", req))

	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	if _, err := client.Do(req); err != nil {
		fmt.Println(err)
		fmt.Println("(GO error): ", err.Error())
		// TODO: Get more exception strings for major browsers
		errString := strings.ToLower(err.Error())
		if strings.Contains(errString, "exceeded while awaiting") ||
			strings.Contains(errString, "ssl") ||
			strings.Contains(errString, "cors") ||
			strings.Contains(errString, "invalid") ||
			strings.Contains(errString, "protocol") {
			fmt.Println(port, "<filtered (open)>")
			portsMapping[port] = true
			return
		} else {
			fmt.Println(port, "<closed>")
			portsMapping[port] = false
			return
		}
	}

	fmt.Println(port, "<open>")
	portsMapping[port] = true
	return
}

func (ps *PortScanner) Start(f int, l int, timeout time.Duration, portsMapping map[int]bool) {
	wg := sync.WaitGroup{}
	for port := f; port <= l; port++ {
		// GO in WASM must be SYNC as of today
		ps.lock.Acquire(context.TODO(), 1)
		wg.Add(1)
		go func(port int) {
			defer ps.lock.Release(1)
			ScanPort(ps.ip, port, timeout, portsMapping)
			defer wg.Done()
		}(port)
	}

	time.Sleep(5 * time.Second)
	wg.Wait()
}

func main() {
	portsMapping := make(map[int]bool)
	ps := &PortScanner{
		ip:           "0.0.0.0",
		lock:         semaphore.NewWeighted(10),
		portsMapping: portsMapping,
	}

	// TODO: Enable port range input
	// ps.Start(1, 65535, 10*time.Millisecond, portsMapping)
	document := js.Global().Get("document")
	documentTitle := document.Call("createElement", "h1")
	documentTitle.Set("innerText", "WebAssembly TCP Port Scanner")
	document.Get("body").Call("appendChild", documentTitle)
	placeHolder := document.Call("createElement", "h3")
	placeHolder.Set("innerText", "Scanning...")
	document.Get("body").Call("appendChild", placeHolder)

	// Edit the ports,
	ps.Start(1, 10000, 1000*time.Millisecond, portsMapping)
	fmt.Println("Finished. Ports Mapping:")

	var openPorts []string
	for k, v := range portsMapping {
		if v == true {
			portString := strconv.Itoa(k)
			openPorts = append(openPorts, portString)
			openPortsParagraph := document.Call("createElement", "li")
			openPortsParagraph.Set("innerText", portString)
			document.Get("body").Call("appendChild", openPortsParagraph)
		}
	}
	fmt.Println("Scanned Ports: ", portsMapping)
	fmt.Println("Open Ports", portsMapping)
	placeHolder.Set("innerText", "Open Ports:")
}
