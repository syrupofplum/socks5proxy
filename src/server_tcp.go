package src

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strconv"
)

func (s *Server) handleSend(proxyConn net.Conn, clientConn net.Conn, errChan chan error, outerChan chan error) {
	var buf = make([]byte, 32*1024)
	for {
		_, err := io.CopyBuffer(proxyConn, clientConn, buf)
		if nil != err {
			outerChan <- err
			errChan <- err
			break
		}
	}
}

func (s *Server) handleRecv(clientConn net.Conn, proxyConn net.Conn, errChan chan error, outerChan chan error) {
	var buf = make([]byte, 32*1024)
	for {
		_, err := io.CopyBuffer(clientConn, proxyConn, buf)
		if nil != err {
			outerChan <- err
			errChan <- err
			break
		}
	}
}

func (s *Server) HandleTcpConn(clientConn net.Conn, errChan chan error) {
	fmt.Println("in handle conn")
	defer clientConn.Close()

	var pClientTCPMethodsIns protocolClientTCPMethods
	var pServerTCPMethodsIns protocolServerTCPMethods
	var pClientTCPRequestsIns protocolClientTCPRequests
	var pServerTCPRepliesIns protocolServerTCPReplies

	rd := bufio.NewReader(clientConn)
	wr := bufio.NewWriter(clientConn)
	var err error

	err = readProtocol(&pClientTCPMethodsIns, rd)
	if nil != err {
		errChan <- fmt.Errorf("receive client methods buffer error, err: %v", err)
		return
	}
	pClientTCPMethodsIns.setFieldsFromBuf()
	fmt.Printf("methods buffer length: %v\n", pClientTCPMethodsIns.getByteLength())
	fmt.Printf("methods buffer: %v\n", pClientTCPMethodsIns.getBuf())

	pServerTCPMethodsIns.Ver = 0x05
	pServerTCPMethodsIns.Method = 0x00
	pServerTCPMethodsIns.setByteLength(2)
	pServerTCPMethodsIns.setBufFromFields()
	err = writeProtocol(&pServerTCPMethodsIns, wr)
	if nil != err {
		errChan <- fmt.Errorf("send server methods resp buffer error, err: %v", err)
		return
	}
	err = wr.Flush()
	if nil != err {
		errChan <- fmt.Errorf("send server methods resp buffer flush error, err: %v", err)
		return
	}
	fmt.Printf("methods resp buffer length: %v\n", pServerTCPMethodsIns.getByteLength())
	fmt.Printf("methods resp buffer: %v\n", pServerTCPMethodsIns.getBuf())

	err = readProtocol(&pClientTCPRequestsIns, rd)
	if nil != err {
		errChan <- fmt.Errorf("receive client request buffer error, err: %v", err)
		return
	}
	pClientTCPRequestsIns.setFieldsFromBuf()
	fmt.Printf("request buffer length: %v\n", pClientTCPRequestsIns.getByteLength())
	fmt.Printf("request buffer: %v\n", pClientTCPRequestsIns.getBuf())

	var proxyTcp ProxyTcp
	proxyTcp.proxyConfig.Network = "tcp"
	proxyTcp.proxyConfig.Address = fmt.Sprintf("%v:%v", string(pClientTCPRequestsIns.DstAddr[1:]), strconv.Itoa(int(binary.BigEndian.Uint16(pClientTCPRequestsIns.DstPort))))
	err = proxyTcp.connect()
	if nil != err {
		errChan <- fmt.Errorf("establish proxy tcp conn error, err: %v", err)
		return
	}
	defer proxyTcp.proxyConn.Close()

	pServerTCPRepliesIns.Ver = 0x05
	pServerTCPRepliesIns.Rep = 0x00
	pServerTCPRepliesIns.Rsv = 0x00
	pServerTCPRepliesIns.Atyp = pClientTCPRequestsIns.Atyp
	copy(pServerTCPRepliesIns._buf[4:], pClientTCPRequestsIns.DstAddr)
	copy(pServerTCPRepliesIns._buf[4+len(pClientTCPRequestsIns.DstAddr):], pClientTCPRequestsIns.DstPort)
	pServerTCPRepliesIns.setByteLength(4 + len(pClientTCPRequestsIns.DstAddr) + 2)
	pServerTCPRepliesIns.setBufFromFields()
	err = writeProtocol(&pServerTCPRepliesIns, wr)
	if nil != err {
		errChan <- fmt.Errorf("send server replies resp buffer error, err: %v", err)
		return
	}
	err = wr.Flush()
	if nil != err {
		errChan <- fmt.Errorf("send server replies resp buffer flush error, err: %v", err)
		return
	}
	fmt.Printf("replies resp buffer length: %v\n", pServerTCPRepliesIns.getByteLength())
	fmt.Printf("replies resp buffer: %v\n", pServerTCPRepliesIns.getBuf())

	var outerChan = make(chan error, 2)
	defer close(outerChan)
	var outerChanCounter = 0
	go s.handleSend(proxyTcp.proxyConn, clientConn, errChan, outerChan)
	go s.handleRecv(clientConn, proxyTcp.proxyConn, errChan, outerChan)
	for true {
		select {
		case <-outerChan:
			outerChanCounter += 1
			if outerChanCounter >= 2 {
				fmt.Printf("exit from %v\n", proxyTcp.proxyConfig.Address)
				break
			}
		}
	}
}
