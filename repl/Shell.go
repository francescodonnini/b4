package repl

import (
	"b4/gossip"
	"b4/shared"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"
	"time"
)

type Shell interface {
	Execute(line string) ([]byte, error)
}

var Usage = []byte(`Type "help name" to find out more about the function name.
peers							returns the list of nodes "ip:port" that this node knows.
coord (--ip=address | --all)	returns the coordinate either of specified ip address or of all the coordinates known
dist --ip=address				returns the distance between this node and the one specified by the --ip option
error expected-rtt				return node error`)

type Impl struct {
	id    shared.Node
	store gossip.Store
}

func NewShell(id shared.Node, store gossip.Store) Shell {
	return &Impl{id: id, store: store}
}

func (i *Impl) Execute(line string) ([]byte, error) {
	// help
	// peers
	// coord [--ip=<ip> | --all]
	// dist --ip=<ip>
	fields := strings.Fields(line)
	if len(fields) == 1 && fields[0] == "peers" {
		return i.parsePeers(fields[1:])
	}
	if len(fields) != 2 {
		return Usage, nil
	}
	switch fields[0] {
	case "help":
		return i.parseHelp(fields[1:])
	case "coord":
		return i.parseCoord(fields[1:])
	case "dist":
		return i.parseDist(fields[1:])
	case "error":
		return i.parseError(fields[1:])
	default:
		return Usage, nil
	}
}

func (i *Impl) parseDist(fields []string) ([]byte, error) {
	if strings.HasPrefix(fields[0], "--ip=") {
		_, ip, ok := strings.Cut(fields[0], "=")
		if ok {
			self, ok := i.store.Read(i.id)
			coord, ok := i.store.Read(shared.Node{
				Ip:   ip,
				Port: 5050,
			})
			if !ok {
				return make([]byte, 0), nil
			}
			dist := self.Coord.Sub(coord.Coord).Magnitude()
			bytes, err := json.Marshal(dist)
			if err != nil {
				return make([]byte, 0), err
			}
			return bytes, nil
		}
	}
	return make([]byte, 0), fmt.Errorf("invalid argument for \"dist\"")
}

func (i *Impl) parseCoord(fields []string) ([]byte, error) {
	if fields[0] == "--all" {
		bytes, err := json.Marshal(i.store.Items())
		if err != nil {
			return make([]byte, 0), err
		}
		return bytes, nil
	} else if strings.HasPrefix(fields[0], "--ip=") {
		_, ip, ok := strings.Cut(fields[0], "=")
		if ok {
			var node shared.Node
			if ip == "self" {
				node = i.id
			} else {
				node = shared.Node{
					Ip:   ip,
					Port: 5050,
				}
			}
			coord, ok := i.store.Read(node)
			if !ok {
				return make([]byte, 0), nil
			}
			bytes, err := json.Marshal(coord)
			if err != nil {
				return make([]byte, 0), err
			}
			return bytes, nil
		}
	}
	return make([]byte, 0), fmt.Errorf("invalid argument for \"coord\"")
}

func (i *Impl) parseHelp(_ []string) ([]byte, error) {
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

func (i *Impl) parseError(fields []string) ([]byte, error) {
	param, err := strconv.ParseInt(fields[0], 10, 32)
	if err != nil {
		return nil, errors.New("invalid rtt")
	}
	self, ok := i.store.Read(i.id)
	if !ok {
		return nil, errors.New("self coord does not exists")
	}
	peers := i.store.Items()
	if len(peers) <= 0 {
		return nil, errors.New("no other coordinates known")
	}
	e := make([]float64, 0)
	rtt := time.Duration(param) * time.Millisecond
	for _, it := range peers {
		if it.Owner == self.Owner {
			continue
		}
		dist := time.Duration(self.Coord.Sub(it.Coord).Magnitude()) * time.Second
		e = append(e, math.Abs(dist.Seconds()-rtt.Seconds()))
	}
	slices.Sort(e)
	bytes, err := json.Marshal(e[len(e)/2])
	if err != nil {
		return make([]byte, 0), errors.New("cannot retrieve node error")
	}
	return bytes, nil
}
