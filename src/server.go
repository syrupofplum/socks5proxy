package src

import (
	"fmt"
	"net"
	"strings"
)

type Server struct {
	BindConfigs []ServerBindConfig
}

func (s *Server) Listen() {
	errChan := make(chan error, 1)
	for _, bindConfig := range s.BindConfigs {
		bindConfig := bindConfig
		go func(errChan chan error) {
			switch strings.ToLower(bindConfig.Network) {
			case "tcp":
				ln, err := net.Listen(bindConfig.Network, bindConfig.Address)
				if nil != err {
					errChan <- err
					return
				}
				for true {
					conn, err := ln.Accept()
					if nil != err {
						errChan <- err
						return
					}
					go s.HandleTcpConn(conn, errChan)
				}
			case "udp":
				ln, err := net.Listen(bindConfig.Network, bindConfig.Address)
				if nil != err {
					errChan <- err
					return
				}
				for true {
					conn, err := ln.Accept()
					if nil != err {
						errChan <- err
						return
					}
					go s.HandleUdpConn(conn, errChan)
				}
			default:
				errChan <- fmt.Errorf("unknown network protocol")
				return
			}
		}(errChan)
	}
	for {
		select {
		case err := <-errChan:
			fmt.Println(err)
			break
		}
	}
}

func (s *Server) HandleUdpConn(conn net.Conn, errChan chan error) {

}
