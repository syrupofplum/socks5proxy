package tests

import (
	src "../src"
	"testing"
)

func TestSocks5_Server(t *testing.T) {
	s := src.Server{}
	s.BindConfigs = make([]src.ServerBindConfig, 1, 1)
	s.BindConfigs[0].Network = "tcp"
	s.BindConfigs[0].Address = "127.0.0.1:1080"
	s.Listen()
}
