package main

import "fmt"

// IPPacket offers some functions working with IPv4 (!) IP packets
// packed for transmission wrapped into UDP
type IPPacket []byte

// IsIPv4 return true if packet is ipv4
func (p *IPPacket) IsIPv4() bool {
	return ((*p)[0] >> 4) == 4
}

// Dst returns [4]byte for destination of package
func (p *IPPacket) Dst() [4]byte {
	return [4]byte{(*p)[16], (*p)[17], (*p)[18], (*p)[19]}
}

// DstAddress returns [4]byte for destination of package
func (p *IPPacket) DstAddress() string {
	return fmt.Sprintf("%d.%d.%d.%d", (*p)[16], (*p)[17], (*p)[18], (*p)[19])
}

// DstRoute returns first 3 octets for destination of package
func (p *IPPacket) DstRoute() string {
	return fmt.Sprintf("%d.%d.%d", (*p)[16], (*p)[17], (*p)[18])
}

//tcp https [
//69 0 Version IHL DSCP	ECN
//0 64 Total Length
//0 0 Identification
//64 0 Flags Fragment Offset
//64 Time to Live
//6 Protocol = tcp
//85 100 Header Checksum
//192 168 50 1 Source address
//192 168 50 2 Destination address
// start tcp
//215 162 Source Port
//1 187 Destination Port
//204 93 36 138 Sequence Number
//0 0 0 0 Acknowledgement Number (meaningful when ACK bit set)
//176 Data Offset Reserved
//2 flags
//255 255 Window
//63 224 Checksum
//0 0 Urgent Pointer (meaningful when URG bit set)
//2 4 5 80 1 3 3 6 1 1 8 10 232 197 95 32 0 0 0 0 4 2 0 0]

//udp dns request ya.ru [
//69 0 Version IHL DSCP	ECN
//0 51 Total Length
//164 178 Identification
//0 0 Flags Fragment Offset
//64 Time to Live
//17 Protocol = udp
//240 179 Header Checksum
//192 168 50 1 Source address
//192 168 50 2 Destination address
// start udp
//243 143 0 53 0 31 134 243 198 177 1 0 0 1 0 0 0 0 0 0 2 121 97 2 114 117 0 0 1 0 1]

//ping  [
//69 0 Version IHL DSCP	ECN
//0 84 Total Length
//78 39 Identification
//0 0 Flags Fragment Offset
//64 Time to Live
//1 Protocot = icmp
//71 46 Header Checksum
//192 168 50 1 Source address
//192 168 50 2 Destination address
// start icmp
//8 0 9 237 210 81 0 2 103 32 3 121 0 11 198 23 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26 27 28 29 30 31 32 33 34 35 36 37 38 39 40 41 42 43 44 45 46 47 48 49 50 51 52 53 54 55]
