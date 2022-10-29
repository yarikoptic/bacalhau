package model

import (
	"fmt"
	"strconv"
)

type Cluster struct {
	Servers []*Server
}

func NewCluster(args []string) (*Cluster, error) {
	// there is not point in an empty cluster!
	if len(args) <= 0 {
		return nil, fmt.Errorf("need arguments >= 3")
	}
	// is len(args) divisible by 3
	if len(args)%3 != 0 {
		return nil, fmt.Errorf("need arguments 3 at a time, e.g. " +
			"10.0.0.1 10000 10099 10.0.0.2 10000 10099 10.0.0.3 10000 10099")
	}
	numServers := len(args) / 3
	servers := []*Server{}
	for i := 0; i < numServers; i++ {
		address := args[i*3]
		start, err := strconv.Atoi(args[i*3+1])
		if err != nil {
			return nil, fmt.Errorf("can't interpret start port %s as uint: %s", args[i+1], err)
		}
		end, err := strconv.Atoi(args[i*3+2])
		if err != nil {
			return nil, fmt.Errorf("can't interpret end port %s as uint: %s", args[i+2], err)
		}
		if end < start {
			return nil, fmt.Errorf("end port (%d) must be >= start port (%d)", end, start)
		}

		for port := start; port <= end; port++ {
			server, err := NewServer(address, port)
			if err != nil {
				return nil, err
			}
			servers = append(servers, server)
		}
	}
	return &Cluster{
		Servers: servers,
	}, nil
}
