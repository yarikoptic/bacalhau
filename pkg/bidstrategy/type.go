package bidstrategy

import (
	"context"

	"github.com/bacalhau-project/bacalhau/pkg/model"
)

type BidStrategyRequest struct {
	NodeID string
	Job    model.Job
}

type BidStrategyResponse struct {
	ShouldBid  bool   `json:"shouldBid"`
	ShouldWait bool   `json:"shouldWait"`
	Reason     string `json:"reason"`
}

func NewShouldBidResponse() BidStrategyResponse {
	return BidStrategyResponse{
		ShouldBid: true,
	}
}

type BidStrategy interface {
	ShouldBid(ctx context.Context, request BidStrategyRequest) (BidStrategyResponse, error)
	ShouldBidBasedOnUsage(ctx context.Context, request BidStrategyRequest, resourceUsage model.ResourceUsageData) (BidStrategyResponse, error)
}

// the JSON data we send to http or exec probes
// TODO: can we just use the BidStrategyRequest struct?
type JobSelectionPolicyProbeData struct {
	NodeID string     `json:"node_id"`
	JobID  string     `json:"job_id"`
	Spec   model.Spec `json:"spec"`
}

// Return JobSelectionPolicyProbeData for the given request
func getJobSelectionPolicyProbeData(request BidStrategyRequest) JobSelectionPolicyProbeData {
	return JobSelectionPolicyProbeData{
		NodeID: request.NodeID,
		JobID:  request.Job.Metadata.ID,
		Spec:   request.Job.Spec,
	}
}
