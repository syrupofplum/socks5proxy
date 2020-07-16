package src

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strconv"
)

func (s *Server) handleSend(wr *bufio.Writer, bufChan chan []byte, errChan chan error) {
	for {
		n, err := wr.Write(<-bufChan)
		if nil != err {
			errChan <- fmt.Errorf("send client data fail, err: %v\n", err)
			break
		}
		err = wr.Flush()
		if nil != err {
			errChan <- fmt.Errorf("send client data fail, err: %v\n", err)
			break
		}
		fmt.Printf("send client data length: %v\n", n)
	}
}

func (s *Server) HandleTcpConn(conn net.Conn, errChan chan error) {
	fmt.Println("in handle conn")
	defer conn.Close()

	var pClientTCPMethodsIns protocolClientTCPMethods
	var pServerTCPMethodsIns protocolServerTCPMethods
	var pClientTCPRequestsIns protocolClientTCPRequests
	var pServerTCPRepliesIns protocolServerTCPReplies

	//var clientBuf [40960]byte
	//var proxyBuf [40960]byte
	//var clientBufChan chan []byte = make(chan []byte, 1)
	//var proxyBufChan chan []byte = make(chan []byte, 1)

	rd := bufio.NewReader(conn)
	wr := bufio.NewWriter(conn)
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
	//go proxyTcp.handleSend(clientBufChan, errChan)
	//go proxyTcp.handleRecv(proxyBuf[:], proxyBufChan, errChan)

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

	//go s.handleSend(wr, proxyBufChan, errChan)
	//for {
	//	n, err := rd.Read(clientBuf[:])
	//	if nil != err {
	//		errChan <- fmt.Errorf("receive client data buffer fail, err: %v", err)
	//		return
	//	}
	//	fmt.Printf("receive client data buffer length: %v, %v\n", n, clientBuf[:5])
	//	clientBufChan <- clientBuf[:n]
	//}
	go io.Copy(conn, proxyTcp.proxyConn)
	io.Copy(proxyTcp.proxyConn, conn)
}
