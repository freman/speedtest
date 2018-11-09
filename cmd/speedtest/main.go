package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/freman/speedtest"
)

func resolveBindAddr(addr string) string {
	_, port, err := net.SplitHostPort(addr)
	if err != nil && err.(*net.AddrError).Err == "missing port in address" {
		return addr + ":0"
	}
	if port == "" {
		return addr + "0"
	}
	return addr
}

func main() {
	bindAddr := flag.String("bind", os.Getenv("SPEEDTEST_BIND"), "Address to bind to")

	flag.Parse()

	var options []speedtest.Option

	if *bindAddr != "" {
		localTCPAddr, err := net.ResolveTCPAddr("tcp", resolveBindAddr(*bindAddr))
		if err != nil {
			fmt.Fprintln(os.Stderr, "Unable to resolve bind address:", err)
			os.Exit(1)
		}
		options = append(options, speedtest.HTTPClient(&http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					LocalAddr: localTCPAddr,
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
					DualStack: true,
				}).DialContext,
				MaxIdleConns:          100,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			},
		}))
	}

	fmt.Println("Speetest Demo")
	c := speedtest.NewClient(options...)
	fmt.Println("Getting platform config")
	pc, err := c.GetPlatformConfig()
	fmt.Println("Testing from:", pc.Client.IP)
	fmt.Println("Getting server list")
	sl, err := c.GetServerList()
	if err != nil {
		panic(err)
	}

	fmt.Println("Finding the closest, fastest server")
	fastest := sl.Fastest(5)

	for _, s := range fastest {
		fmt.Println(s)
	}

	server := fastest[0]

	fmt.Printf("Found %v with %v latency\n", server, server.TestLatency())

	fmt.Println("Downloading...")
	fmt.Printf("  - %0.2fmbit/s\n", server.TestDownload())

	fmt.Println("Uploading...")
	fmt.Printf("  - %0.2fmbit/s\n", server.TestUpload())

}
