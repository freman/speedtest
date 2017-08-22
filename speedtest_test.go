package speedtest_test

import (
	"testing"

	"github.com/freman/speedtest"
)

func TestNewClient(t *testing.T) {
	if speedtest.NewClient() == nil {
		t.Error("Expected client got nil")
	}
}

func TestGetClientConfig(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping online tests")
	}
	c := speedtest.NewClient()
	pc, err := c.GetPlatformConfig()
	if err != nil {
		t.Error("unexpected error", err)
	}
	if pc == nil {
		t.Error("expected pc to not be nil")
	}
}

func TestGetServerList(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping online tests")
	}
	c := speedtest.NewClient()
	sl, err := c.GetServerList()
	if err != nil {
		t.Error("unexpected error", err)
	}
	if sl == nil {
		t.Error("expected sl to not be nil")
	}
	if len(sl) == 0 {
		t.Error("expected list to be != 0")
	}
}

func TestGetServerLatency(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping online tests")
	}
	c := speedtest.NewClient()
	sl, err := c.GetServerList()
	if err != nil {
		t.Error("unexpected error", err)
	}
	s := sl[0]

	latency := s.TestLatency()
	if latency == 0 {
		t.Error("expected latencyh to be != 0")
	}
}

func TestDownload(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping online tests")
	}
	c := speedtest.NewClient()
	sl, err := c.GetServerList()
	if err != nil {
		t.Fatal(err)
	}
	fastest := sl.Fastest(3)
	s := fastest[0]
	if s.TestDownload() == 0 {
		t.Error("expected a value greater than 0")
	}
}

func TestUpload(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping online tests")
	}
	c := speedtest.NewClient()
	sl, err := c.GetServerList()
	if err != nil {
		t.Fatal(err)
	}

	fastest := sl.Fastest(3)
	s := fastest[0]

	if s.TestUpload() == 0 {
		t.Error("expected a value greater than 0")
	}
}
