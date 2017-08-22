package speedtest

import (
	"encoding/xml"
	"strconv"
	"strings"

	gct "github.com/freman/go-commontypes"
)

var emptyStruct struct{}

type idfilter map[int]struct{}

func (f idfilter) exists(id int) bool {
	_, exists := f[id]
	return exists
}

func (f *idfilter) UnmarshalXMLAttr(attr xml.Attr) error {
	(*f) = make(map[int]struct{})
	for _, v := range strings.Split(attr.Value, ",") {
		d, _ := strconv.Atoi(v)
		(*f)[d] = emptyStruct
	}

	return nil
}

type PlatformConfig struct {
	XMLName  xml.Name       `xml:"settings"`
	Client   ClientConfig   `xml:"client"`
	Server   ServerConfig   `xml:"server-config"`
	Upload   UploadConfig   `xml:"upload"`
	Download DownloadConfig `xml:"download"`
	Latency  LatencyConfig  `xml:"latency"`
}

type ClientConfig struct {
	IP  gct.IP  `xml:"ip,attr"`
	Lat float64 `xml:"lat,attr"`
	Lon float64 `xml:"lon,attr"`
	Isp string  `xml:"isp,attr"`
}

type ServerConfig struct {
	ThreadCount int      `xml:"threadcount,attr"`
	IgnoreIDs   idfilter `xml:"ignoreids,attr"`
}

type UploadConfig struct {
	TestLength    int    `xml:"testlength,attr"`
	Ratio         int    `xml:"ratio,attr"`
	InitialTest   string `xml:"initialtest,attr"`
	MinTestSize   string `xml:"mintestsize,attr"`
	Threads       int    `xml:"threads,attr"`
	MaxChunkSize  string `xml:"maxchunksize,attr"`
	MaxChunkCount int    `xml:"maxchunkcount,attr"`
	ThreadsPerURL int    `xml:"threadsperurl,attr"`
}

type DownloadConfig struct {
	TestLength    int    `xml:"testlength,attr"`
	InitialTest   string `xml:"initialtest,attr"`
	MinTestSize   string `xml:"mintestsize,attr"`
	ThreadsPerURL int    `xml:"threadsperurl,attr"`
}

type LatencyConfig struct {
	TestLength int `xml:"testlength,attr"`
	WaitTime   int `xml:"waittime,attr"`
	Timeout    int `xml:"timeout,attr"`
}
