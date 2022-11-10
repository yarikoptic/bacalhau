package publicapi

import (
	"context"
	"net/http"
	"strings"

	"github.com/filecoin-project/bacalhau/pkg/model"
	"github.com/filecoin-project/bacalhau/pkg/system"
	"github.com/pkg/errors"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/rs/zerolog/log"
)

//nolint:unused
type webSocketRequest struct {
	ClientID string `json:"client_id"`
}

//nolint:unused
type webSocketResponse struct {
	WSEvent *model.WebSocketEvent `json:"web_socket_event"`
}

//nolint:unused
func (apiServer *APIServer) websocket(res http.ResponseWriter, req *http.Request) {
	ctx, span := system.GetSpanFromRequest(req, "apiServer/websocket")
	defer span.End()

	conn, _, _, err := ws.UpgradeHTTP(req, res)
	if err != nil {
		log.Error().AnErr("error", errors.Wrap(err, "events ws handler: failed to upgrade"))
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer conn.Close()
	_ = make(chan *model.WebSocketEvent)
	_, rootCancel := context.WithCancel(ctx)

	// listening for control messages and updating subscriptions
	go func() {
		var _ context.Context
		var cancel context.CancelFunc

		for {
			select {
			case <-req.Context().Done():
				if cancel != nil {
					cancel()
				}
				return
			default:
				_, op, err := wsutil.ReadClientData(conn)
				if err != nil {
					if !strings.Contains(err.Error(), "EOF") && !strings.Contains(err.Error(), "closed") {
						log.Warn().Msgf("error from wsutil.ReadClientData(conn): %s", err)
					}

					// handle error
					if cancel != nil {
						cancel()
					}
					rootCancel()
					return
				}
				if op.IsData() {
					// err = s.jsonSerializer.Decode(msg, &wc)
					if err != nil {
						log.Warn().Msgf("error from Decode: %s", err)
						continue
					}

					// switch wc.ActionType {
					// case "subscribe":
					// 	// canceling previous subscription (if any)
					// 	if cancel != nil {
					// 		cancel()
					// 	}
					// 	ctx, cancel = context.WithCancel(req.Context())
					// 	s.subscribeToEventsFeed(ctx, account.ID, events)

					// case "unsubscribe":
					// 	if cancel != nil {
					// 		cancel()
					// 	}
					// }
				}
			}
		}
	}()

	_ = rootCancel
	// for {
	// 	select {
	// 	case <-rootCtx.Done():
	// 		// finishing work
	// 		return
	// 	case event, ok := <-events:
	// 		if !ok {
	// 			// closing conn
	// 			return
	// 		}
	// 		encoded, err := s.jsonSerializer.Encode(event)
	// 		if err != nil {
	// 			continue
	// 		}
	// 		err = wsutil.WriteServerMessage(conn, ws.OpText, encoded)
	// 		if err != nil {
	// 			// handle error
	// 			return
	// 		}
	// 	}
	// }
}
