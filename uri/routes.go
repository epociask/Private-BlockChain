package uri

import (
	"chain/uri/handlers"
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Start",
		"GET",
		"/start",
		handlers.Start,
	},
	Route{
		"Get Block",
		"GET",
		"/block/{height}/{hash}",
		handlers.GetBlock,
	},
	Route{
		"Show",
		"GET",
		"/show",
		handlers.Show,
	},
	Route{
		"Receive Block",
		"POST",
		"/block/receive",
		handlers.ReceiveBlock,
	},
	Route{
		"Peer",
		"POST",
		"/peers",
		handlers.Register,
	},
	Route{
		"ShowPeers",
		"GET",
		"/showpeers",
		handlers.ShowPeers,
	},
	Route{
		"Download",
		"GET",
		"/download",
		handlers.Download,
	},
}
