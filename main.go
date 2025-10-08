package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/KaganAsl/go-traceroute/capture"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

var ttl int = 1
var packetCount int = 3
var maxTtl = 30
var timeout = 1

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

	icmpConn, err := capture.OpenICMPCapture()
	if err != nil {
		panic(err)
	}
	fmt.Println("ICMP connection opened.")
	defer icmpConn.Close()

	pconn := ipv4.NewPacketConn(conn)

	for range maxTtl {
		times := make([]float64, packetCount)
		source := ""
		var code int
		if err := pconn.SetTTL(ttl); err != nil {
			panic(err)
		}
		for i := range packetCount {
			msg := []byte("Hello")
			startTime := time.Now()
			if _, err := pconn.WriteTo(msg, nil, addr); err != nil {
				panic(err)
			}

			res := make([]byte, 1500)
			icmpConn.SetReadDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
			n, src, err := icmpConn.ReadFrom(res)
			if err, ok := err.(net.Error); ok && err.Timeout() {
				times[i] = -1
				continue
			}
			endTime := time.Since(startTime)
			endTimeMs := float64(endTime.Microseconds()) / 1000
			times[i] = endTimeMs
			if err != nil {
				panic(err)
			}
			res = res[:n]
			parsedMsg, err := icmp.ParseMessage(1, res)
			if err != nil {
				panic(err)
			}
			code = parsedMsg.Code
			if src != nil {
				source = src.String()
			}
		}
		result := source
		for _, t := range times {
			if t <= 0 {
				result += " *"
			} else {
				result += fmt.Sprintf(" %.3fms", t)
			}
		}
		fmt.Println(result)
		if code == 3 {
			return
		}
		ttl += 1
	}
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
