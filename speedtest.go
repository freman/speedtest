package speedtest

import (
	"encoding/xml"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	IByte = 1
	KByte = IByte * 1000
	MByte = KByte * 1000

	IBit = 1
	KBit = IBit * 1000
	MBit = KBit * 1000
)

type Client struct {
	config   *Config
	platform *PlatformConfig
}

func NewClient(options ...Option) *Client {
	c := &Client{
		config: defaultConfig(),
	}
	c.Options(options...)

	return c
}

func (c *Client) Options(options ...Option) {
	for _, o := range options {
		o(c.config)
	}
}

func (c *Client) doRequest(method string, u *url.URL, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, makeURLRandom(u).String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("User-Agent", c.config.UserAgent)

	return c.config.HTTPClient.Do(req)
}

func (c *Client) discardRequest(method string, u *url.URL, body io.Reader) (int64, error) {
	resp, err := c.doRequest(method, u, body)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	return io.Copy(ioutil.Discard, resp.Body)
}

func (c *Client) getXMLObject(u *url.URL, v interface{}) error {
	resp, err := c.doRequest(http.MethodGet, u, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	dec := xml.NewDecoder(resp.Body)
	if err := dec.Decode(v); err != nil {
		panic(err)
	}

	return nil
}

func (c *Client) GetPlatformConfig() (*PlatformConfig, error) {
	if c.platform == nil {
		c.platform = &PlatformConfig{}
		if err := c.getXMLObject(c.config.PlatformConfigURL, c.platform); err != nil {
			return nil, err
		}
	}
	return c.platform, nil
}

func (c *Client) GetServerList() (ServerList, error) {
	pc, err := c.GetPlatformConfig()
	if err != nil {
		return nil, err
	}

	list := &rawXMLServerList{}
	if err := c.getXMLObject(c.config.ServerListURL, list); err != nil {
		return nil, err
	}

	return list.unwrap(c, pc), nil
}
