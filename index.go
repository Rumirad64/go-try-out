package main

import (
	"fmt"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

func main() {
	//two goroutines, one for sending, one for listening
	go send()
	go listen()
	go pcapListen()

	// wait for the goroutines to finish
	var input string
	fmt.Scanln(&input)

}

func send() {
	fmt.Println("Sending...")
	// define the address and port to send the packet to
	addr, _ := net.ResolveUDPAddr("udp", "54.179.1.70:3306")

	// create a connection to the address
	conn, _ := net.DialUDP("udp", nil, addr)

	// send the packet with a string
	_, _ = conn.Write([]byte("Hello, World!"))

}

func listen() {
	fmt.Println("PCAP Listening...")
	// listen for ICMP error replies
	packetConn, _ := net.ListenPacket("udp4", "0.0.0.0:15000")
	var buffer [1024]byte
	n, _, _ := packetConn.ReadFrom(buffer[:])
	fmt.Println("ICMP Error:", n)

}

func pcapListen() {
	//list all devices
	devices, err := pcap.FindAllDevs()
	if err != nil {
		panic(err)
	}
	fmt.Println("Devices:")
	/* 	for _, device := range devices {
		fmt.Println("Name: ", device.Name)
		fmt.Println("Description: ", device.Description)
		fmt.Println("Devices addresses: ", device.Description)
		for _, address := range device.Addresses {
			fmt.Println("- IP address: ", address.IP)
			fmt.Println("- Subnet mask: ", address.Netmask)
		}
	} */

	fmt.Println("Listening...")
	var device string = devices[5].Name
	fmt.Println("Device: ", device)
	if handle, err := pcap.OpenLive(device, 1600, true, pcap.BlockForever); err != nil {
		panic(err)
	} else if err := handle.SetBPFFilter("tcp and port 80"); err != nil { // optional
		panic(err)
	} else {
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		fmt.Println("Packet source: ", packetSource)
		for packet := range packetSource.Packets() {
			// print packet information
			fmt.Println("Packet: ", packet)
			fmt.Println(packet)
		}
	}

}
