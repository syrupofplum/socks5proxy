package src

import (
	"bufio"
	"fmt"
	"net"
)

type ProxyTcp struct {
	proxyConfig ProxyConfig
	proxyConn   net.Conn
	rd          *bufio.Reader
	wr          *bufio.Writer
}

func (p *ProxyTcp) connect() error {
	fmt.Println(p.proxyConfig.Network)
	fmt.Println(p.proxyConfig.Address)
	var err error
	p.proxyConn, err = net.Dial(p.proxyConfig.Network, p.proxyConfig.Address)
	if nil != err {
		return err
	}
	p.rd = bufio.NewReader(p.proxyConn)
	p.wr = bufio.NewWriter(p.proxyConn)
	return nil
}

func (p *ProxyTcp) send(data []byte) (int, error) {
	n, err := p.wr.Write(data)
	if nil != err {
		return -1, err
	}
	err = p.wr.Flush()
	if nil != err {
		return -1, err
	}
	fmt.Printf("send proxy data length: %v\n", n)
	return n, nil
}

func (p *ProxyTcp) handleSend(bufChan chan []byte, errChan chan error) {
	for {
		buf := <-bufChan
		fmt.Printf("send proxy: %v\n", buf[:5])
		_, err := p.send(buf)
		if nil != err {
			errChan <- fmt.Errorf("send proxy fail, err: %v\n", err)
			break
		}
	}
}

func (p *ProxyTcp) recv(data []byte) (int, error) {
	n, err := p.rd.Read(data)
	if nil != err {
		return -1, err
	}
	fmt.Printf("recv proxy data length: %v\n", n)
	return n, nil
}

func (p *ProxyTcp) handleRecv(buf []byte, bufChan chan []byte, errChan chan error) {
	for {
		n, err := p.recv(buf)
		if nil != err {
			errChan <- fmt.Errorf("recv proxy fail, err: %v\n", err)
			break
		}
		bufChan <- buf[:n]
	}
}
