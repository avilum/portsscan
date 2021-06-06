
package main

import (
	"context"
	"errors"
	"fmt"
	//"net"
	"net/http"
	"crypto/tls"
	"strings"
	"sync"
	"time"
	"golang.org/x/sync/semaphore"
)

type PortScanner struct {
	ip   string
	lock *semaphore.Weighted
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
	    Timeout: timeout,
	}

	resp, err := client.Head(fmt.Sprintf("http://%s", target))
	if err != nil {
		fmt.Println("(GO error): ", errors.Unwrap(fmt.Errorf("%w", err)))
		fmt.Println("Client: ", client)
		// fmt.Println("Go Error: ", err);

		if strings.Contains(err.Error(), "exceeded while awaiting") ||
		   strings.Contains(err.Error(), "ssl") ||
		   strings.Contains(err.Error(), "protocol") {
			fmt.Println(port,"<filtered (open)>")
			portsMapping[port] = true
			return
		} else {
			fmt.Println(port, "<closed>")
			portsMapping[port] = false
			return
		}
	}
	
	defer resp.Body.Close()
	
	//conn.Close()
	fmt.Println(port, "<open>")
	portsMapping[port] = true
	return
}

func (ps *PortScanner) Start(f, l int, timeout time.Duration,  portsMapping map[int]bool) {
	wg := sync.WaitGroup{}
	defer wg.Wait()
	for port := f; port <= l; port++ {
		ps.lock.Acquire(context.TODO(), 1)
		wg.Add(1)
		go func(port int) {
			defer ps.lock.Release(1)
			defer wg.Done()
			ScanPort(ps.ip, port, timeout, portsMapping)
		}(port)
	}
}

func main() {
	portsMapping := make(map[int]bool)
	ps := &PortScanner{
		ip:   "0.0.0.0",
		lock: semaphore.NewWeighted(5),
		portsMapping: portsMapping,
	}
//	ps.Start(1, 65535, 100*time.Millisecond)
	ps.Start(4999, 5001, 100 * time.Millisecond, portsMapping)
	fmt.Println("Finished. Ports Mapping:")
	fmt.Println(portsMapping)
}

