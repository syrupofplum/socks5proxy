package src

func GetByteLength(p ProtocolBase) uint {
	return p.GetByteLength()
}

type ProtocolBase interface {
	GetDepth() uint
	GetByteLength() uint
}

type ProtocolValidator interface {
	Validate() bool
}

type ClientTCPMethods struct {
	VER      byte
	NMETHODS byte
	METHODS  []byte
}

func (c *ClientTCPMethods) GetDepth() uint {
	return 0
}

func (c *ClientTCPMethods) GetByteLength() uint {
	return uint(2 + c.NMETHODS)
}

func (c *ClientTCPMethods) Validate() bool {
	if c.VER != byte(5) {
		return false
	}
	if c.NMETHODS < 1 || c.NMETHODS > 255 {
		return false
	}
	return true
}

type ClientRemoteHeader struct {
	VER     byte
	CMD     byte
	RSV     byte
	ATYP    byte
	DSTADDR []byte
	DSTPORT []byte
}

func (c *ClientRemoteHeader) GetDepth() uint {
	return 100
}

type ClientRemoteDetail struct {
}

func (c *ClientRemoteDetail) GetDepth() uint {
	return 3
}

type ServerMethods struct {
	VER    byte
	METHOD byte
}

func (s *ServerMethods) GetDepth() uint {
	return 0
}
