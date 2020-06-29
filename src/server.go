package src

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type ServerBindConfig struct {
	Network string
	Address string
}

type Server struct {
	BindConfigs []ServerBindConfig
}

func (s *Server) Listen() {
	errChan := make(chan error, 1)
	go func(errChan chan error) {
		for _, bindConfig := range s.BindConfigs {
			switch strings.ToLower(bindConfig.Network) {
			case "tcp":
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
				go s.HandleTcpConn(conn, errChan)
			case "udp":
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
				go s.HandleUdpConn(conn, errChan)
			default:
				errChan <- fmt.Errorf("unknown network protocol")
				return
			}
		}
	}(errChan)
	select {
	case err := <-errChan:
		fmt.Println(err)
		break
	}
}

func (s *Server) HandleTcpConn(conn net.Conn, errChan chan error) {
	fmt.Println("in handle conn")
	defer conn.Close()

	var methodsBuffer [257]byte
	var methodsByteLength uint
	var methodsRespBuffer [2]byte
	var requestBuffer [32]byte
	var requestByteLength uint
	var _ [4096]byte

	rd := bufio.NewReader(conn)
	wr := bufio.NewWriter(conn)
	var n int
	var err error

	n, err = rd.Read(methodsBuffer[:])
	if nil != err {
		errChan <- fmt.Errorf("receive client methods buffer error")
		return
	}
	methodsByteLength = uint(n)
	fmt.Printf("methods buffer length: %v\n", methodsByteLength)
	fmt.Printf("methods buffer: %v\n", methodsBuffer[:16])

	methodsRespBuffer[0] = 0x05
	methodsRespBuffer[1] = 0x00
	n, err = wr.Write(methodsRespBuffer[:])
	if nil != err {
		errChan <- fmt.Errorf("send client methods resp buffer error")
		return
	}
	err = wr.Flush()
	if nil != err {
		errChan <- fmt.Errorf("send client methods resp buffer flush error")
		return
	}
	fmt.Printf("methods resp buffer length: %v\n", n)
	fmt.Printf("methods resp buffer: %v\n", methodsRespBuffer)

	n, err = rd.Read(requestBuffer[:])
	if nil != err {
		errChan <- fmt.Errorf("receive client request buffer error")
		return
	}
	requestByteLength = uint(n)
	fmt.Printf("request buffer length: %v\n", requestByteLength)
	fmt.Printf("request buffer: %v\n", requestBuffer)
}

func (s *Server) HandleUdpConn(conn net.Conn, errChan chan error) {

}
