package main

import (
	"fmt"

	"github.com/freman/speedtest"
)

func main() {
	fmt.Println("Speetest Demo")
	c := speedtest.NewClient()
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
	fmt.Printf("  - %0.2fmbit/s", server.TestDownload())

	fmt.Println("Uploading...")
	fmt.Printf("  - %0.2fmbit/s", server.TestUpload())

}
