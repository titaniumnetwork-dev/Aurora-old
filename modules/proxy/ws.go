package proxy

import (
	"context"
	"fmt"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/titaniumnetwork-dev/Aurora/modules/config"
	"net/http"
)

// Server used for ws proxy
// Even after this is done there are a lot of improvements to be made
// Actually send websocket error however that works
func WSServer(w http.ResponseWriter, r *http.Request) {
	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		log.Println(err)
		return
	}
	go func() {
		defer conn.Close()

		for {
			msg, op, err := wsutil.ReadClientData(serverConn)
			if err != nil {
				log.Println(err)
				return
			}

			clientConn, _, _, err := ws.DefaultDialer.Dial(ctx, config.WSProxyURL.String())
			if err != nil {
				log.Println(err)
				return
			}

			clientMsg, clientOp, err = wsutil.ReadServerData(clientConn)
			if err != nil {
				log.Println(err)
				return
			}
		}
	}()
}
