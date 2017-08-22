package speedtest

import (
	"net/http"
	"net/url"
)

type Option func(*Config)

func ServerListURL(u *url.URL) Option {
	return func(c *Config) {
		c.ServerListURL = u
	}
}

func PlatformConfigURL(u *url.URL) Option {
	return func(c *Config) {
		c.PlatformConfigURL = u
	}
}

func UserAgent(ua string) Option {
	return func(c *Config) {
		c.UserAgent = ua
	}
}

func UploadSizes(size ...int) Option {
	return func(c *Config) {
		c.UploadSizes = size
	}
}

func DownloadSizes(size ...int) Option {
	return func(c *Config) {
		c.DownloadSizes = size
	}
}

func HTTPClient(h *http.Client) Option {
	return func(c *Config) {
		c.HTTPClient = h
	}
}

func Log(l Logger) Option {
	return func(c *Config) {
		c.Logger = l
	}
}

func defaultConfig() *Config {
	client := http.DefaultClient
	client.Transport = &http.Transport{
		DisableCompression: true,
	}

	return &Config{
		ServerListURL: &url.URL{
			Scheme: "http",
			Host:   "c.speedtest.net",
			Path:   "/speedtest-servers-static.php",
		},
		PlatformConfigURL: &url.URL{
			Scheme: "http",
			Host:   "c.speedtest.net",
			Path:   "/speedtest-config.php",
		},
		UserAgent:     "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.21 Safari/537.36",
		UploadSizes:   []int{256 * KByte, 512 * KByte, MByte, int(1.5 * MByte), 2 * MByte},
		DownloadSizes: []int{350, 500, 750, 1000, 1500, 2000, 2500, 3000, 3500, 4000},
		HTTPClient:    client,
		Logger:        &voidLogger{},
	}
}

type Config struct {
	ServerListURL     *url.URL
	PlatformConfigURL *url.URL
	UserAgent         string
	UploadSizes       []int
	DownloadSizes     []int
	HTTPClient        *http.Client
	Logger            Logger
}
