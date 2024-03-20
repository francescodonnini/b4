package repl

import (
	"b4/sampling"
	"strings"
)

type Shell interface {
	Execute(line string) string
}

var Usage = []byte(`Type "help name" to find out more about the function name.
peers							returns the list of nodes "ip:port" that this node knows.
coord (--ip=address | --all)	returns the coordinate either of specified ip address or of all the coordinates known
dist --ip=address				returns the distance between this node and the one specified by the --ip option`)

type Impl struct {
	sampling sampling.PeerSamplingService
}

func (i *Impl) Execute(line string) []byte {
	// help
	// peers
	// coord [--ip=<ip> | --all]
	// dist --ip=<ip>
	fields := strings.Fields(line)
	if len(fields) < 1 {
		return Usage
	}
	switch fields[0] {
	case "help":
		return i.parseHelp(fields)
	case "peers":
		return i.parsePeers(fields)
	case "coord":
		return i.parseCoord(fields)
	case "dist":
		return i.parseDist(fields)
	default:
		return Usage
	}
}

func (i *Impl) parseDist(fields []string) []byte {
	panic("Not implemented yet!")
}

func (i *Impl) parseCoord(fields []string) []byte {
	panic("Not implemented yet!")
}

func (i *Impl) parseHelp(fields []string) []byte {
	panic("Not implemented yet!")
}

func (i *Impl) parsePeers(fields []string) []byte {
	panic("Not implemented yet!")
}
