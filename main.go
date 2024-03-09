package main

import (
	"log"
	"net"
)

func main() {
	ip := GetIp()
	log.Println(ip.String())
}

func GetIp() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	return conn.LocalAddr().(*net.UDPAddr).IP
}
