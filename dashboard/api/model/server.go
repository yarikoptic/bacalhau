package model

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/filecoin-project/bacalhau/pkg/model"
	"github.com/filecoin-project/bacalhau/pkg/publicapi"
)

type Server struct {
	Address string `json:"address"`
	Port    int    `json:"port"`
	Api     *publicapi.APIClient
}

type JobInfo struct {
	Job     model.Job               `json:"job"`
	State   model.JobState          `json:"state"`
	Events  []model.JobEvent        `json:"events"`
	Results []model.PublishedResult `json:"results"`
}

func NewServer(address string, port int) (*Server, error) {
	server := &Server{
		Address: address,
		Port:    port,
	}
	server.Api = publicapi.NewAPIClient(server.GetApiAddress(""))
	return server, nil
}

func (server *Server) GetApiAddress(path string) string {
	return fmt.Sprintf("http://%s:%d%s", server.Address, server.Port, path)
}

func (server *Server) GetJobInfo(ctx context.Context, id string) (*JobInfo, error) {
	info := &JobInfo{}

	job, _, err := server.Api.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	info.Job = *job
	id = job.ID

	errorChan := make(chan error, 1)
	doneChan := make(chan bool, 1)
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		events, err := server.Api.GetEvents(ctx, id)
		if err != nil {
			errorChan <- err
		}
		info.Events = events
		wg.Done()
	}()
	go func() {
		state, err := server.Api.GetJobState(ctx, id)
		if err != nil {
			errorChan <- err
		}
		info.State = state
		wg.Done()
	}()
	go func() {
		results, err := server.Api.GetResults(ctx, id)
		if err != nil {
			errorChan <- err
		}
		info.Results = results
		wg.Done()
	}()
	go func() {
		wg.Wait()
		doneChan <- true
	}()
	select {
	case <-doneChan:
		return info, nil
	case err := <-errorChan:
		return nil, err
	}
}

func (server *Server) GetID() (string, error) {
	return HttpGet[string](server.GetApiAddress("/id"))
}

func (server *Server) GetPeers() ([]string, error) {
	data, err := HttpGet[map[string][]string](server.GetApiAddress("/peers"))
	if err != nil {
		return nil, err
	}
	peers, ok := data["bacalhau-job-event"]
	if !ok {
		return nil, fmt.Errorf("could not extract peers from bacalhau-job-event key")
	}
	sort.Strings(peers)
	return peers, nil
}

func (server *Server) GetDebug() (publicapi.DebugResponse, error) {
	return HttpGet[publicapi.DebugResponse](server.GetApiAddress("/debug"))
}
