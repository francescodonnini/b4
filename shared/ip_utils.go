package shared

import (
	"log"
	"net"
)

func GetId() Node {
	ip := getIp()
	return Node{Ip: ip.String(), Port: 5050}
}

func GetEndpoint() Node {
	return Node{Ip: "10.0.0.253", Port: 5050}
}

func getIp() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = conn.Close() }()
	return conn.LocalAddr().(*net.UDPAddr).IP
}
