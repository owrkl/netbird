package forwarder

import (
	log "github.com/sirupsen/logrus"
	wgdevice "golang.zx2c4.com/wireguard/device"
	"gvisor.dev/gvisor/pkg/tcpip"
	"gvisor.dev/gvisor/pkg/tcpip/header"
	"gvisor.dev/gvisor/pkg/tcpip/stack"
)

// endpoint implements stack.LinkEndpoint and handles integration with the wireguard device
type endpoint struct {
	dispatcher stack.NetworkDispatcher
	device     *wgdevice.Device
	mtu        uint32
}

func (e *endpoint) Attach(dispatcher stack.NetworkDispatcher) {
	e.dispatcher = dispatcher
}

func (e *endpoint) IsAttached() bool {
	return e.dispatcher != nil
}

func (e *endpoint) MTU() uint32 {
	return e.mtu
}

func (e *endpoint) Capabilities() stack.LinkEndpointCapabilities {
	return stack.CapabilityNone
}

func (e *endpoint) MaxHeaderLength() uint16 {
	return 0
}

func (e *endpoint) LinkAddress() tcpip.LinkAddress {
	return ""
}

func (e *endpoint) WritePackets(pkts stack.PacketBufferList) (int, tcpip.Error) {
	var written int
	for _, pkt := range pkts.AsSlice() {
		netHeader := header.IPv4(pkt.NetworkHeader().View().AsSlice())

		data := stack.PayloadSince(pkt.NetworkHeader())
		if data == nil {
			continue
		}

		// Send the packet through WireGuard
		address := netHeader.DestinationAddress()

		// TODO: handle dest ip addresses outside our network
		err := e.device.CreateOutboundPacket(data.AsSlice(), address.AsSlice())
		if err != nil {
			log.Errorf("CreateOutboundPacket: %v", err)
			continue
		}
		written++
	}

	return written, nil
}

func (e *endpoint) Wait() {
}

func (e *endpoint) ARPHardwareType() header.ARPHardwareType {
	return header.ARPHardwareNone
}

func (e *endpoint) AddHeader(*stack.PacketBuffer) {
}

func (e *endpoint) ParseHeader(*stack.PacketBuffer) bool {
	return true
}
