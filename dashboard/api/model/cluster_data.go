package model

import (
	"log"
	"sort"
	"sync"
	"time"

	"github.com/filecoin-project/bacalhau/pkg/publicapi"
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

type ClusterData struct {
	Cluster        *Cluster
	clusterMapData ClusterMapResult
	debugData      []publicapi.DebugResponse
	mutex          sync.Mutex
}

func NewClusterData(cluster *Cluster) *ClusterData {
	return &ClusterData{
		Cluster:        cluster,
		clusterMapData: ClusterMapResult{},
		debugData:      []publicapi.DebugResponse{},
	}
}

func (clusterMap *ClusterData) GetClusterMapData() ClusterMapResult {
	clusterMap.mutex.Lock()
	defer clusterMap.mutex.Unlock()
	return clusterMap.clusterMapData
}

func (clusterMap *ClusterData) GetDebugData() []publicapi.DebugResponse {
	clusterMap.mutex.Lock()
	defer clusterMap.mutex.Unlock()
	return clusterMap.debugData
}

func (clusterMap *ClusterData) Loop() {
	for {
		peers, err := clusterMap.LoadPeers()
		if err != nil {
			log.Println(err.Error())
			continue
		}
		debug, err := clusterMap.LoadDebug()
		if err != nil {
			log.Println(err.Error())
			continue
		}
		func() {
			clusterMap.mutex.Lock()
			defer clusterMap.mutex.Unlock()
			clusterMap.clusterMapData = clusterMap.ProcessPeers(peers)
			clusterMap.debugData = debug
		}()
		time.Sleep(1 * time.Second)
	}
}

func (clusterMap *ClusterData) LoadDebug() ([]publicapi.DebugResponse, error) {
	return clusterMap.Cluster.LoadDebugData()
}

func (clusterMap *ClusterData) LoadPeers() (ClusterPeers, error) {
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

func (clusterMap *ClusterData) ProcessPeers(theMap ClusterPeers) ClusterMapResult {
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
