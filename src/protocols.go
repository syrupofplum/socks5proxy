package src

type ProtocolBase interface {
	GetDepth() uint8
}

type ProtocolValidator interface {
	Validate() bool
}

type ClientTCPMethodsHeader struct {
	VER      byte
	NMETHODS byte
}

func (c *ClientTCPMethodsHeader) GetDepth() uint8 {
	return 0
}

func (c *ClientTCPMethodsHeader) Validate() bool {
	if c.VER != byte(5) {
		return false
	}
	if c.NMETHODS < 1 || c.NMETHODS > 255 {
		return false
	}
	return true
}

type ClientTCPMethods struct {
	METHODS  []byte
}

func (c *ClientTCPMethods) GetDepth() uint8 {
	return 1
}

type ClientRemoteHeader struct {
	VER byte
	CMD byte
	RSV byte
	ATYP byte
}

func (c *ClientRemoteHeader) GetDepth() uint8 {
	return 2
}

type ClientRemoteDetail struct {
	DSTADDR []byte
	DSTPORT []byte
}

func (c *ClientRemoteDetail) GetDepth() uint8 {
	return 3
}

type ServerMethods struct {
	VER    byte
	METHOD byte
}

func (s *ServerMethods) GetDepth() uint8 {
	return 0
}
