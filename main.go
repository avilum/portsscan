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
        base := fmt.Sprintf("%s:%d", ip, port)

        placeHolder := js.Global().Get("document").Call("getElementById", "counter")
        placeHolder.Set("innerText", "Scanning "+ base + " / 10,000 ports")
        if port == 10000 {
                placeHolder.Set("innerText", "Scanned "+base+"; Waiting for the remaining responses...")
        }

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

        // Trying both http and https
        bind := fmt.Sprintf("http://%s", base)
        var req *http.Request
        var er error
        req, er = http.NewRequest("GET", bind, nil)
        if er != nil {
                fmt.Print("Failed to instantiate target over HTTP, trying HTTPS", er)
                bind = fmt.Sprintf("https://%s", base)
                req, er = http.NewRequest("GET", bind, nil)
                if er != nil {
                        fmt.Print("Failed to initiate target over HTTPS ", er)
                }
                // target = bind
                //      } else {
                //              target = bind
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
        // fmt.Println("(GO request): ", fmt.Sprintf("%+v", req))

        req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
        if _, err := client.Do(req); err != nil {
                // fmt.Println(err)
                // fmt.Println("(GO error): ", err.Error())
                // TODO: Get more exception strings for major browsers
                errString := strings.ToLower(err.Error())
                if strings.Contains(errString, "sent an invalid response") ||
                        strings.Contains(errString, "ERR_SSL_PROTOCOL_ERROR") ||
                        //                        errString, "exceeded while awaiting") ||
                        //                      strings.Contains(errString, "tls") ||
                        strings.Contains(errString, "ssl") ||
                        //                      strings.Contains(errString, "timeout") ||
                        strings.Contains(errString, "cors") ||
                        //                      strings.Contains(errString, "REFUSED") ||
                        strings.Contains(errString, "invalid") ||
                        //                      strings.Contains(errString, "https") ||
                        //                      strings.Contains(errString, "handshake") ||
                        strings.Contains(errString, "protocol") {
                        fmt.Println(port, "<filtered (open)>")
                        portsMapping[port] = true

                        // Append JS list element
//                      portString := strconv.Itoa(port)
//                      openPortsParagraph := js.Global().Get("document").Call("getElementByID", "openPort")
//                      openPortsParagraph.Set("innerText", portString)
//                      js.Global().Get("document").Get("body").Call("appendChild", openPortsParagraph)
                        return
                } else {
                        fmt.Println(port, "<closed>")
                        portsMapping[port] = false
                        return
                }
        }

        fmt.Println(port, "<open>")
        portsMapping[port] = true
//      portString := strconv.Itoa(port)
//      openPortsParagraph := js.Global().Get("document").Call("getElementByID", "openPort")
//      openPortsParagraph.Set("innerText", portString)
//      js.Global().Get("document").Get("body").Call("appendChild", openPortsParagraph)
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
        //time.Sleep(10 * time.Second)
        wg.Wait()
}

func main() {
        portsMapping := make(map[int]bool)
        ps := &PortScanner{
                ip:           "0.0.0.0",
                lock:         semaphore.NewWeighted(200),
                portsMapping: portsMapping,
        }

        // TODO: Enable port range input
        document := js.Global().Get("document")
        documentTitle := document.Call("createElement", "h1")
        documentTitle.Set("innerText", "TCP Port Scanner, Written in Go, Compiled to WebAssembly.")
        document.Get("body").Call("appendChild", documentTitle)

        placeHolder := document.Call("createElement", "h1")
        placeHolder.Set("innerText", "Scanning...")
        placeHolder.Set("id", "counter")
        document.Get("body").Call("appendChild", placeHolder)

        foundPort := document.Call("createElement", "h3")
        foundPort.Set("innerText", "Open: ")
        foundPort.Set("id", "openPort")
        document.Get("body").Call("appendChild", foundPort)

        //      scanned := document.Call("createElement", "h3")
        //      scanned.Set("innerText", "Done: ")
        //      scanned.Set("id", "scanned")
        //      document.Get("body").Call("appendChild", scanned)

        // Edit the ports,
        // ps.Start(1, 10000, 30000 * time.Millisecond, portsMapping)
        ps.Start(1, 10000, 500*time.Millisecond, portsMapping)
        fmt.Println("Finished. Ports Mapping:")

        var openPorts []string
        for k, v := range portsMapping {
                if v == true {
                        portString := strconv.Itoa(k)
                        openPorts = append(openPorts, portString)
                        document := js.Global().Get("document")
                        openPortsParagraph := document.Call("createElement", "li")
                        openPortsParagraph.Set("innerText", portString)
                        document.Get("body").Call("appendChild", openPortsParagraph)
                }
        }
        // fmt.Println("Open Ports", portsMapping)
        placeHolder.Set("innerText", "Open Ports:")
        fmt.Println("Scanned Ports: ", portsMapping)
        placeHolder.Set("innerText", portsMapping)
}
