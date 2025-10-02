package main

import (
	"fmt"
	"net"
	"os"

	"golang.org/x/net/ipv4"
)

// type lookup struct {
// 	name    string
// 	results [][]net.IP
// }

const PORT int = 33434

func main() {
	domain := os.Args[1]
	fmt.Printf("Searching: %s\n", domain)
	ips, err := dnsLookUp(domain)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get IPs: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("IP: %s\n", ips[0])
	addr := &net.UDPAddr{IP: ips[0], Port: PORT}

	conn, err := net.ListenUDP("udp4", nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("UDP connection opened.")
	defer conn.Close()

	icmpConn, err := net.ListenPacket("ip:icmp", "0.0.0.0")
	if err != nil {
		panic(err)
	}
	fmt.Println("ICMP connection opened.")
	defer icmpConn.Close()

	pconn := ipv4.NewPacketConn(conn)
	if err := pconn.SetTTL(1); err != nil {
		panic(err)
	}
	fmt.Println("Set ttl.")

	msg := []byte("Hello")
	if _, err := pconn.WriteTo(msg, nil, addr); err != nil {
		panic(err)
	}
	fmt.Println("Msg write.")

	res := make([]byte, 1500)
	n, src, err := icmpConn.ReadFrom(res)
	if err != nil {
		panic(err)
	}
	res = res[:n]
	fmt.Println("Msg read:")
	fmt.Println(string(res))
	fmt.Println(src.String())
}

func dnsLookUp(domain string) ([]net.IP, error) {
	var ips, err = net.LookupIP(domain)
	if err != nil {
		return nil, err
	}
	var ipv4s []net.IP
	for _, ip := range ips {
		if ipv4 := ip.To4(); ipv4 != nil {
			ipv4s = append(ipv4s, ipv4)
		}
	}
	return ipv4s, nil
}
