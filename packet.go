package divert

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type Packet struct {
	Content []byte
	Addr    *Address
}

func (p *Packet) IPVersion() uint8 {
	return p.Content[0] >> 4
}

func (p *Packet) Decode(options gopacket.DecodeOptions) gopacket.Packet {
	switch p.IPVersion() {
	case 4:
		return gopacket.NewPacket(p.Content, layers.LayerTypeIPv4, options)
	case 6:
		return gopacket.NewPacket(p.Content, layers.LayerTypeIPv6, options)
	}

	return nil
}
