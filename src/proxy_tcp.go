package src

import (
	"bufio"
	"net"
)

type ProxyTcp struct {
	proxyConfig ProxyConfig
	rd *bufio.Reader
	wr *bufio.Writer
}

func (p *ProxyTcp) connect() error {
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
	return n, err
}

func (p *ProxyTcp) handleSend(data chan []byte, errChan chan error) {
}

func (p *ProxyTcp) recv(data []byte) (int, error) {
	n, err := p.rd.Read(data)
	if nil != err {
		return -1, err
	}
	return n, err
}
