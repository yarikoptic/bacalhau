package model

import (
	"log"
	"sort"
	"sync"
	"time"
)

type ClusterPeers map[string][]string

type ClusterMapNode struct {
	ID    string `json:"id"`
	Group int    `json:"group"`
}

type ClusterMapLink struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

type ClusterMapResult struct {
	Nodes []ClusterMapNode `json:"nodes"`
	Links []ClusterMapLink `json:"links"`
}

type ClusterMap struct {
	Cluster *Cluster
	data    ClusterMapResult
	mutex   sync.Mutex
}

func NewClusterMap(cluster *Cluster) *ClusterMap {
	return &ClusterMap{
		Cluster: cluster,
		data:    ClusterMapResult{},
	}
}

func (clusterMap *ClusterMap) GetResults() ClusterMapResult {
	clusterMap.mutex.Lock()
	defer clusterMap.mutex.Unlock()
	return clusterMap.data
}

func (clusterMap *ClusterMap) Loop() {
	for {
		peers, err := clusterMap.LoadPeers()
		if err != nil {
			log.Println(err.Error())
			continue
		}
		func() {
			clusterMap.mutex.Lock()
			defer clusterMap.mutex.Unlock()
			clusterMap.data = clusterMap.ProcessPeers(peers)
		}()
		time.Sleep(1 * time.Second)
	}
}

func (clusterMap *ClusterMap) LoadPeers() (ClusterPeers, error) {
	newPeerMap := ClusterPeers{}
	for _, server := range clusterMap.Cluster.Servers {
		id, err := server.GetID()
		if err != nil {
			return nil, err
		}
		peers, err := server.GetPeers()
		if err != nil {
			return nil, err
		}
		newPeerMap[id] = peers
	}
	return newPeerMap, nil
}

func (clusterMap *ClusterMap) ProcessPeers(theMap ClusterPeers) ClusterMapResult {
	result := ClusterMapResult{}

	// keys of theMap
	keys := []string{}
	for k := range theMap {
		keys = append(keys, k)
	}
	// sort keys
	sort.Strings(keys)

	for _, node := range keys {
		links := theMap[node]
		result.Nodes = append(result.Nodes, ClusterMapNode{ID: node, Group: 0})
		for _, link := range links {
			result.Links = append(result.Links, ClusterMapLink{Source: node, Target: link})
		}
	}
	return result
}
