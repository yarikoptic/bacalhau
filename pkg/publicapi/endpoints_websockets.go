package publicapi

import (
	"net/http"

	"github.com/filecoin-project/bacalhau/pkg/system"
	"github.com/filecoin-project/bacalhau/pkg/websockets"
)

func (a *APIServer) websocket(res http.ResponseWriter, req *http.Request) {
	_, span := system.GetSpanFromRequest(req, "apiServer/websocket")
	defer span.End()

	websockets.AddClientListener(h, res, req)
}
