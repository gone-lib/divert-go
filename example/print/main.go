package main

import (
	"fmt"

	"github.com/gone-lib/divert-go"
	"golang.org/x/net/ipv4"
)

func checkPacket(handle *divert.Handle, packetChan <-chan *divert.Packet) {
	for packet := range packetChan {
		fmt.Println("Sniffed", packet.Addr.Sniffed())
		fmt.Println("Outbound", packet.Addr.Outbound())
		fmt.Println("Loopback", packet.Addr.Loopback())
		fmt.Println("Impostor", packet.Addr.Impostor())
		fmt.Println("IPv6", packet.Addr.IPv6())
		fmt.Println("IPChecksum", packet.Addr.IPChecksum())
		fmt.Println("TCPChecksum", packet.Addr.TCPChecksum())
		fmt.Println("UDPChecksum", packet.Addr.UDPChecksum())
		header, _ := ipv4.ParseHeader(packet.Content)
		fmt.Println(header)
		// divert.HelperCalcChecksum(packet, 0)
		handle.Send(packet.Content, packet.Addr)
	}
}

func main() {
	// var filter = "(outbound  and tcp.DstPort == 1800) or (inbound  and  tcp.SrcPort == 1800)"
	var filter = "!loopback"
	handle, err := divert.Open(filter, divert.LayerNetwork, divert.PriorityLowest, 0)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	handle.SetParam(divert.QueueLength, divert.QueueLengthMax)
	handle.SetParam(divert.QueueTime, divert.QueueTimeMax)

	packetChan, err := handle.Packets()
	if err != nil {
		panic(err)
	}
	//defer handle.Close()
	checkPacket(handle, packetChan)
}
