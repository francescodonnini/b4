package shared

import (
	"log"
	"net"
)

func GetAddress() Node {
	ip := getIp()
	settings := NewSettings()
	port, ok := settings.GetInt("PORT")
	if !ok {
		log.Fatalf("cannot read PORT\n")
	}
	return Node{
		Ip:   ip.String(),
		Port: port,
	}

}

func getIp() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = conn.Close() }()
	return conn.LocalAddr().(*net.UDPAddr).IP
}

func GetRegistryAddress() Node {
	settings := NewSettings()
	addr, _ := settings.GetString("REGISTRY")
	ip := net.ParseIP(addr)
	if ip == nil {
		log.Fatalln("invalid ip address of registry")
	}
	port, ok := settings.GetInt("REGISTRY_PORT")
	if !ok {
		log.Fatalf("cannot read port of registry\n")
	}
	return Node{
		Ip:   ip.String(),
		Port: port,
	}
}
