package server

import (
	"net/http"

	"github.com/bilibili/twirp"
	"github.com/busyfree/shorturl-go/cmd/server/hook"
	"github.com/busyfree/shorturl-go/server/webgin"
)

var hooks = twirp.ChainHooks(
	hook.NewRequestID(),
	hook.NewLog(),
)

var privateHooks = twirp.ChainHooks(
	hook.NewRequestID(),
	hook.NewLog(),
)

var allowGetHooks = twirp.ChainHooks(
	hook.NeAllowGet(),
	hook.NewRequestID(),
	hook.NewLog(),
)

func initMux(mux *http.ServeMux, isInternal bool) {
	{
		webgin.InitWebGin()
		mux.Handle("/", webgin.GinRoute)
	}
}

func initInternalMux(mux *http.ServeMux) {
}
