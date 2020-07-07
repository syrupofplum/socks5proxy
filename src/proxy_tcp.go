package src

import (
	"bufio"
	"fmt"
	"net"
)

type ProxyTcp struct {
	proxyConfig ProxyConfig
	rd          *bufio.Reader
	wr          *bufio.Writer
}

func (p *ProxyTcp) connect() error {
	fmt.Println(p.proxyConfig.Network)
	fmt.Println(p.proxyConfig.Address)
	proxyConn, err := net.Dial(p.proxyConfig.Network, p.proxyConfig.Address)
	if nil != err {
		return err
	}
	p.rd = bufio.NewReader(proxyConn)
	p.wr = bufio.NewWriter(proxyConn)
	return nil
}

func (p *ProxyTcp) send(data []byte) (int, error) {
	n, err := p.wr.Write(data)
	if nil != err {
		return -1, err
	}
	fmt.Printf("send proxy data length: %v\n", n)
	return n, nil
}

func (p *ProxyTcp) handleSend(bufChan chan []byte, errChan chan error) {
	for {
		_, err := p.send(<-bufChan)
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
