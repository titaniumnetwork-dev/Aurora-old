package proxy

import (
	//"context"
	"fmt"
	//"github.com/gobwas/ws"
	//"github.com/gobwas/ws/wsutil"
	//"github.com/titaniumnetwork-dev/Aurora/modules/config"
	"net/http"
)

// Server used for ws proxy
// Even after this is done there are a lot of improvements to be made
// Do not use pre-forwarding blocking or get information at that time
// TODO: Send websocket error however that works
func WSServer(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "Websocket support is not finished yet")
	return
	/*
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

			// TODO: STtore something in proxy proxyurl variable
			proxyConn, _, _, err := ws.DefaultDialer.Dial(ctx, config.ProxyURL.String())
			if err != nil {
				log.Println(err)
				return
			}

			proxyMsg, proxyOp, err = wsutil.ReadServerData(clientConn)
			if err != nil {
				log.Println(err)
				return
			}
		}
	}()
	*/
}
