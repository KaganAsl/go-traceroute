package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	domain := os.Args[1]
	ips, err := dnsLookUp(domain)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get IPs: %v\n", err)
		os.Exit(1)
	}
	for _, ip := range ips {
		fmt.Printf("%s. IN A %s\n", domain, ip.String())
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
