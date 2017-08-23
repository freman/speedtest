package speedtest

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	gct "github.com/freman/go-commontypes"
)

type Server struct {
	XMLName  xml.Name `xml:"server"`
	URL      gct.URL  `xml:"url,attr"`
	Lat      float64  `xml:"lat,attr"`
	Lon      float64  `xml:"lon,attr"`
	Name     string   `xml:"name,attr"`
	Country  string   `xml:"country,attr"`
	CC       string   `xml:"cc,attr"`
	Sponsor  string   `xml:"sponsor,attr"`
	ID       int      `xml:"id,attr"`
	Distance float64
	Latency  time.Duration
	client   *Client
}

func (x Server) String() string {
	return fmt.Sprintf("%d: %s in %s by %s (%0.2fkm)", x.ID, x.Name, x.Country, x.Sponsor, x.Distance)
}

func (x *Server) LatencyURL() *url.URL {
	return x.URL.ResolveReference(&url.URL{Path: "latency.txt"})
}

func (x *Server) DownloadURL(size int) *url.URL {
	return x.URL.ResolveReference(&url.URL{Path: fmt.Sprintf("random%[1]dx%[1]d.jpg", size)})
}

func (x *Server) TestLatency() (latency time.Duration) {
	for i := 0; i < 3; i++ {
		start := time.Now()

		_, err := x.client.discardRequest(context.Background(), http.MethodGet, x.LatencyURL(), nil)
		if err != nil {
			if i == 0 {
				x.client.config.Logger.Warnf("failed to poll server %v due to %v", x, err)
				return time.Hour
			}
			latency += time.Hour
			continue
		}

		latency += time.Now().Sub(start)
	}
	latency = latency / time.Duration(6)
	return
}

func (x Server) TestDownload() (speed float64) {
	c := x.client
	cc := c.config
	pc, _ := c.GetPlatformConfig()

	threads := pc.Server.ThreadCount * 2
	requests := len(cc.DownloadSizes) * pc.Download.ThreadsPerURL

	cc.Logger.Debugf("launching %d download threads for %d requests", threads, requests)

	inch := make(chan *url.URL, threads)
	defer close(inch)

	ouch := make(chan int64, requests)
	defer close(ouch)

	timeout := time.Duration(pc.Download.TestLength) * time.Second
	ctx, cf := context.WithDeadline(context.Background(), time.Now().Add(timeout))
	defer cf()

	for i := 0; i < threads; i++ {
		go func(thread int) {
			cc.Logger.Debugf("\t[%d] waiting", thread)
			for url := range inch {
				cc.Logger.Debugf("\t[%d] downloading %s", thread, url)
				size, err := c.discardRequest(ctx, http.MethodGet, url, nil)
				if err != nil && err != context.DeadlineExceeded {
					cc.Logger.Warnf("unable to download from server: %v", err)
				}
				ouch <- size
				cc.Logger.Debugf("\t[%d] downloaded %d", thread, size)
			}
			cc.Logger.Debugf("\t[%d] shutting down", thread)
		}(i)
	}

	cc.Logger.Debugf("launching queuer thread for %d requests", requests)

	start := time.Now()

	go func() {
		for _, size := range cc.DownloadSizes {
			for i := 0; i < pc.Download.ThreadsPerURL; i++ {
				inch <- x.DownloadURL(size)
			}
		}
	}()

	size := int64(0)
	for i := 0; i < requests; i++ {
		size += <-ouch
	}

	seconds := time.Now().Sub(start).Seconds()
	megabits := float64(size*8) / float64(MBit)
	speed = megabits / seconds

	if size > 100000 {
		pc.Upload.Threads = 8
	}

	return
}

func (x Server) TestUpload() (speed float64) {
	c := x.client
	cc := c.config
	pc, _ := c.GetPlatformConfig()

	uploadSizes := cc.UploadSizes[pc.Upload.Ratio-1:]
	requests := pc.Upload.MaxChunkCount
	threads := pc.Upload.Threads

	buf := randomBytes(uploadSizes[len(uploadSizes)-1])

	cc.Logger.Debugf("launching %d download threads for %d requests", threads, requests)

	inch := make(chan int, threads)
	defer close(inch)

	ouch := make(chan int64, requests)
	defer close(ouch)

	timeout := time.Duration(pc.Upload.TestLength) * time.Second
	ctx, cf := context.WithDeadline(context.Background(), time.Now().Add(timeout))
	defer cf()

	for i := 0; i < threads; i++ {
		go func(thread int) {
			cc.Logger.Debugf("\t[%d] waiting", thread)
			for size := range inch {
				cc.Logger.Debugf("\t[%d] uploading %d bytes", thread, size)
				reader := bytes.NewReader(buf[:size])
				_, err := c.discardRequest(ctx, http.MethodPost, x.URL.URL, reader)
				if err != nil && err != context.DeadlineExceeded {
					cc.Logger.Warnf("unable to upload to server: %v", err)
				}
				read, _ := reader.Seek(0, io.SeekCurrent)
				ouch <- read
				cc.Logger.Debugf("\t[%d] uploaded %d", thread, size)
			}
			cc.Logger.Debugf("\t[%d] shutting down", thread)
		}(i)
	}

	cc.Logger.Debugf("launching queuer thread for %d requests", requests)

	start := time.Now()
	go func() {
		c := 0
		for {
			for _, size := range uploadSizes {
				inch <- size
				c++
				if c >= requests {
					return
				}
			}
		}
	}()

	size := int64(0)
	for i := 0; i < requests; i++ {
		size += <-ouch
	}

	seconds := time.Now().Sub(start).Seconds()
	megabits := float64(size*8) / float64(MBit)
	speed = megabits / seconds
	return
}
