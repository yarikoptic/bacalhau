package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/filecoin-project/bacalhau/dashboard/api/model"
	"github.com/filecoin-project/bacalhau/pkg/system"
)

// serve local files on web server
// fs := http.FileServer(http.Dir("./static"))
// http.Handle("/", fs)

func main() {
	if err := system.InitConfig(); err != nil {
		log.Fatal(err)
	}

	cluster, err := model.NewCluster(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	clusterMap := model.NewClusterMap(cluster)
	go clusterMap.Loop()

	fmt.Printf("servers: %d\n", len(cluster.Servers))

	http.Handle("/api/map", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(w).Encode(clusterMap.GetResults())
		if err != nil {
			log.Print(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}))

	http.Handle("/api/jobs", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jobsReq, err := model.GetRequestBody[struct {
			IDFilter    string `json:"idFilter"`
			MaxJobs     int    `json:"maxJobs"`
			ReturnAll   bool   `json:"returnAll"`
			SortBy      string `json:"sortBy"`
			SortReverse bool   `json:"sortReverse"`
		}](w, r)
		if err != nil {
			return
		}

		results, err := cluster.Servers[0].Api.List(
			context.Background(),
			jobsReq.IDFilter,
			jobsReq.MaxJobs,
			jobsReq.ReturnAll,
			jobsReq.SortBy,
			jobsReq.SortReverse,
		)

		if err != nil {
			log.Print(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(results)
		if err != nil {
			log.Print(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}))

	http.Handle("/api/jobinfo", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jobReq, err := model.GetRequestBody[struct {
			ID string `json:"id"`
		}](w, r)
		if err != nil {
			return
		}

		info, err := cluster.Servers[0].GetJobInfo(context.Background(), jobReq.ID)
		if err != nil {
			log.Print(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(info)
		if err != nil {
			log.Print(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}))

	log.Print("Listening on :31337...")
	err = http.ListenAndServe(":31337", nil)
	if err != nil {
		log.Fatal(err)
	}
}
