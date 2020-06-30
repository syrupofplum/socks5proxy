package src

import "io"

func readProtocol(p protocolBase, rd io.Reader) error {
	n, err := rd.Read(p.getFullBuf())
	if nil != err {
		return err
	}
	p.setByteLength(n)
	return nil
}

func writeProtocol(p protocolBase, wr io.Writer) error {
	_, err := wr.Write(p.getBuf())
	if nil != err {
		return err
	}
	return nil
}

type protocolBase interface {
	setBufFromFields()
	setFieldsFromBuf()
	getFullBuf() []byte
	getBuf() []byte
	setByteLength(int)
	getByteLength() int
}

type protocolValidator interface {
	validate() bool
}

type protocolClientTCPMethods struct {
	_buf     [257]byte
	_n       int
	Ver      byte
	NMethods byte
	Methods  []byte
}

func (p *protocolClientTCPMethods) setBufFromFields() {
	p._buf[0] = p.Ver
	p._buf[1] = p.NMethods
	copy(p._buf[2:], p.Methods)
}

func (p *protocolClientTCPMethods) setFieldsFromBuf() {
	p.Ver = p._buf[0]
	p.NMethods = p._buf[1]
	p.Methods = p._buf[2 : 2+p.NMethods]
}

func (p *protocolClientTCPMethods) getFullBuf() []byte {
	return p._buf[:]
}

func (p *protocolClientTCPMethods) getBuf() []byte {
	return p._buf[:p._n]
}

func (p *protocolClientTCPMethods) getByteLength() int {
	return p._n
}

func (p *protocolClientTCPMethods) setByteLength(n int) {
	p._n = n
}

func (p *protocolClientTCPMethods) validate() bool {
	if p.Ver != byte(5) {
		return false
	}
	if p.NMethods < 1 || p.NMethods > 255 {
		return false
	}
	return true
}

type protocolClientTCPRequests struct {
	_buf    [4096]byte
	_n      int
	Ver     byte
	Cmd     byte
	Rsv     byte
	Atyp    byte
	DstAddr []byte
	DstPort []byte
}

func (p *protocolClientTCPRequests) setBufFromFields() {
	p._buf[0] = p.Ver
	p._buf[1] = p.Cmd
	p._buf[2] = p.Rsv
	p._buf[3] = p.Atyp
	copy(p._buf[4:], p.DstAddr)
	copy(p._buf[4+len(p.DstAddr):], p.DstPort)
}

func (p *protocolClientTCPRequests) setFieldsFromBuf() {
	p.Ver = p._buf[0]
	p.Cmd = p._buf[1]
	p.Rsv = p._buf[2]
	p.Atyp = p._buf[3]
	nDstAddr := p._buf[4]
	p.DstAddr = p._buf[4 : 4+1+nDstAddr]
	p.DstPort = p._buf[4+1+nDstAddr : 4+1+nDstAddr+2]
}

func (p *protocolClientTCPRequests) getFullBuf() []byte {
	return p._buf[:]
}

func (p *protocolClientTCPRequests) getBuf() []byte {
	return p._buf[:p._n]
}

func (p *protocolClientTCPRequests) getByteLength() int {
	return p._n
}

func (p *protocolClientTCPRequests) setByteLength(n int) {
	p._n = n
}

type protocolServerTCPMethods struct {
	_buf   [2]byte
	_n     int
	Ver    byte
	Method byte
}

func (p *protocolServerTCPMethods) setBufFromFields() {
	p._buf[0] = p.Ver
	p._buf[1] = p.Method
}

func (p *protocolServerTCPMethods) setFieldsFromBuf() {
	p.Ver = p._buf[0]
	p.Method = p._buf[1]
}

func (p *protocolServerTCPMethods) getFullBuf() []byte {
	return p._buf[:]
}

func (p *protocolServerTCPMethods) getBuf() []byte {
	return p._buf[:p._n]
}

func (p *protocolServerTCPMethods) getByteLength() int {
	return p._n
}

func (p *protocolServerTCPMethods) setByteLength(n int) {
	p._n = n
}

type protocolServerTCPReplies struct {
	_buf     [4096]byte
	_n       int
	Ver      byte
	Rep      byte
	Rsv      byte
	Atyp     byte
	BindAddr []byte
	BindPort []byte
}

func (p *protocolServerTCPReplies) setBufFromFields() {
	p._buf[0] = p.Ver
	p._buf[1] = p.Rep
	p._buf[2] = p.Rsv
	p._buf[3] = p.Atyp
	copy(p._buf[4:], p.BindAddr)
	copy(p._buf[4+len(p.BindAddr):], p.BindPort)
}

func (p *protocolServerTCPReplies) setFieldsFromBuf() {
	p.Ver = p._buf[0]
	p.Rep = p._buf[1]
	p.Rsv = p._buf[2]
	p.Atyp = p._buf[3]
	nBindAddr := p._buf[4]
	p.BindAddr = p._buf[4 : 4+1+nBindAddr]
	p.BindPort = p._buf[4+1+nBindAddr : 4+1+nBindAddr+2]
}

func (p *protocolServerTCPReplies) getFullBuf() []byte {
	return p._buf[:]
}

func (p *protocolServerTCPReplies) getBuf() []byte {
	return p._buf[:p._n]
}

func (p *protocolServerTCPReplies) getByteLength() int {
	return p._n
}

func (p *protocolServerTCPReplies) setByteLength(n int) {
	p._n = n
}
