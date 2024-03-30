package repl

import (
	"b4/gossip"
	"b4/shared"
	"encoding/json"
	"errors"
	"strings"
)

type Shell interface {
	Execute(line string) ([]byte, error)
}

var Usage = []byte(`Type "help name" to find out more about the function name.
peers							returns the list of nodes "ip:port" that this node knows.
coord (--ip=address | --all)	returns the coordinate either of specified ip address or of all the coordinates known
dist --ip=address				returns the distance between this node and the one specified by the --ip option`)

type Impl struct {
	store gossip.Store
}

func NewShell(store gossip.Store) Shell {
	return &Impl{store: store}
}

func (i *Impl) Execute(line string) ([]byte, error) {
	// help
	// peers
	// coord [--ip=<ip> | --all]
	// dist --ip=<ip>
	fields := strings.Fields(line)
	if len(fields) < 1 {
		return Usage, nil
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
		return Usage, nil
	}
}

func (i *Impl) parseDist(fields []string) ([]byte, error) {
	panic("Not implemented yet!")
}

func (i *Impl) parseCoord(fields []string) ([]byte, error) {
	coord, ok := i.store.Read(shared.Node{
		Ip:   fields[0],
		Port: 5050,
	})
	if !ok {
		return make([]byte, 0), nil
	}
	bytes, err := json.Marshal(coord)
	if err != nil {
		return make([]byte, 0), errors.New("cannot retrieve list of peers")
	}
	return bytes, nil
}

func (i *Impl) parseHelp(fields []string) ([]byte, error) {
	panic("Not implemented yet!")
}

func (i *Impl) parsePeers(_ []string) ([]byte, error) {
	view := i.store.Peers()
	bytes, err := json.Marshal(view)
	if err != nil {
		return make([]byte, 0), errors.New("cannot retrieve list of peers")
	}
	return bytes, nil
}
