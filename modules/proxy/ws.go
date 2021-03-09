package ws

import (
	"github.com/titaniumnetwork-dev/Aurora/modules/config"
	"github.com/gobwas/ws"
	"net/http"
}

// HTTP Server used to forward websocket connection
func Server(w http.ResponseWrite, r *http.Request) {
	// TODO: Redirect to websocket server and proxify connects	
}
