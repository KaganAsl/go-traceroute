//go:build windows

package capture

import (
	"fmt"
	"net"
)

// TODO: Require a function to implement icmp packet captures. Following code won't work in windows
func OpenICMPCapture() (net.PacketConn, error) {
	icmpConn, err := net.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return nil, fmt.Errorf("failed to open ICMP socket: %w", err)
	}
	return icmpConn, nil
}
