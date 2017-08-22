package speedtest

import (
	"encoding/xml"
	"errors"
	"sort"
	"time"

	geo "github.com/kellydunn/golang-geo"
)

type rawXMLServerList struct {
	XMLName xml.Name `xml:"settings"`
	Servers struct {
		XMLName xml.Name   `xml:"servers"`
		List    ServerList `xml:"server"`
	} `xml:"servers"`
}

type ServerList []Server
type Fastest ServerList

func (r rawXMLServerList) unwrap(c *Client, pc *PlatformConfig) ServerList {
	p := geo.NewPoint(pc.Client.Lat, pc.Client.Lon)
	b := r.Servers.List[:0]
	for _, s := range r.Servers.List {
		if !pc.Server.IgnoreIDs.exists(s.ID) {
			p2 := geo.NewPoint(s.Lat, s.Lon)
			s.client = c
			s.Distance = p.GreatCircleDistance(p2)
			b = append(b, s)
		}
	}

	sort.Slice(b, func(i int, j int) bool {
		return b[i].Distance < b[j].Distance
	})

	return b
}

func (l ServerList) ByID(id int) (*Server, error) {
	for _, s := range l {
		if s.ID == id {
			return &s, nil
		}
	}
	return nil, errors.New("server not found")
}

func (x ServerList) Fastest(n int) []Server {
	l := make(map[Server]time.Duration, n)
	topN := x[:n]
	sort.Slice(topN, func(i, j int) bool {
		if _, ok := l[topN[i]]; !ok {
			l[topN[i]] = x[i].TestLatency()
		}

		if _, ok := l[topN[j]]; !ok {
			l[topN[j]] = x[j].TestLatency()
		}

		return l[topN[i]] > l[topN[j]]
	})

	return []Server(topN[:])
}
