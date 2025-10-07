//go:build !windows

package capture

import (
	"fmt"
	"net"
)

func OpenICMPCapture() (net.PacketConn, error) {
	icmpConn, err := net.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return nil, fmt.Errorf("failed to open ICMP socket: %w", err)
	}
	return icmpConn, nil
}
