package util

import (
	fmt "fmt"
	"net"
	"strings"
)

// Connect dials the given address and returns a net.Conn. The protoAddr argument should be prefixed with the protocol,
// eg. "tcp://127.0.0.1:8080" or "unix:///tmp/test.sock"
func Connect(protoAddr string) (net.Conn, error) {
	proto, address := ProtocolAndAddress(protoAddr)
	conn, err := net.Dial(proto, address)
	return conn, err
}

// ProtocolAndAddress splits an address into the protocol and address components.
// For instance, "tcp://127.0.0.1:8080" will be split into "tcp" and "127.0.0.1:8080".
// If the address has no protocol prefix, the default is "tcp".
func ProtocolAndAddress(listenAddr string) (string, string) {
	protocol, address := "tcp", listenAddr
	parts := strings.SplitN(address, "://", 2)
	if len(parts) == 2 {
		protocol, address = parts[0], parts[1]
	}
	return protocol, address
}

var localIpSegment = []byte{192, 168}

// 如设置为 172.31,直接用 172, 31 作为变量传入就行
// SetLocalIPSegment(172, 31)
func SetLocalIPSegment(segs ...byte) {
	localIpSegment = []byte{}
	for i := 0; i < len(segs) && i < 4; i++ {
		localIpSegment = append(localIpSegment, segs[i])
	}
}

func GetLocalIPSegment() [2]byte {
	var segs [2]byte
	for i := 0; i < len(localIpSegment) && i < 2; i++ {
		segs[i] = localIpSegment[i]
	}
	return segs
}

// 获取本机 ip
// 默认本地的 ip 段为 192.168，如果不是，调用此方法前先调用 SetLocalIPSegment 方法设置本地 ip 段
func LocalIP() (net.IP, error) {
	tables, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, t := range tables {
		addrs, err := t.Addrs()
		if err != nil {
			return nil, err
		}
		for _, a := range addrs {
			ipnet, ok := a.(*net.IPNet)
			if !ok {
				continue
			}
			if v4 := ipnet.IP.To4(); v4 != nil {
				var matchd = true
				for i := 0; i < len(localIpSegment); i++ {
					if v4[i] != localIpSegment[i] {
						matchd = false
					}
				}
				if matchd {
					return v4, nil
				}
			}
		}
	}
	return nil, fmt.Errorf("cannot find local IP address")
}
