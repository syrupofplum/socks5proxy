package src

import (
	"fmt"
	"net"
)

type ServerBindConfig struct {
	Network string
	Address string
}

type Server struct {
	BindConfigs []ServerBindConfig
}

func (s *Server) Listen() {
	errChan := make(chan error)
	go func(errChan chan error) {
		for _, bindConfig := range s.BindConfigs {
			ln, err := net.Listen(bindConfig.Network, bindConfig.Address)
			if nil != err {
				errChan <- err
				return
			}
			conn, err := ln.Accept()
			if nil != err {
				errChan <- err
				return
			}
			go s.HandleConn(conn, errChan)
		}
	}(errChan)
	select {
	case err := <- errChan:
		fmt.Println(err)
		break
	}
}

func (s *Server) HandleConn(conn net.Conn, errChan chan error) {
	defer conn.Close()
	// todo
}
